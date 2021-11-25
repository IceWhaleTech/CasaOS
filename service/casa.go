package service

import (
	json2 "encoding/json"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/tidwall/gjson"
)

type CasaService interface {
	GetServerList(index, size, tp, categoryId, key string) ([]model.ServerAppList, int64)
	GetServerCategoryList() []model.ServerCategoryList
	GetTaskList(size int) []model2.TaskDBModel
	GetServerAppInfo(id string) model.ServerAppList
	ShareAppFile(body []byte) string
}

type casaService struct {
}

func (o *casaService) ShareAppFile(body []byte) string {
	head := make(map[string]string)

	head["Authorization"] = GetToken()

	content := httper2.Post(config.ServerInfo.ServerApi+"/v1/community/add", body, "application/json", head)
	return content
}

func (o *casaService) GetTaskList(size int) []model2.TaskDBModel {
	head := make(map[string]string)

	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/task/list/"+strconv.Itoa(size), head)

	list := []model2.TaskDBModel{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}

func (o *casaService) GetServerList(index, size, tp, categoryId, key string) ([]model.ServerAppList, int64) {

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/list?index="+index+"&size="+size+"&type="+tp+"&category_id="+categoryId+"&key="+key, head)

	list := []model.ServerAppList{}

	count := gjson.Get(listS, "data.count").Int()
	json2.Unmarshal([]byte(gjson.Get(listS, "data.items").String()), &list)

	return list, count
}

func (o *casaService) GetServerCategoryList() []model.ServerCategoryList {

	head := make(map[string]string)
	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/category", head)

	list := []model.ServerCategoryList{}

	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}
func (o *casaService) GetServerAppInfo(id string) model.ServerAppList {

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Get(config.ServerInfo.ServerApi+"/v1/app/info/"+id, head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

	return info
}
func GetToken() string {
	t := make(chan string)
	keyName := "casa_token"

	var auth string
	if result, ok := Cache.Get(keyName); ok {
		auth, ok = result.(string)
		if ok {

			return auth
		}
	}
	go func() {
		str := httper2.Get(config.ServerInfo.ServerApi+"/token", nil)

		t <- gjson.Get(str, "data").String()
	}()
	auth = <-t

	Cache.SetDefault(keyName, auth)
	return auth
}

func NewOasisService() CasaService {
	return &casaService{}
}
