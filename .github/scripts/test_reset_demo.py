"""Exercise replacement and failure recovery without contacting AWS."""

import contextlib
import io
import os
import subprocess
import unittest
from unittest.mock import Mock, patch

import reset_demo


ACTIVE = "CasaOS-Demo-200"
ORPHAN = "CasaOS-Demo-100"
NEW = "CasaOS-Demo-999-1"
ADDRESS = "203.0.113.10"
CANDIDATE_ADDRESS = "203.0.113.11"
ENVIRONMENT = {
    "AWS_REGION": "us-west-2",
    "DEMO_INSTANCE_NAME": NEW,
    "DEMO_SNAPSHOT_NAME": "source-snapshot",
    "DEMO_STATIC_IP_NAME": "CasaOS-Demo-IP",
    "DEMO_AVAILABILITY_ZONE": "us-west-2a",
    "DEMO_BUNDLE_ID": "medium_2_0",
    "DEMO_HOST": "demo.casaos.io",
}


class FakeLightsail:
    def __init__(self):
        self.owner = ACTIVE
        self.instances = {ACTIVE, ORPHAN, "CasaOS-Demo-database"}
        self.protected = set()
        self.events = []
        self.fail_create = False
        self.fail_cutover = False
        self.fail_rollback = False
        self.fail_delete = None
        self.failed_operation = None
        self.healthy_candidate = True
        self.healthy_cutover = True
        self.snapshot_state = "available"
        self.clock = 0

    def sleep(self, seconds):
        self.clock += seconds

    def instance(self, name):
        return {
            "name": name,
            "state": {"name": "running"},
            "publicIpAddress": CANDIDATE_ADDRESS,
            "isStaticIp": name == self.owner or name in self.protected,
        }

    def aws(self, command, *arguments):
        self.events.append((command, arguments))
        if command == "get-static-ip":
            return {"staticIp": {"attachedTo": self.owner, "ipAddress": ADDRESS}}
        if command == "get-instance":
            name = arguments[1]
            if name not in self.instances:
                raise subprocess.CalledProcessError(1, "aws", stderr="NotFoundException")
            return {"instance": self.instance(name)}
        if command == "get-instances":
            return {"instances": [self.instance(name) for name in sorted(self.instances)]}
        if command == "get-instance-snapshot":
            return {"instanceSnapshot": {"state": self.snapshot_state}}
        if command == "get-operation":
            status = "Failed" if arguments[1] == self.failed_operation else "Succeeded"
            return {"operation": {"status": status, "errorDetails": "operation failed"}}
        if command == "create-instances-from-snapshot":
            self.instances.add(arguments[arguments.index("--instance-names") + 1])
            if self.fail_create:
                raise RuntimeError("Create response lost after AWS accepted the request")
        elif command == "attach-static-ip":
            name = arguments[arguments.index("--instance-name") + 1]
            if name == ACTIVE and self.fail_rollback:
                raise RuntimeError("Rollback unavailable")
            self.owner = name
            if name == NEW and self.fail_cutover:
                raise RuntimeError("Attach response lost after AWS applied the request")
        elif command == "delete-instance":
            name = arguments[1]
            if name == self.fail_delete:
                raise RuntimeError("Deletion unavailable")
            self.instances.remove(name)
        elif command != "put-instance-public-ports":
            raise AssertionError(f"Unexpected command {command}")
        operation = {"id": command}
        # The firewall API uses a singular response, unlike create/attach/delete.
        return {"operation": operation} if command == "put-instance-public-ports" else {"operations": [operation]}

    def http_ready(self, address):
        self.events.append(("http", (address, self.owner)))
        if address == CANDIDATE_ADDRESS:
            return self.healthy_candidate
        return self.owner == ACTIVE or self.healthy_cutover


