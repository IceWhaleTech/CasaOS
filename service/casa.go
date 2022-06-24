package service

import (
	"encoding/json"
	json2 "encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type CasaService interface {
	GetServerList(index, size, tp, categoryId, key string) model.ServerAppListCollection
	GetServerCategoryList() []model.CategoryList
	GetServerAppInfo(id, t string, language string) model.ServerAppList
	ShareAppFile(body []byte) string
	PushHeart(id, t string, language string)

	PushConnectionStatus(uuid, err string, from, to, event string)
	PushUserInfo()
	GetUserInfoByShareId(shareId string) model.UserInfo
	GetPersonPublic() (list []model.FriendsModel)
	GetCasaosVersion() model.Version
	AsyncGetServerList() (collection model.ServerAppListCollection)
	AsyncGetServerCategoryList() []model.CategoryList
}

type casaService struct {
}

func (o *casaService) ShareAppFile(body []byte) string {
	head := make(map[string]string)

	head["Authorization"] = GetToken()

	content := httper2.Post(config.ServerInfo.ServerApi+"/v1/community/add", body, "application/json", head)
	return content
}

func (o *casaService) GetServerList(index, size, tp, categoryId, key string) model.ServerAppListCollection {

	keyName := fmt.Sprintf("list_%s_%s_%s_%s_%s", index, size, tp, categoryId, "en")
	collection := model.ServerAppListCollection{}
	if result, ok := Cache.Get(keyName); ok {
		res, ok := result.(string)
		if ok {
			json2.Unmarshal([]byte(res), &collection)
			return collection
		}
	}

	collectionStr := file.ReadFullFile(config.AppInfo.DBPath + "/app_list.json")

	err := json2.Unmarshal(collectionStr, &collection)
	if err != nil {
		loger.Error("marshal error", zap.Any("err", err), zap.Any("content", string(collectionStr)))
		collection = o.AsyncGetServerList()
	}

	go o.AsyncGetServerList()

	if categoryId != "0" {
		categoryInt, _ := strconv.Atoi(categoryId)
		nList := []model.ServerAppList{}
		for _, v := range collection.List {
			if v.CategoryId == categoryInt {
				nList = append(nList, v)
			}
		}
		collection.List = nList
		nCommunity := []model.ServerAppList{}
		for _, v := range collection.Community {
			if v.CategoryId == categoryInt {
				nCommunity = append(nCommunity, v)
			}
		}
		collection.Community = nCommunity
	}
	if tp != "name" {
		if tp == "new" {
			sort.Slice(collection.List, func(i, j int) bool {
				return collection.List[i].CreatedAt.After(collection.List[j].CreatedAt)
			})
			sort.Slice(collection.Community, func(i, j int) bool {
				return collection.Community[i].CreatedAt.After(collection.Community[j].CreatedAt)
			})
		} else if tp == "rank" {
			sort.Slice(collection.List, func(i, j int) bool {
				return collection.List[i].QueryCount > collection.List[j].QueryCount
			})
			sort.Slice(collection.Community, func(i, j int) bool {
				return collection.Community[i].QueryCount > collection.Community[j].QueryCount
			})
		}
	}
	sizeInt, _ := strconv.Atoi(size)

	if index != "1" {
		indexInt, _ := strconv.Atoi(index)
		collection.List = collection.List[(indexInt-1)*sizeInt : indexInt*sizeInt]
		collection.Community = collection.Community[(indexInt-1)*sizeInt : indexInt*sizeInt]
	} else {
		if len(collection.List) > sizeInt {
			collection.List = collection.List[:sizeInt]
		}
		if len(collection.Community) > sizeInt {
			collection.Community = collection.Community[:sizeInt]
		}
	}

	if len(collection.List) > 0 {
		by, _ := json.Marshal(collection)
		Cache.Set(keyName, string(by), time.Minute*10)
	}

	return collection

}

