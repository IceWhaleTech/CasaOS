package route

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/model/system_app"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/docker"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/env_helper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	uuid "github.com/satori/go.uuid"
)

func InitFunction() {
	go checkSystemApp()
	Update2_3()
	CheckSerialDiskMount()

}

var syncIsExistence = false

func installSyncthing(appId string) {

	var appInfo model.ServerAppList
	m := model.CustomizationPostData{}
	var dockerImage string
	var dockerImageVersion string
	appInfo = service.MyService.OAPI().GetServerAppInfo(appId, "system", "us_en")
	dockerImage = appInfo.Image
	dockerImageVersion = appInfo.ImageVersion

	if len(appInfo.ImageVersion) == 0 {
		dockerImageVersion = "latest"
	}

	if appInfo.NetworkModel != "host" {
		for i := 0; i < len(appInfo.Ports); i++ {
			if p, _ := strconv.Atoi(appInfo.Ports[i].ContainerPort); port.IsPortAvailable(p, appInfo.Ports[i].Protocol) {
				appInfo.Ports[i].CommendPort = strconv.Itoa(p)
			} else {
				if appInfo.Ports[i].Protocol == "tcp" {
					if p, err := port.GetAvailablePort("tcp"); err == nil {
						appInfo.Ports[i].CommendPort = strconv.Itoa(p)
					}
				} else if appInfo.Ports[i].Protocol == "upd" {
					if p, err := port.GetAvailablePort("udp"); err == nil {
						appInfo.Ports[i].CommendPort = strconv.Itoa(p)
					}
				}
			}

			if appInfo.Ports[i].Type == 0 {
				appInfo.PortMap = appInfo.Ports[i].CommendPort
			}
		}
	}

	for i := 0; i < len(appInfo.Devices); i++ {
		if !file.CheckNotExist(appInfo.Devices[i].ContainerPath) {
			appInfo.Devices[i].Path = appInfo.Devices[i].ContainerPath
		}
	}
	if len(appInfo.Tip) > 0 {
		appInfo.Tip = env_helper.ReplaceStringDefaultENV(appInfo.Tip)
	}

	appInfo.MaxMemory = service.MyService.ZiMa().GetMemInfo().Total >> 20

	id := uuid.NewV4().String()

	installLog := model2.AppNotify{}

	// step：下载镜像
	err := service.MyService.Docker().DockerPullImage(dockerImage+":"+dockerImageVersion, installLog)
	if err != nil {
		//pull image error
		fmt.Println("pull image error", err, dockerImage, dockerImageVersion)
		return
	}
	for !service.MyService.Docker().IsExistImage(dockerImage + ":" + dockerImageVersion) {
		time.Sleep(time.Second)
	}

	m.CpuShares = 50
	m.Envs = appInfo.Envs
	m.Memory = int64(appInfo.MaxMemory)
	m.Origin = "system"
	m.PortMap = appInfo.PortMap
	m.Ports = appInfo.Ports
	m.Restart = "always"
	m.Volumes = appInfo.Volumes

	containerId, err := service.MyService.Docker().DockerContainerCreate(dockerImage+":"+dockerImageVersion, id, m, appInfo.NetworkModel)
	if err != nil {
		fmt.Println("container create error", err)
		// create container error
		return
	}

	//step：start container
	err = service.MyService.Docker().DockerContainerStart(id)
	if err != nil {
		//start container error
		return
	}

	portsStr, _ := json.Marshal(appInfo.Ports)
	envsStr, _ := json.Marshal(appInfo.Envs)
	volumesStr, _ := json.Marshal(appInfo.Volumes)
	devicesStr, _ := json.Marshal(appInfo.Devices)
	//step: 保存数据到数据库
	md := model2.AppListDBModel{
		CustomId: id,
		Title:    appInfo.Title,
		//ScreenshotLink: appInfo.ScreenshotLink,
		Slogan:      appInfo.Tagline,
		Description: appInfo.Description,
		//Tags:           appInfo.Tags,
		Icon:        appInfo.Icon,
		Version:     dockerImageVersion,
		ContainerId: containerId,
		Image:       dockerImage,
		Index:       appInfo.Index,
		PortMap:     appInfo.PortMap,
		Label:       appInfo.Title,
		EnableUPNP:  false,
		Ports:       string(portsStr),
		Envs:        string(envsStr),
		Volumes:     string(volumesStr),
		Position:    true,
		NetModel:    appInfo.NetworkModel,
		Restart:     m.Restart,
		CpuShares:   50,
		Memory:      int64(appInfo.MaxMemory),
		Devices:     string(devicesStr),
		Origin:      m.Origin,
		CreatedAt:   strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt:   strconv.FormatInt(time.Now().Unix(), 10),
	}
	service.MyService.App().SaveContainer(md)

	checkSystemApp()
}

// check if the system application is installed
func checkSystemApp() {
	list := service.MyService.App().GetSystemAppList()
	for _, v := range *list {
		if v.Image == "linuxserver/syncthing" {
			if v.State != "running" {
				//step：start container
				service.MyService.Docker().DockerContainerStart(v.CustomId)
			}
			syncIsExistence = true
			if config.SystemConfigInfo.SyncPort != v.Port {
				config.SystemConfigInfo.SyncPort = v.Port
			}
			var paths []model.PathMap
			json.Unmarshal([]byte(v.Volumes), &paths)
			path := ""
			for _, i := range paths {
				if i.ContainerPath == "/config" {
					path = docker.GetDir(v.CustomId, i.Path) + "/config.xml"
					for i := 0; i < 10; i++ {
						if file.CheckNotExist(path) {
							time.Sleep(1 * time.Second)
						} else {
							break
						}
					}
					break
				}
			}
			content := file.ReadFullFile(path)
			syncConfig := &system_app.SyncConfig{}
			xml.Unmarshal(content, &syncConfig)
			config.SystemConfigInfo.SyncKey = syncConfig.Key
		}
	}
	if !syncIsExistence {
		installSyncthing("74")
	}
}
func CheckSerialDiskMount() {
	// check mount point
	dbList := service.MyService.Disk().GetSerialAll()

	list := service.MyService.Disk().LSBLK()
	mountPoint := make(map[string]string, len(dbList))
	//remount
	for _, v := range dbList {
		mountPoint[v.Path] = v.MountPoint
	}
	for _, v := range list {
		command.ExecEnabledSMART(v.Path)
		if v.Children != nil {
			for _, h := range v.Children {
				if len(h.MountPoint) == 0 && len(v.Children) == 1 && h.FsType == "ext4" {
					if m, ok := mountPoint[h.Path]; ok {
						//mount point check
						volume := m
						if !file.CheckNotExist(m) {
							for i := 0; file.CheckNotExist(volume); i++ {
								volume = m + strconv.Itoa(i+1)
							}
						}
						service.MyService.Disk().MountDisk(h.Path, volume)
						if volume != m {
							ms := model2.SerialDisk{}
							ms.Serial = v.Serial
							service.MyService.Disk().UpdateMountPoint(ms)
						}

					}
				}
			}
		}
	}
	service.MyService.Disk().RemoveLSBLKCache()
	command.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;AutoRemoveUnuseDir")

}
func Update2_3() {
	command.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/assist.sh")
}
