"""Replace the Lightsail demo, keeping its predecessor until cutover is healthy."""

import http.client
import json
import os
import re
import signal
import subprocess
import time


# Accept the original timestamp names and the new run-ID/attempt names only.
DEMO_NAME = re.compile(r"CasaOS-Demo-[0-9]+(?:-[0-9]+)?")


def aws(command, *arguments):
    result = subprocess.run(
        [
            "aws", "lightsail", command, *arguments,
            "--region", os.environ["AWS_REGION"],
            "--output", "json", "--no-cli-pager",
            "--cli-connect-timeout", "10", "--cli-read-timeout", "30",
        ],
        capture_output=True,
        text=True,
        check=True,
        timeout=120,
    )
    return json.loads(result.stdout)


def wait_for(description, check, timeout=120):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        result = check()
        if result:
            return result
        time.sleep(5)
    raise TimeoutError(f"Timed out waiting for {description}")


def mutate(command, *arguments, timeout=120):
    # A successful API request can still describe a pending or failed operation.
    response = aws(command, *arguments)
    operations = response.get("operations")
    if operations is None:
        operations = [response["operation"]]
    if not operations:
        raise RuntimeError(f"{command} returned no operations")

    def completed():
        statuses = []
        for operation in operations:
            current = aws("get-operation", "--operation-id", operation["id"])["operation"]
            if current["status"] == "Failed":
                raise RuntimeError(f"{command} failed: {current.get('errorDetails', current)}")
            statuses.append(current["status"])
        return all(status in ("Succeeded", "Completed") for status in statuses)

    wait_for(f"{command} operations", completed, timeout)


def instance(name):
    try:
        return aws("get-instance", "--instance-name", name)["instance"]
    except subprocess.CalledProcessError as error:
        if "NotFoundException" in error.stderr:
            return None
        raise


def static_ip():
    return aws("get-static-ip", "--static-ip-name", os.environ["DEMO_STATIC_IP_NAME"])["staticIp"]


def http_ready(address):
    # Connect directly to the candidate, without DNS caching, proxies, or redirects.
    connection = http.client.HTTPConnection(address, timeout=10)
    try:
        connection.request("GET", "/", headers={"Host": os.environ["DEMO_HOST"]})
        response = connection.getresponse()
        body = response.read(65536)
        return response.status == 200 and bool(
            re.search(br"<title>\s*CasaOS\s*</title>", body, re.IGNORECASE)
        )
    except (OSError, http.client.HTTPException):
        return False
    finally:
        connection.close()


def delete_unused(name, serving_instance):
    if static_ip().get("attachedTo") != serving_instance:
        raise RuntimeError("Static IP ownership changed; aborting instance cleanup")
    current = instance(name)
    if current is None:
        return
    if name == serving_instance or current["isStaticIp"]:
        raise RuntimeError(f"Refusing to delete an instance with an attached static IP: {name}")
    print(f"Deleting unused demo instance {name}", flush=True)
    mutate("delete-instance", "--instance-name", name)
    wait_for(f"{name} deletion", lambda: instance(name) is None)


def recover(previous_name, new_name, cutover_attempted, address):
    if cutover_attempted:
        owner = static_ip().get("attachedTo")
        if owner not in (previous_name, new_name, None):
            raise RuntimeError("Static IP moved outside this reset; preserving both instances")
        print(f"Restoring static IP to {previous_name}", flush=True)
        mutate(
            "attach-static-ip", "--static-ip-name", os.environ["DEMO_STATIC_IP_NAME"],
            "--instance-name", previous_name,
        )
        wait_for("rollback IP attachment", lambda: static_ip().get("attachedTo") == previous_name)
        wait_for("previous demo HTTP readiness", lambda: http_ready(address))
    delete_unused(new_name, previous_name)


def reset():
    new_name = os.environ["DEMO_INSTANCE_NAME"]
    if not DEMO_NAME.fullmatch(new_name):
        raise ValueError(f"Invalid demo instance name: {new_name}")
    original_ip = static_ip()
    previous_name = original_ip.get("attachedTo") or ""
    if not DEMO_NAME.fullmatch(previous_name):
        raise RuntimeError("Demo static IP must be attached to an existing demo instance")
    previous = instance(previous_name)
    if previous is None or previous["state"]["name"] != "running":
        raise RuntimeError("The current demo must be running before it can be replaced")
    if instance(new_name) is not None:
        raise RuntimeError(f"Candidate already exists: {new_name}")
    snapshot = aws(
        "get-instance-snapshot", "--instance-snapshot-name", os.environ["DEMO_SNAPSHOT_NAME"],
    )["instanceSnapshot"]
    if snapshot["state"] != "available":
        raise RuntimeError("Demo source snapshot is not available")
    # Capture candidates before creation so an external/newer deployment is never swept.
    old_names = {
        item["name"] for item in aws("get-instances")["instances"]
        if DEMO_NAME.fullmatch(item["name"])
    }
    if previous_name not in old_names:
        raise RuntimeError("Current demo is missing from instance inventory")

    cutover_attempted = False
    try:
        print(f"Creating demo candidate {new_name}", flush=True)
        mutate(
            "create-instances-from-snapshot",
            "--instance-snapshot-name", os.environ["DEMO_SNAPSHOT_NAME"],
            "--instance-names", new_name,
            "--availability-zone", os.environ["DEMO_AVAILABILITY_ZONE"],
            "--bundle-id", os.environ["DEMO_BUNDLE_ID"],
            timeout=600,
        )

        def running_candidate():
            candidate = instance(new_name)
            if candidate and candidate["state"]["name"] == "running":
                return candidate.get("publicIpAddress")
            return None

        candidate_ip = wait_for("candidate running with a public IPv4 address", running_candidate, 180)
        mutate(
            "put-instance-public-ports", "--instance-name", new_name,
            "--port-infos", "fromPort=0,toPort=65535,protocol=all",
        )
        wait_for("candidate CasaOS HTTP readiness", lambda: http_ready(candidate_ip), 300)
        if static_ip().get("attachedTo") != previous_name:
            raise RuntimeError("Static IP ownership changed before cutover")
        # Set this before the API call, which may apply a change then return an error.
        cutover_attempted = True
        mutate(
            "attach-static-ip", "--static-ip-name", os.environ["DEMO_STATIC_IP_NAME"],
            "--instance-name", new_name,
        )
        wait_for("candidate IP attachment", lambda: static_ip().get("attachedTo") == new_name)
        wait_for("demo HTTP readiness through static IP", lambda: http_ready(original_ip["ipAddress"]))
    except BaseException:
        try:
            recover(previous_name, new_name, cutover_attempted, original_ip["ipAddress"])
        except Exception as error:
            print(f"::error::Recovery incomplete; inspect both instances before deleting: {error}", flush=True)
        raise

    # Cutover is committed. Cleanup failure must not roll traffic back to a deleted host.
    # Retry all pre-existing matching instances, including orphans from failed runs.
    cleanup_errors = []
    for old_name in sorted(old_names):
        try:
            delete_unused(old_name, new_name)
        except Exception as error:
            cleanup_errors.append(f"{old_name}: {error}")
    if cleanup_errors:
        raise RuntimeError("Demo is serving, but retirement needs retry: " + "; ".join(cleanup_errors))
    print(f"Demo reset completed; serving instance is {new_name}", flush=True)


def terminate(signum, _frame):
    raise RuntimeError(f"Reset interrupted by signal {signum}")


if __name__ == "__main__":
    signal.signal(signal.SIGTERM, terminate)
    reset()
