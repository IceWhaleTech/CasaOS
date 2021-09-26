package service

import (
	"context"
	json2 "encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	client2 "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"oasis/model"
	"oasis/pkg/config"
	"oasis/pkg/docker"
	"oasis/pkg/utils/command"
	httper2 "oasis/pkg/utils/httper"
	loger2 "oasis/pkg/utils/loger"
	model2 "oasis/service/model"
	"strings"
	"time"
)

type AppService interface {
	GetMyList(index, size int, position bool) *[]model2.MyAppList
	SaveContainer(m model2.AppListDBModel)
	GetServerAppInfo(id string) model.ServerAppList
	GetUninstallInfo(id string) model2.AppListDBModel
	RemoveContainerById(id string)
	GetContainerInfo(name string) (types.Container, error)
	GetAppDBInfo(id string) model2.AppListDBModel
	UpdateApp(m model2.AppListDBModel)
	GetSimpleContainerInfo(name string) (types.Container, error)
	DelAppConfigDir(id string)
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
	a.db.Table(model2.CONTAINERTABLENAME).Select("title,icon,port_map,`index`,container_id,position,label,slogan").Find(&lm)

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

		if lMap[container.ID] != nil {
			var m model2.AppListDBModel
			m = lMap[container.ID].(model2.AppListDBModel)
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
				Slogan:   m.Slogan,
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

func (a *appStruct) GetServerAppInfo(id string) model.ServerAppList {

	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	infoS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/info/"+id, head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

	return info
}

//创建容器成功后保存容器
func (a *appStruct) SaveContainer(m model2.AppListDBModel) {
	a.db.Table(model2.CONTAINERTABLENAME).Create(&m)
}

func (a *appStruct) UpdateApp(m model2.AppListDBModel) {
	a.db.Table(model2.CONTAINERTABLENAME).Save(&m)
}

func (a *appStruct) DelAppConfigDir(id string) {
	command.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;DelAppConfigDir " + docker.GetDir(id, "/config"))
}

func (a *appStruct) RemoveContainerById(id string) {
	a.db.Table(model2.CONTAINERTABLENAME).Where("custom_id = ?", id).Delete(&model2.AppListDBModel{})
}

func NewAppService(db *gorm.DB, logger loger2.OLog) AppService {
	return &appStruct{db: db, log: logger}
}