class ResetTest(unittest.TestCase):
    def setUp(self):
        self.cloud = FakeLightsail()
        patches = [
            patch.dict(os.environ, ENVIRONMENT),
            patch.object(reset_demo, "aws", self.cloud.aws),
            patch.object(reset_demo, "http_ready", self.cloud.http_ready),
            patch.object(reset_demo.time, "monotonic", lambda: self.cloud.clock),
            patch.object(reset_demo.time, "sleep", self.cloud.sleep),
            contextlib.redirect_stdout(io.StringIO()),
        ]
        for item in patches:
            item.__enter__()
            self.addCleanup(item.__exit__, None, None, None)

    def deleted(self):
        return [args[1] for command, args in self.cloud.events if command == "delete-instance"]

    def test_healthy_replacement_removes_all_legacy_orphans_after_cutover(self):
        reset_demo.reset()
        self.assertEqual(self.cloud.owner, NEW)
        self.assertEqual(self.cloud.instances, {NEW, "CasaOS-Demo-database"})
        self.assertEqual(self.deleted(), [ORPHAN, ACTIVE])
        healthy = self.cloud.events.index(("http", (ADDRESS, NEW)))
        deletion = next(i for i, event in enumerate(self.cloud.events) if event[0] == "delete-instance")
        self.assertLess(healthy, deletion)

    def test_unhealthy_candidate_never_receives_static_ip(self):
        self.cloud.healthy_candidate = False
        with self.assertRaises(TimeoutError):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, ACTIVE)
        self.assertEqual(self.deleted(), [NEW])
        self.assertFalse(any(command == "attach-static-ip" for command, _ in self.cloud.events))

    def test_unhealthy_static_ip_rolls_back_before_deleting_candidate(self):
        self.cloud.healthy_cutover = False
        with self.assertRaises(TimeoutError):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, ACTIVE)
        self.assertEqual(self.deleted(), [NEW])
        rollback_health = self.cloud.events.index(("http", (ADDRESS, ACTIVE)))
        deletion = self.cloud.events.index(("delete-instance", ("--instance-name", NEW)))
        self.assertLess(rollback_health, deletion)

    def test_attach_error_after_side_effect_still_rolls_back(self):
        self.cloud.fail_cutover = True
        with self.assertRaisesRegex(RuntimeError, "Attach response lost"):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, ACTIVE)
        self.assertEqual(self.deleted(), [NEW])

    def test_failed_rollback_preserves_both_instances(self):
        self.cloud.healthy_cutover = False
        self.cloud.fail_rollback = True
        with self.assertRaises(TimeoutError):
            reset_demo.reset()
        self.assertIn(ACTIVE, self.cloud.instances)
        self.assertIn(NEW, self.cloud.instances)
        self.assertEqual(self.deleted(), [])

    def test_create_error_after_side_effect_cleans_candidate_only(self):
        self.cloud.fail_create = True
        with self.assertRaisesRegex(RuntimeError, "Create response lost"):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, ACTIVE)
        self.assertEqual(self.deleted(), [NEW])

    def test_failed_firewall_operation_is_not_treated_as_success(self):
        self.cloud.failed_operation = "put-instance-public-ports"
        with self.assertRaisesRegex(RuntimeError, "operation failed"):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, ACTIVE)
        self.assertEqual(self.deleted(), [NEW])

    def test_partial_retirement_keeps_new_demo_and_next_reset_retries_orphans(self):
        self.cloud.fail_delete = ACTIVE
        with self.assertRaisesRegex(RuntimeError, "Deletion unavailable"):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, NEW)
        self.assertNotIn(ORPHAN, self.cloud.instances)
        self.assertIn(ACTIVE, self.cloud.instances)
        self.cloud.fail_delete = None
        with patch.dict(os.environ, {"DEMO_INSTANCE_NAME": "CasaOS-Demo-1000-1"}):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, "CasaOS-Demo-1000-1")
        self.assertEqual(self.cloud.instances, {"CasaOS-Demo-1000-1", "CasaOS-Demo-database"})

    def test_different_static_ip_protects_matching_instance(self):
        self.cloud.protected.add(ORPHAN)
        with self.assertRaisesRegex(RuntimeError, "Refusing to delete"):
            reset_demo.reset()
        self.assertIn(ORPHAN, self.cloud.instances)
        self.assertEqual(self.cloud.owner, NEW)
        self.assertNotIn(ACTIVE, self.cloud.instances)

    def test_one_failed_orphan_does_not_block_other_retirements(self):
        self.cloud.fail_delete = ORPHAN
        with self.assertRaisesRegex(RuntimeError, "retirement needs retry"):
            reset_demo.reset()
        self.assertEqual(self.cloud.owner, NEW)
        self.assertIn(ORPHAN, self.cloud.instances)
        self.assertNotIn(ACTIVE, self.cloud.instances)

    def test_foreign_static_ip_owner_aborts_before_mutation(self):
        self.cloud.owner = "another-service"
        with self.assertRaisesRegex(RuntimeError, "existing demo"):
            reset_demo.reset()
        self.assertEqual(self.deleted(), [])
        self.assertFalse(any(command.startswith("create-") for command, _ in self.cloud.events))

    def test_unavailable_snapshot_aborts_before_creation(self):
        self.cloud.snapshot_state = "pending"
        with self.assertRaisesRegex(RuntimeError, "not available"):
            reset_demo.reset()
        self.assertNotIn(NEW, self.cloud.instances)

    def test_existing_candidate_is_not_deleted_by_retry(self):
        self.cloud.instances.add(NEW)
        with self.assertRaisesRegex(RuntimeError, "already exists"):
            reset_demo.reset()
        self.assertEqual(self.deleted(), [])

    def test_changed_static_ip_aborts_retirement(self):
        self.cloud.owner = "another-service"
        with self.assertRaisesRegex(RuntimeError, "ownership changed"):
            reset_demo.delete_unused(ORPHAN, NEW)
        self.assertEqual(self.deleted(), [])


