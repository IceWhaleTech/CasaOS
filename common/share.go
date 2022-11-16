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
	APICasaOSShare = "/v1/samba/shares"
)

type ShareService interface {
	DeleteShare(id string) error
}
type shareService struct {
	address string
}

func (n *shareService) DeleteShare(id string) error {
	url := strings.TrimSuffix(n.address, "/") + APICasaOSShare + "/" + id
	fmt.Println(url)
	message := "{}"
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Fetch Request
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("failed to send share (status code: " + fmt.Sprint(response.StatusCode) + ")")
	}
	return nil

}

func NewShareService(runtimePath string) (ShareService, error) {
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

	return &shareService{
		address: address,
	}, nil
}
