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
	APICasaOSNotify   = "/v1/notiry"
)

type NotifyService interface {
	SendNotify(path string, message map[string]interface{}) error
	SendSystemNotify(message map[string]interface{}) error
}
type notifyService struct {
	address string
}

func (n *notifyService) SendNotify(path string, message map[string]interface{}) error {

	url := strings.TrimSuffix(n.address, "/") + "/" + APICasaOSNotify + "/" + path
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("failed to send notify (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}
	return nil

}
func (n *notifyService) SendSystemNotify(message map[string]interface{}) error {

	url := strings.TrimSuffix(n.address, "/") + "/" + APICasaOSNotify
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
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

	if response.StatusCode != 200 {
		return nil, errors.New("failed to ping casaos service")
	}

	return &notifyService{
		address: address,
	}, nil
}