class ApiTest(unittest.TestCase):
    def test_missing_operations_are_not_treated_as_success(self):
        with patch.object(reset_demo, "aws", return_value={"operations": []}):
            with self.assertRaisesRegex(RuntimeError, "no operations"):
                reset_demo.mutate("create-instances-from-snapshot")

    def test_operation_waits_through_started_status(self):
        responses = [
            {"operations": [{"id": "create-1"}]},
            {"operation": {"status": "Started"}},
            {"operation": {"status": "Succeeded"}},
        ]
        with patch.object(reset_demo, "aws", side_effect=responses) as aws:
            with patch.object(reset_demo.time, "sleep") as sleep:
                reset_demo.mutate("create-instances-from-snapshot")
        self.assertEqual(aws.call_count, 3)
        sleep.assert_called_once_with(5)

    def test_access_denied_is_not_mistaken_for_deleted_instance(self):
        error = subprocess.CalledProcessError(1, "aws", stderr="AccessDeniedException")
        with patch.object(reset_demo, "aws", side_effect=error):
            with self.assertRaises(subprocess.CalledProcessError):
                reset_demo.instance(ACTIVE)

    def test_http_requires_casaos_homepage_and_status_200(self):
        connection = Mock()
        response = connection.getresponse.return_value
        with patch.dict(os.environ, ENVIRONMENT):
            with patch.object(reset_demo.http.client, "HTTPConnection", return_value=connection):
                response.status = 200
                response.read.return_value = b"<title>\n CasaOS\n</title>"
                self.assertTrue(reset_demo.http_ready(CANDIDATE_ADDRESS))
                connection.request.assert_called_with("GET", "/", headers={"Host": "demo.casaos.io"})
                response.read.return_value = b"<title>Other service</title>"
                self.assertFalse(reset_demo.http_ready(CANDIDATE_ADDRESS))
                response.status = 503
                response.read.return_value = b"<title>CasaOS</title>"
                self.assertFalse(reset_demo.http_ready(CANDIDATE_ADDRESS))


if __name__ == "__main__":
    unittest.main()
