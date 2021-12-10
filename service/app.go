package service

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	client2 "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AppService interface {
	GetMyList(index, size int, position bool) *[]model2.MyAppList
	SaveContainer(m model2.AppListDBModel)
	GetUninstallInfo(id string) model2.AppListDBModel
	RemoveContainerById(id string)
	GetContainerInfo(name string) (types.Container, error)
	GetAppDBInfo(id string) model2.AppListDBModel
	UpdateApp(m model2.AppListDBModel)
	GetSimpleContainerInfo(name string) (types.Container, error)
	DelAppConfigDir(path string)
	GetSystemAppList() *[]model2.MyAppList
	GetHardwareUsageSteam()
	GetHardwareUsage() []model.DockerStatsModel
	GetAppStats(id string) string
}

type appStruct struct {
	db  *gorm.DB
	log loger2.OLog
}

//获取我的应用列表
func (a *appStruct) GetMyList(index, size int, position bool) *[]model2.MyAppList {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		a.log.Error("初始化client失败", "app.getmylist", "line:36", err)
	}
	defer cli.Close()
	fts := filters.NewArgs()
	fts.Add("label", "origin")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: fts})
	if err != nil {
		a.log.Error("获取docker容器失败", "app.getmylist", "line:42", err)
	}
	//获取本地数据库应用

	var lm []model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Select("title,icon,port_map,`index`,container_id,position,label,slogan,image").Find(&lm)

	list := []model2.MyAppList{}
	lMap := make(map[string]interface{})
	for _, dbModel := range lm {
		if position {
			if dbModel.Position {
				lMap[dbModel.ContainerId] = dbModel
			}
		} else {
			lMap[dbModel.ContainerId] = dbModel
		}
	}
	for _, container := range containers {

		if lMap[container.ID] != nil && container.Labels["origin"] != "system" {
			m := lMap[container.ID].(model2.AppListDBModel)
			if len(m.Label) == 0 {
				m.Label = m.Title
			}

			info, err := cli.ContainerInspect(context.Background(), container.ID)
			var tm string
			if err != nil {
				tm = time.Now().String()
			} else {
				tm = info.State.StartedAt
			}
			list = append(list, model2.MyAppList{
				Name:     m.Label,
				Icon:     m.Icon,
				State:    container.State,
				CustomId: strings.ReplaceAll(container.Names[0], "/", ""),
				Port:     m.PortMap,
				Index:    m.Index,
				UpTime:   tm,
				Image:    m.Image,
				Slogan:   m.Slogan,
				//Rely:     m.Rely,
			})
		}

	}

	return &list

}

//system application list
func (a *appStruct) GetSystemAppList() *[]model2.MyAppList {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		a.log.Error("初始化client失败", "app.getmylist", "line:36", err)
	}
	defer cli.Close()
	fts := filters.NewArgs()
	fts.Add("label", "origin=system")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: fts})
	if err != nil {
		a.log.Error("获取docker容器失败", "app.getmylist", "line:42", err)
	}

	//获取本地数据库应用

	var lm []model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Select("title,icon,port_map,`index`,container_id,position,label,slogan,image,volumes").Find(&lm)

	list := []model2.MyAppList{}
	lMap := make(map[string]interface{})
	for _, dbModel := range lm {
		lMap[dbModel.ContainerId] = dbModel
	}
	for _, container := range containers {

		if lMap[container.ID] != nil {
			m := lMap[container.ID].(model2.AppListDBModel)
			if len(m.Label) == 0 {
				m.Label = m.Title
			}

			info, err := cli.ContainerInspect(context.Background(), container.ID)
			var tm string
			if err != nil {
				tm = time.Now().String()
			} else {
				tm = info.State.StartedAt
			}
			list = append(list, model2.MyAppList{
				Name:     m.Label,
				Icon:     m.Icon,
				State:    container.State,
				CustomId: strings.ReplaceAll(container.Names[0], "/", ""),
				Port:     m.PortMap,
				Index:    m.Index,
				UpTime:   tm,
				Image:    m.Image,
				Slogan:   m.Slogan,
				Volumes:  m.Volumes,
				//Rely:     m.Rely,
			})
		}
	}

	return &list

}

//获取我的应用列表
func (a *appStruct) GetContainerInfo(name string) (types.Container, error) {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		a.log.Error("初始化client失败", "app.getmylist", "line:36", err)
	}
	filters := filters.NewArgs()
	filters.Add("name", name)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		a.log.Error("获取docker容器失败", "app.getmylist", "line:42", err)
	}

	if len(containers) > 0 {
		return containers[0], nil
	}
	return types.Container{}, nil

}

func (a *appStruct) GetSimpleContainerInfo(name string) (types.Container, error) {
	//获取docker应用
	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return types.Container{}, err
	}
	defer cli.Close()
	filters := filters.NewArgs()
	filters.Add("name", name)
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
	command.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;DelAppConfigDir " + path)
}

func (a *appStruct) RemoveContainerById(id string) {
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

	steam := true
	for !isFinish {
		if steam {
			steam = false
			go func() {
				a.GetHardwareUsageSteam()
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

func (a *appStruct) GetHardwareUsageSteam() {

	cli, err := client2.NewClientWithOpts(client2.FromEnv)
	if err != nil {
		return
	}
	defer cli.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var lm []model2.AppListDBModel
	a.db.Table(model2.CONTAINERTABLENAME).Select("label,title,icon,container_id").Where("origin != ?", "system").Find(&lm)
	for i := 0; i < 100; i++ {
		if config.CasaOSGlobalVariables.AppChange {
			lm = []model2.AppListDBModel{}
			config.CasaOSGlobalVariables.AppChange = false
			a.db.Table(model2.CONTAINERTABLENAME).Select("label,title,icon,container_id").Where("origin != ?", "system").Find(&lm)
			dataApps := dataStats
			dataStats.Range(func(key, value interface{}) bool {
				dataStats.Delete(key)
				return true
			})
			for _, v := range lm {
				m, _ := dataApps.Load(v.ContainerId)
				if m != nil {
					dataStats.Store(v.ContainerId, m)
				}
			}
		}
		var wg sync.WaitGroup
		for _, v := range lm {
			wg.Add(1)
			go func(v model2.AppListDBModel, i int) {
				defer wg.Done()
				stats, err := cli.ContainerStats(ctx, v.ContainerId, true)
				if err != nil {
					return
				}
				decode := json.NewDecoder(stats.Body)
				var data interface{}
				if err := decode.Decode(&data); err == io.EOF {
					return
				}
				m, _ := dataStats.Load(v.ContainerId)
				dockerStats := model.DockerStatsModel{}
				if m != nil {
					dockerStats.Pre = m.(model.DockerStatsModel).Data
				}
				dockerStats.Data = data
				dockerStats.Icon = v.Icon
				if len(v.Label) > 0 {
					dockerStats.Title = v.Label
				} else {
					dockerStats.Title = v.Title
				}
				dataStats.Store(v.ContainerId, dockerStats)
				if i == 99 {
					stats.Body.Close()
				}
			}(v, i)
		}
		wg.Wait()
		isFinish = true
		time.Sleep(time.Second * 3)
	}
	isFinish = false
	cancel()
}

func NewAppService(db *gorm.DB, logger loger2.OLog) AppService {
	return &appStruct{db: db, log: logger}
}
