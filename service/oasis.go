package service

import (
	json2 "encoding/json"
	"github.com/tidwall/gjson"
	"oasis/model"
	"oasis/pkg/config"
	httper2 "oasis/pkg/utils/httper"
	model2 "oasis/service/model"
	"strconv"
)

type OasisService interface {
	GetServerList(index, size, tp, categoryId, key string) ([]model.ServerAppList, int64)
	GetServerCategoryList() []model.ServerCategoryList
	GetTaskList(size int) []model2.TaskDBModel
}

type oasisService struct {
}

func (o *oasisService) GetTaskList(size int) []model2.TaskDBModel {
	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/task/list/"+strconv.Itoa(size), head)

	list := []model2.TaskDBModel{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}

func (o *oasisService) GetServerList(index, size, tp, categoryId, key string) ([]model.ServerAppList, int64) {

	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/list?index="+index+"&size="+size+"&type="+tp+"&category_id="+categoryId+"&key="+key, head)

	list := []model.ServerAppList{}

	count := gjson.Get(listS, "data.count").Int()
	json2.Unmarshal([]byte(gjson.Get(listS, "data.items").String()), &list)

	return list, count
}

func (o *oasisService) GetServerCategoryList() []model.ServerCategoryList {

	head := make(map[string]string)

	t := make(chan string)

	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	head["Authorization"] = <-t

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/category", head)

	list := []model.ServerCategoryList{}

	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}

func NewOasisService() OasisService {
	return &oasisService{}
}
