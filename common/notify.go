package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	CasaOSURLFilename = "casaos.url"
	APICasaOSNotify   = "/v1/notify"
)

type NotifyService interface {
	SendNotify(path string, message map[string]interface{}) error
	SendSystemStatusNotify(message map[string]interface{}) error
}
type notifyService struct {
	address string
}

func (n *notifyService) SendNotify(path string, message map[string]interface{}) error {
	url := strings.TrimSuffix(n.address, "/") + APICasaOSNotify + "/" + path
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("failed to send notify (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}
	return nil
}

// disk: "sys_disk":{"size":56866869248,"avail":5855485952,"health":true,"used":48099700736}
// usb:   "sys_usb":[{"name": "sdc","size": 7747397632,"model": "DataTraveler_2.0","avail": 7714418688,"children": null}]
func (n *notifyService) SendSystemStatusNotify(message map[string]interface{}) error {
	url := strings.TrimSuffix(n.address, "/") + APICasaOSNotify + "/system_status"

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("failed to send notify (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}
	return nil
}

func NewNotifyService(runtimePath string) (NotifyService, error) {
	casaosAddressFile := filepath.Join(runtimePath, CasaOSURLFilename)

	buf, err := os.ReadFile(casaosAddressFile)
	if err != nil {
		return nil, err
	}

	address := string(buf)

	response, err := http.Get(address + "/ping")
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to ping casaos service")
	}

	return &notifyService{
		address: address,
	}, nil
}
