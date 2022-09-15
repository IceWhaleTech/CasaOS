package common

import "testing"

func TestSendNotify(t *testing.T) {
	notify, err := NewNotifyService("/var/run/casaos")
	if err != nil {
		t.Fatal(err)
	}
	err = notify.SendNotify("test", map[string]interface{}{
		"test": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendSystemStatusNotify(t *testing.T) {
	notify, err := NewNotifyService("/var/run/casaos")
	if err != nil {
		t.Fatal(err)
	}
	err = notify.SendSystemStatusNotify(map[string]interface{}{
		"sys_usb": `[{"name": "sdc","size": 7747397632,"model": "DataTraveler_2.0","avail": 7714418688,"children": null}]`,
	})
	if err != nil {
		t.Fatal(err)
	}
}
