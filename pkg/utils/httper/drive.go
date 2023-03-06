package httper

import (
	"time"
)

type MountList struct {
	MountPoints []MountPoints `json:"mountPoints"`
}
type MountPoints struct {
	MountPoint string `json:"MountPoint"`
	Fs         string `json:"Fs"`
	Icon       string `json:"Icon"`
	Name       string `json:"Name"`
}
type MountPoint struct {
	MountPoint string `json:"mount_point"`
	Fs         string `json:"fs"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
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

// func NewRestyClient() *resty.Client {

// 	unixSocket := "/var/run/rclone/rclone.sock"

// 	transport := http.Transport{
// 		Dial: func(_, _ string) (net.Conn, error) {
// 			return net.Dial("unix", unixSocket)
// 		},
// 	}

// 	client := resty.New()

// 	client.SetTransport(&transport).SetBaseURL("http://localhost")
// 	client.SetRetryCount(3).SetRetryWaitTime(5*time.Second).SetTimeout(DefaultTimeout).SetHeader("User-Agent", UserAgent)
// 	return client
// }

//	func GetMountList() (MountList, error) {
//		var result MountList
//		res, err := NewRestyClient().R().Post("/mount/listmounts")
//		if err != nil {
//			return result, err
//		}
//		if res.StatusCode() != 200 {
//			return result, fmt.Errorf("get mount list failed")
//		}
//		json.Unmarshal(res.Body(), &result)
//		for i := 0; i < len(result.MountPoints); i++ {
//			result.MountPoints[i].Fs = result.MountPoints[i].Fs[:len(result.MountPoints[i].Fs)-1]
//		}
//		return result, err
//	}
//
//	func Mount(mountPoint string, fs string) error {
//		res, err := NewRestyClient().R().SetFormData(map[string]string{
//			"mountPoint": mountPoint,
//			"fs":         fs,
//			"mountOpt":   `{"AllowOther": true}`,
//		}).Post("/mount/mount")
//		if err != nil {
//			return err
//		}
//		if res.StatusCode() != 200 {
//			return fmt.Errorf("mount failed")
//		}
//		logger.Info("mount then", zap.Any("res", res.Body()))
//		return nil
//	}
// func Unmount(mountPoint string) error {
// 	res, err := NewRestyClient().R().SetFormData(map[string]string{
// 		"mountPoint": mountPoint,
// 	}).Post("/mount/unmount")
// 	if err != nil {
// 		logger.Error("when unmount", zap.Error(err))
// 		return err
// 	}
// 	if res.StatusCode() != 200 {
// 		logger.Error("then unmount failed", zap.Any("res", res.Body()))
// 		return fmt.Errorf("unmount failed")
// 	}
// 	logger.Info("unmount then", zap.Any("res", res.Body()))
// 	return nil
// }

// func CreateConfig(data map[string]string, name, t string) error {
// 	data["config_is_local"] = "false"
// 	dataStr, _ := json.Marshal(data)
// 	res, err := NewRestyClient().R().SetFormData(map[string]string{
// 		"name":       name,
// 		"parameters": string(dataStr),
// 		"type":       t,
// 	}).Post("/config/create")
// 	logger.Info("when create config then", zap.Any("res", res.Body()))
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode() != 200 {
// 		return fmt.Errorf("create config failed")
// 	}

// 	return nil
// }

// func GetConfigByName(name string) (map[string]string, error) {

//		res, err := NewRestyClient().R().SetFormData(map[string]string{
//			"name": name,
//		}).Post("/config/get")
//		if err != nil {
//			return nil, err
//		}
//		if res.StatusCode() != 200 {
//			return nil, fmt.Errorf("create config failed")
//		}
//		var result map[string]string
//		json.Unmarshal(res.Body(), &result)
//		return result, nil
//	}
// func GetAllConfigName() (RemotesResult, error) {
// 	var result RemotesResult
// 	res, err := NewRestyClient().R().SetFormData(map[string]string{}).Post("/config/listremotes")
// 	if err != nil {
// 		return result, err
// 	}
// 	if res.StatusCode() != 200 {
// 		return result, fmt.Errorf("get config failed")
// 	}

//		json.Unmarshal(res.Body(), &result)
//		return result, nil
//	}
// func DeleteConfigByName(name string) error {
// 	res, err := NewRestyClient().R().SetFormData(map[string]string{
// 		"name": name,
// 	}).Post("/config/delete")
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode() != 200 {
// 		return fmt.Errorf("delete config failed")
// 	}
// 	return nil
// }
