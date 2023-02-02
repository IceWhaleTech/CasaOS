package httper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type MountList struct {
	MountPoints []struct {
		MountPoint string `json:"MountPoint"`
		Fs         string `json:"Fs"`
		Icon       string `json:"Icon"`
	} `json:"mountPoints"`
}
type MountPoint struct {
	MountPoint string `json:"mount_point"`
	Fs         string `json:"fs"`
	Icon       string `json:"icon"`
}
type MountResult struct {
	Error string `json:"error"`
	Input struct {
		Fs         string `json:"fs"`
		MountPoint string `json:"mountPoint"`
	} `json:"input"`
	Path   string `json:"path"`
	Status int    `json:"status"`
}

type RemotesResult struct {
	Remotes []string `json:"remotes"`
}

var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
var DefaultTimeout = time.Second * 30

func NewRestyClient() *resty.Client {

	unixSocket := "/tmp/rclone.sock"

	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", unixSocket)
		},
	}

	client := resty.New()

	client.SetTransport(&transport).SetBaseURL("http://localhost")
	client.SetRetryCount(3).SetRetryWaitTime(5*time.Second).SetTimeout(DefaultTimeout).SetHeader("User-Agent", UserAgent)
	return client
}

func GetMountList() (MountList, error) {
	var result MountList
	res, err := NewRestyClient().R().Post("/mount/listmounts")
	if err != nil {
		return result, err
	}
	if res.StatusCode() != 200 {
		return result, fmt.Errorf("get mount list failed")
	}
	json.Unmarshal(res.Body(), &result)
	for i := 0; i < len(result.MountPoints); i++ {
		result.MountPoints[i].Fs = result.MountPoints[i].Fs[:len(result.MountPoints[i].Fs)-1]
	}
	return result, err
}
func Mount(mountPoint string, fs string) error {
	res, err := NewRestyClient().R().SetFormData(map[string]string{
		"mountPoint": mountPoint,
		"fs":         fs,
	}).Post("/mount/mount")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("mount failed")
	}
	return nil
}
func Unmount(mountPoint string) error {
	res, err := NewRestyClient().R().SetFormData(map[string]string{
		"mountPoint": mountPoint,
	}).Post("/mount/unmount")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("unmount failed")
	}
	return nil
}

func CreateConfig(data map[string]string, name, t string) error {
	data["config_is_local"] = "false"
	dataStr, _ := json.Marshal(data)
	res, err := NewRestyClient().R().SetFormData(map[string]string{
		"name":       name,
		"parameters": string(dataStr),
		"type":       t,
	}).Post("/config/create")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("create config failed")
	}
	return nil
}

func GetConfigByName(name string) (map[string]string, error) {

	res, err := NewRestyClient().R().SetFormData(map[string]string{
		"name": name,
	}).Post("/config/get")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("create config failed")
	}
	var result map[string]string
	json.Unmarshal(res.Body(), &result)
	return result, nil
}
func GetAllConfigName() (RemotesResult, error) {
	var result RemotesResult
	res, err := NewRestyClient().R().SetFormData(map[string]string{}).Post("/config/listremotes")
	if err != nil {
		return result, err
	}
	if res.StatusCode() != 200 {
		return result, fmt.Errorf("get config failed")
	}

	json.Unmarshal(res.Body(), &result)
	return result, nil
}
func DeleteConfigByName(name string) error {
	res, err := NewRestyClient().R().SetFormData(map[string]string{
		"name": name,
	}).Post("/config/delete")
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("delete config failed")
	}
	return nil
}
