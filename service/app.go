package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	client2 "github.com/docker/docker/client"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppService interface {
	GetMyList(index, size int, position bool) (*[]model2.MyAppList, *[]model2.MyAppList)
	SaveContainer(m model2.AppListDBModel)
	GetUninstallInfo(id string) model2.AppListDBModel
	DeleteApp(id string)
	GetContainerInfo(id string) (types.Container, error)
	GetAppDBInfo(id string) model2.AppListDBModel
	UpdateApp(m model2.AppListDBModel)
	GetSimpleContainerInfo(id string) (types.Container, error)
	DelAppConfigDir(path string)
	GetSystemAppList() []types.Container
	GetHardwareUsageStream()
	GetHardwareUsage() []model.DockerStatsModel
	GetAppStats(id string) string
	GetAllDBApps() []model2.AppListDBModel
	ImportApplications(casaApp bool)
	CheckNewImage()
}

type appStruct struct {
	db *gorm.DB
}

func (a *appStruct) CheckNewImage() {
	list := MyService.Docker().DockerContainerList()
	for _, v := range list {
		inspect, err := MyService.Docker().DockerImageInfo(strings.Split(v.Image, ":")[0])
		if err != nil {
			NewVersionApp[v.ID] = inspect.ID
			continue
		}
		if inspect.ID == v.ImageID {
			delete(NewVersionApp, v.ID)
			continue
		}
		NewVersionApp[v.ID] = inspect.ID
	}

}
func (a *appStruct) ImportApplications(casaApp bool) {
	if casaApp {
		list := MyService.App().GetAllDBApps()
		for _, app := range list {
			info, err := MyService.Docker().DockerContainerInfo(app.CustomId)
			if err != nil {
				MyService.App().DeleteApp(app.CustomId)
				continue
			}
			//info.NetworkSettings
			info.Config.Labels["casaos"] = "casaos"
			info.Config.Labels["web"] = app.PortMap
			info.Config.Labels["icon"] = app.Icon
			info.Config.Labels["desc"] = app.Description
			info.Config.Labels["index"] = app.Index
			info.Config.Labels["custom_id"] = app.CustomId
			info.Name = app.Title
			container_id, err := MyService.Docker().DockerContainerCopyCreate(info)
			if err != nil {
				fmt.Println(err)
				continue
			}
			MyService.App().DeleteApp(app.CustomId)
			MyService.Docker().DockerContainerStop(app.CustomId)
			MyService.Docker().DockerContainerRemove(app.CustomId, false)
			MyService.Docker().DockerContainerStart(container_id)

		}
	} else {
		list := MyService.Docker().DockerContainerList()
		for _, app := range list {
			info, err := MyService.Docker().DockerContainerInfo(app.ID)
			if err != nil || info.Config.Labels["casaos"] == "casaos" {
				continue
			}
			info.Config.Labels["casaos"] = "casaos"
			info.Config.Labels["web"] = ""
			info.Config.Labels["icon"] = ""
			info.Config.Labels["desc"] = ""
			info.Config.Labels["index"] = ""
			info.Config.Labels["custom_id"] = uuid.NewV4().String()

			_, err = MyService.Docker().DockerContainerCopyCreate(info)
			if err != nil {
				continue
			}

		}
	}

	// allcontainer := MyService.Docker().DockerContainerList()
	// for _, app := range allcontainer {
	// 	info, err := MyService.Docker().DockerContainerInfo(app.ID)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	MyService.Docker().DockerContainerStop(app.ID)
	// 	MyService.Docker().DockerContainerRemove(app.ID, false)
	// 	//info.NetworkSettings
	// 	info.Config.Labels["custom_id"] = uuid.NewV4().String()
	// 	container_id, err := MyService.Docker().DockerContainerCopyCreate(info)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	MyService.Docker().DockerContainerStart(container_id)
	//}
}