func (o *casaService) AsyncGetServerList() (collection model.ServerAppListCollection) {

	results := file.ReadFullFile(config.AppInfo.DBPath + "/app_list.json")
	err := json2.Unmarshal(results, &collection)
	if err != nil {
		loger.Error("marshal error", zap.Any("err", err), zap.Any("content", string(results)))
	}

	if collection.Version == o.GetCasaosVersion().Version {
		return collection
	}

	head := make(map[string]string)

	head["Authorization"] = GetToken()

	listS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/newlist?index=1&size=1000&rank=name&category_id=0&key=&language=en", head)
	listModel := []model.ServerAppList{}
	communityModel := []model.ServerAppList{}
	recommendModel := []model.ServerAppList{}
	json2.Unmarshal([]byte(gjson.Get(listS, "data.list").String()), &listModel)
	json2.Unmarshal([]byte(gjson.Get(listS, "data.recommend").String()), &recommendModel)
	json2.Unmarshal([]byte(gjson.Get(listS, "data.community").String()), &communityModel)

	if len(listModel) > 0 {
		collection.Community = communityModel
		collection.List = listModel
		collection.Recommend = recommendModel
		collection.Version = o.GetCasaosVersion().Version
		by, err := json.Marshal(collection)
		if err != nil {
			loger.Error("marshal error", zap.Any("err", err))
		}
		file.WriteToPath(by, config.AppInfo.DBPath, "app_list.json")
	}
	return
}

// func (o *casaService) GetServerCategoryList() (list []model.ServerCategoryList) {

// 	keyName := fmt.Sprintf("category_list")
// 	if result, ok := Cache.Get(keyName); ok {
// 		res, ok := result.(string)
// 		if ok {
// 			json2.Unmarshal([]byte(gjson.Get(res, "data").String()), &list)
// 			return list
// 		}
// 	}

// 	head := make(map[string]string)
// 	head["Authorization"] = GetToken()

// 	listS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/category", head)

// 	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &list)
// 	if len(list) > 0 {
// 		Cache.Set(keyName, listS, time.Hour*24)
// 	}
// 	return list
// }

func (o *casaService) GetServerCategoryList() (list []model.CategoryList) {
	category := model.ServerCategoryList{}
	results := file.ReadFullFile(config.AppInfo.DBPath + "/app_category.json")
	err := json2.Unmarshal(results, &category)
	if err != nil {
		loger.Error("marshal error", zap.Any("err", err), zap.Any("content", string(results)))
		return o.AsyncGetServerCategoryList()
	}
	go o.AsyncGetServerCategoryList()
	return category.Item
}

func (o *casaService) AsyncGetServerCategoryList() []model.CategoryList {
	list := model.ServerCategoryList{}
	results := file.ReadFullFile(config.AppInfo.DBPath + "/app_category.json")
	err := json2.Unmarshal(results, &list)
	if err != nil {
		loger.Error("marshal error", zap.Any("err", err), zap.Any("content", string(results)))
	}

	if list.Version == o.GetCasaosVersion().Version {
		return nil
	}
	item := []model.CategoryList{}
	head := make(map[string]string)
	head["Authorization"] = GetToken()
	listS := httper2.Get(config.ServerInfo.ServerApi+"/v2/app/category", head)
	json2.Unmarshal([]byte(gjson.Get(listS, "data").String()), &item)
	if len(item) > 0 {
		list.Version = o.GetCasaosVersion().Version
		list.Item = item
		by, err := json.Marshal(list)
		if err != nil {
			loger.Error("marshal error", zap.Any("err", err))
		}
		file.WriteToPath(by, config.AppInfo.DBPath, "app_category.json")
	}
	return item
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

/**
 * @description: get remote version
 * @return {model.Version}
 */
func (o *casaService) GetCasaosVersion() model.Version {
	keyName := "casa_version"
	var dataStr string
	var version model.Version
	if result, ok := Cache.Get(keyName); ok {
		dataStr, ok = result.(string)
		if ok {
			data := gjson.Get(dataStr, "data")
			json2.Unmarshal([]byte(data.String()), &version)
			return version
		}
	}

	v := httper.OasisGet(config.ServerInfo.ServerApi + "/v1/sys/version")
	data := gjson.Get(v, "data")
	json2.Unmarshal([]byte(data.String()), &version)

	if len(version.Version) > 0 {
		Cache.Set(keyName, v, time.Minute*20)
	}

	return version
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
