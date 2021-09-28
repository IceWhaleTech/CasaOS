package service

import (
	json2 "encoding/json"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	loger2 "github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"strconv"
)

type TaskService interface {
	List(desc bool) []model.TaskDBModel
	Delete(id string)
	Add(m *model.TaskDBModel)
	Update(m *model.TaskDBModel)
	Info(id string) model.TaskDBModel
	SyncTaskService()
	GetServerTasks() []model.TaskDBModel
}

type taskService struct {
	db  *gorm.DB
	log loger2.OLog
}

func (s *taskService) List(desc bool) []model.TaskDBModel {
	var list []model.TaskDBModel
	var orderBy string
	if !desc {
		orderBy = "id"
	} else {
		orderBy = "id DESC"
	}
	s.db.Order(orderBy).Where("state=?", types.TASK_STATE_UNCOMPLETE).Find(&list)
	return list
}

func (s *taskService) Delete(id string) {
	var m model.TaskDBModel
	s.db.Where("id = ?", id).Delete(&m)
}

func (s *taskService) Add(m *model.TaskDBModel) {
	s.db.Save(m)
}
func (s *taskService) Update(m *model.TaskDBModel) {
	s.db.Model(&m).Update("state", m.State)
}
func (s *taskService) taskDirService(id string) model.TaskDBModel {
	var m model.TaskDBModel
	s.db.Where("id = ?", id).First(&m)
	return m
}

func (s *taskService) Info(id string) model.TaskDBModel {
	var m model.TaskDBModel
	s.db.Where("id = ?", id).Delete(&m)
	return m
}
func (s *taskService) GetServerTasks() []model.TaskDBModel {
	var count int64
	s.db.Model(&model.TaskDBModel{}).Count(&count)
	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/task/list/0?desc=true", head)

	list := []model.TaskDBModel{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	go func(list []model.TaskDBModel) {
		for _, dbModel := range list {
			dbModel.Id = 0
			s.db.Create(&dbModel)
		}
	}(list)
	return list
}
func (s *taskService) SyncTaskService() {
	var count int64
	s.db.Model(&model.TaskDBModel{}).Count(&count)
	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/task/list/"+strconv.Itoa(int(count)), head)

	list := []model.TaskDBModel{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	go func(list []model.TaskDBModel) {
		for _, dbModel := range list {
			dbModel.Id = 0
			s.db.Create(&dbModel)
		}
	}(list)
}
func SyncTask(db *gorm.DB) {
	var count int64
	db.Model(&model.TaskDBModel{}).Count(&count)
	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/task/list/"+strconv.Itoa(int(count)), head)

	list := []model.TaskDBModel{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	//go func(list []model.TaskDBModel) {
	//	for _, dbModel := range list {
	//		dbModel.Id = 0
	//		db.Create(&dbModel)
	//	}
	//}(list)
}
func NewTaskService(db *gorm.DB, log loger2.OLog) TaskService {
	return &taskService{db: db, log: log}
}