//获取我的应用列表
func (a *appStruct) GetMyList(index, size int, position bool) (*[]model2.MyAppList, *[]model2.MyAppList) {
	cli, err := client2.NewClientWithOpts(client2.FromEnv, client2.WithTimeout(time.Second*5))
	if err != nil {
		loger.Error("Failed to init client", zap.Any("err", err))
	}
	defer cli.Close()
	// fts := filters.NewArgs()
	// fts.Add("label", "casaos=casaos")
	//fts.Add("label", "casaos")
	//fts.Add("casaos", "casaos")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		loger.Error("Failed to get container_list", zap.Any("err", err))
	}
	//获取本地数据库应用

	unTranslation := []model2.MyAppList{}

	list := []model2.MyAppList{}

	for _, m := range containers {
		if m.Labels["casaos"] == "casaos" {

			_, newVersion := NewVersionApp[m.ID]
			name := strings.ReplaceAll(m.Names[0], "/", "")
			icon := m.Labels["icon"]
			if len(m.Labels["name"]) > 0 {
				name = m.Labels["name"]
			}
			if m.Labels["origin"] == "system" {
				name = strings.Split(m.Image, ":")[0]
				if len(strings.Split(name, "/")) > 1 {
					icon = "https://icon.casaos.io/main/all/" + strings.Split(name, "/")[1] + ".png"
				}
			}

			list = append(list, model2.MyAppList{
				Name:     name,
				Icon:     icon,
				State:    m.State,
				CustomId: m.Labels["custom_id"],
				Id:       m.ID,
				Port:     m.Labels["web"],
				Index:    m.Labels["index"],
				//Order:      m.Labels["order"],
				Image:  m.Image,
				Latest: newVersion,
				//Type:   m.Labels["origin"],
				//Slogan: m.Slogan,
				//Rely:     m.Rely,
				Host:     m.Labels["host"],
				Protocol: m.Labels["protocol"],
			})
		} else {
			unTranslation = append(unTranslation, model2.MyAppList{
				Name:     strings.ReplaceAll(m.Names[0], "/", ""),
				Icon:     "",
				State:    m.State,
				CustomId: m.ID,
				Id:       m.ID,
				Port:     "",
				Latest:   false,
				Host:     "",
				Protocol: "",
				Image:    m.Image,
			})
		}
	}

	//lMap := make(map[string]interface{})
	// for _, dbModel := range lm {
	// 	if position {
	// 		if dbModel.Position {
	// 			lMap[dbModel.ContainerId] = dbModel
	// 		}
	// 	} else {
	// 		lMap[dbModel.ContainerId] = dbModel
	// 	}
	// }
	// for _, container := range containers {

	// 	if lMap[container.ID] != nil && container.Labels["origin"] != "system" {
	// 		m := lMap[container.ID].(model2.AppListDBModel)
	// 		if len(m.Label) == 0 {
	// 			m.Label = m.Title
	// 		}

	// 		// info, err := cli.ContainerInspect(context.Background(), container.ID)
	// 		// var tm string
	// 		// if err != nil {
	// 		// 	tm = time.Now().String()
	// 		// } else {
	// 		// 	tm = info.State.StartedAt
	// 		//}
	// 		list = append(list, model2.MyAppList{
	// 			Name:     m.Label,
	// 			Icon:     m.Icon,
	// 			State:    container.State,
	// 			CustomId: strings.ReplaceAll(container.Names[0], "/", ""),
	// 			Port:     m.PortMap,
	// 			Index:    m.Index,
	// 			//UpTime:   tm,
	// 			Image:  m.Image,
	// 			Slogan: m.Slogan,
	// 			//Rely:     m.Rely,
	// 		})
	// 	}

	// }

	return &list, &unTranslation

}

//system application list
func (a *appStruct) GetSystemAppList() []types.Container {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		loger.Error("Failed to init client", zap.Any("err", err))
	}
	defer cli.Close()
	fts := filters.NewArgs()
	fts.Add("label", "origin=system")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: fts})
	if err != nil {
		loger.Error("Failed to get container_list", zap.Any("err", err))
	}

	//获取本地数据库应用

	// var lm []model2.AppListDBModel
	// a.db.Table(model2.CONTAINERTABLENAME).Select("title,icon,port_map,`index`,container_id,position,label,slogan,image,volumes").Find(&lm)

	//list := []model2.MyAppList{}
	//lMap := make(map[string]interface{})
	// for _, dbModel := range lm {
	// 	lMap[dbModel.ContainerId] = dbModel
	// }

	return containers

}
func (a *appStruct) GetAllDBApps() []model2.AppListDBModel {
	var lm []model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Select("custom_id,title,icon,container_id,label,slogan,image,port_map").Find(&lm)
	return lm
}

//获取我的应用列表
func (a *appStruct) GetContainerInfo(id string) (types.Container, error) {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		loger.Error("Failed to init client", zap.Any("err", err))
	}
	filters := filters.NewArgs()
	filters.Add("id", id)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		loger.Error("Failed to get container_list", zap.Any("err", err))
	}

	if len(containers) > 0 {
		return containers[0], nil
	}
	return types.Container{}, nil

}

