package service

import (
	"encoding/json"
	json2 "encoding/json"
	"fmt"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/tidwall/gjson"
)

type CasaService interface {
	GetServerList(index, size, tp, categoryId, key, language string) (recommend, list, community []model.ServerAppList)
	GetServerCategoryList() []model.ServerCategoryList
	GetTaskList(size int) []model2.TaskDBModel
	GetServerAppInfo(id, t string, language string) model.ServerAppList
	ShareAppFile(body []byte) string
	PushHeart(id, t string, language string)
	PushAppAnalyse(uuid, t string, name, language string)
	PushConnectionStatus(uuid, err string, from, to, event string)
	PushUserInfo()
	GetUserInfoByShareId(shareId string) model.UserInfo
	GetPersonPublic() (list []model.FriendsModel)
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

func (o *casaService) GetServerList(index, size, tp, categoryId, key, language string) (recommend, list, community []model.ServerAppList) {

	keyName := fmt.Sprintf("list_%s_%s_%s_%s_%s", index, size, tp, categoryId, language)

	if result, ok := Cache.Get(keyName); ok {
		res, ok := result.(string)
		if ok {
			json2.Unmarshal([]byte(gjson.Get(res, "data.list").String()), &list)
			json2.Unmarshal([]byte(gjson.Get(res, "data.recommend").String()), &recommend)
			json2.Unmarshal([]byte(gjson.Get(res, "data.community").String()), &community)
			return
		}
	}

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/newlist?index="+index+"&size="+size+"&rank="+tp+"&category_id="+categoryId+"&key="+key+"&language="+language, head)

	json2.Unmarshal([]byte(gjson.Get(listS, "data.list").String()), &list)
	json2.Unmarshal([]byte(gjson.Get(listS, "data.recommend").String()), &recommend)
	json2.Unmarshal([]byte(gjson.Get(listS, "data.community").String()), &community)

	if len(list) > 0 {
		Cache.SetDefault(keyName, listS)
	}
	return
}

func (o *casaService) GetServerCategoryList() []model.ServerCategoryList {

	head := make(map[string]string)
	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/category", head)

	list := []model.ServerCategoryList{}

	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}
func (o *casaService) GetServerAppInfo(id, t string, language string) model.ServerAppList {

	head := make(map[string]string)

	head["Authorization"] = GetToken()
	infoS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/info/"+id+"?t="+t+"&language="+language, head)

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

func (o *casaService) PushHeart(id, t string, language string) {

	m := model.CasaOSHeart{}
	m.UuId = id
	m.Type = t
	b, _ := json.Marshal(m)

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Post(config.ServerInfo.ServerApi+"/v1/analyse/heart", b, "application/json", head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

}

func (o *casaService) PushAppAnalyse(uuid, t string, name, language string) {

	m := model.AppAnalyse{}
	m.UUId = uuid
	m.Type = t
	m.Name = name
	m.Language = language
	b, _ := json.Marshal(m)

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Post(config.ServerInfo.ServerApi+"/v1/analyse/app", b, "application/json", head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

}
func (o *casaService) PushConnectionStatus(uuid, err string, from, to, event string) {

	m := model.ConnectionStatus{}
	m.UUId = uuid
	m.Error = err
	m.From = from
	m.To = to
	m.Event = event
	b, _ := json.Marshal(m)

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Post(config.ServerInfo.ServerApi+"/v1/analyse/connect", b, "application/json", head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

}
func (o *casaService) PushUserInfo() {
	m := model.UserInfo{}
	m.Desc = config.UserInfo.Description
	m.Avatar = config.UserInfo.Avatar
	m.NickName = config.UserInfo.NickName
	m.ShareId = config.ServerInfo.Token
	b, _ := json.Marshal(m)

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Post(config.ServerInfo.ServerApi+"/v1/user/info", b, "application/json", head)

	info := model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)

}

func (o *casaService) GetUserInfoByShareId(shareId string) model.UserInfo {

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	infoS := httper2.Get(config.ServerInfo.ServerApi+"/v1/user/info/"+shareId, head)

	info := model.UserInfo{}
	json2.Unmarshal([]byte(gjson.Get(infoS, "data").String()), &info)
	return info
}
func (o *casaService) GetPersonPublic() (list []model.FriendsModel) {
	head := make(map[string]string)

	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v1/person/public", head)

	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)

	return list
}
func NewCasaService() CasaService {
	return &casaService{}
}