func (a *appStruct) GetSimpleContainerInfo(id string) (types.Container, error) {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return types.Container{}, err
	}
	defer cli.Close()
	filters := filters.NewArgs()
	filters.Add("id", id)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		return types.Container{}, err
	}

	if len(containers) > 0 {
		return containers[0], nil
	}
	return types.Container{}, errors.New("container not existent")
}

//获取我的应用列表
func (a *appStruct) GetAppDBInfo(id string) model2.AppListDBModel {
	var m model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Where("custom_id = ?", id).First(&m)
	return m
}

//根据容器id获取镜像名称
func (a *appStruct) GetUninstallInfo(id string) model2.AppListDBModel {
	var m model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Select("image,version,enable_upnp,ports,envs,volumes,origin").Where("custom_id = ?", id).First(&m)
	return m
}

//创建容器成功后保存容器
func (a *appStruct) SaveContainer(m model2.AppListDBModel) {
	a.db.Table(model2.CONTAINERTABLENAME).Create(&m)
}

func (a *appStruct) UpdateApp(m model2.AppListDBModel) {
	a.db.Table(model2.CONTAINERTABLENAME).Save(&m)
}

func (a *appStruct) DelAppConfigDir(path string) {
	command.OnlyExec("source " + config.AppInfo.ShellPath + "/helper.sh ;DelAppConfigDir " + path)
}

func (a *appStruct) DeleteApp(id string) {
	a.db.Table(model2.CONTAINERTABLENAME).Where("custom_id = ?", id).Delete(&model2.AppListDBModel{})
}

var dataStats sync.Map

var isFinish bool = false

func (a *appStruct) GetAppStats(id string) string {
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return ""
	}
	defer cli.Close()
	con, err := cli.ContainerStats(context.Background(), id, false)
	if err != nil {
		return err.Error()
	}
	defer con.Body.Close()
	c, _ := ioutil.ReadAll(con.Body)
	return string(c)
}

func (a *appStruct) GetHardwareUsage() []model.DockerStatsModel {

	stream := true
	for !isFinish {
		if stream {
			stream = false
			go func() {
				a.GetHardwareUsageStream()
			}()
		}
		runtime.Gosched()
	}
	list := []model.DockerStatsModel{}

	dataStats.Range(func(key, value interface{}) bool {
		list = append(list, value.(model.DockerStatsModel))
		return true
	})
	return list

}

func (a *appStruct) GetHardwareUsageStream() {

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return
	}
	defer cli.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	fts := filters.NewArgs()
	fts.Add("label", "casaos=casaos")
	//fts.Add("status", "running")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: fts})
	if err != nil {
		loger.Error("Failed to get container_list", zap.Any("err", err))
	}
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: fts})
			if err != nil {
				loger.Error("Failed to get container_list", zap.Any("err", err))
				continue
			}
		}
		if config.CasaOSGlobalVariables.AppChange {
			config.CasaOSGlobalVariables.AppChange = false
			dataStats.Range(func(key, value interface{}) bool {
				dataStats.Delete(key)
				return true
			})
		}

		var temp sync.Map
		var wg sync.WaitGroup
		for _, v := range containers {
			if v.State != "running" {
				continue
			}
			wg.Add(1)
			go func(v types.Container, i int) {
				defer wg.Done()
				stats, err := cli.ContainerStatsOneShot(ctx, v.ID)
				if err != nil {
					return
				}
				decode := json.NewDecoder(stats.Body)
				var data interface{}
				if err := decode.Decode(&data); err == io.EOF {
					return
				}
				m, _ := dataStats.Load(v.ID)
				dockerStats := model.DockerStatsModel{}
				if m != nil {
					dockerStats.Previous = m.(model.DockerStatsModel).Data
				}
				dockerStats.Data = data
				dockerStats.Icon = v.Labels["icon"]
				dockerStats.Title = strings.ReplaceAll(v.Names[0], "/", "")

				// @tiger - 不建议直接把依赖的数据结构封装返回。
				//          如果依赖的数据结构有变化，应该在这里适配或者保存，这样更加对客户端负责
				temp.Store(v.ID, dockerStats)
				if i == 99 {
					stats.Body.Close()
				}
			}(v, i)
		}
		wg.Wait()
		dataStats = temp
		isFinish = true

		time.Sleep(time.Second * 1)
	}
	isFinish = false
	cancel()
}

func NewAppService(db *gorm.DB) AppService {
	return &appStruct{db: db}
}
