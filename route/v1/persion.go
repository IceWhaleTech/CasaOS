package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func PersonTest(c *gin.Context) {
	token := c.Query("token")
	_, err := uuid.FromString(token)
	fmt.Println(err)

	//service.MyService.Person().GetPersionInfo("fb2333a1-72b2-4cb4-9e31-61ccaffa55b9")

	msg := model.MessageModel{}
	msg.Type = types.PERSONHELLO
	msg.Data = ""
	msg.From = config.ServerInfo.Token
	msg.To = token
	msg.UUId = uuid.NewV4().String()

	dd, err := service.Dial(msg, true)
	if err == nil {
		fmt.Println(err)
	}
	fmt.Println(dd)
	user := service.MyService.Casa().GetUserInfoByShareId(token)
	if reflect.DeepEqual(user, model.UserInfo{}) {
		fmt.Println("空数据")
	}
	fmt.Println(user)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Retry the file that failed to download
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param  uui path string true "download uuid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/refile/{uuid} [get]
func GetPersonReFile(c *gin.Context) {

	uid := c.Param("uuid")
	_, err := uuid.FromString(uid)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	task := service.MyService.Download().GetDownloadById(uid)
	if reflect.DeepEqual(task, model2.PersonDownloadDBModel{}) {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSON_REMOTE_ERROR)})
		return
	}
	token := task.From
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSON_REMOTE_ERROR)})
		return
	}

	m := model.MessageModel{}
	m.Data = task.Path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONDOWNLOAD
	m.UUId = uid
	go service.Dial(m, false)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary download file
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param  share_id query string true "opponent share_id"
// @Param  path query string true "file path"
// @Param  file_name query string true "file name"
// @Param  local_path query string true "local_path"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/file [get]
func GetPersonFile(c *gin.Context) {

	path := c.Query("path")
	localPath := c.Query("local_path")
	token := c.Query("share_id")
	fileName := c.Query("file_name")
	_, err := uuid.FromString(token)
	if len(path) == 0 || err != nil || len(localPath) == 0 || len(fileName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if file.CheckNotExist(localPath) {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.DIR_NOT_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.DIR_NOT_EXISTS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSON_REMOTE_ERROR)})
		return
	}

	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSON_REMOTE_ERROR)})
		return
	}

	// task id
	uuid := uuid.NewV4().String()

	task := model2.PersonDownloadDBModel{}
	task.UUID = uuid
	task.Name = fileName
	task.Length = 0
	task.From = token
	task.Path = path
	task.Size = 0
	task.State = types.DOWNLOADAWAIT
	task.Created = time.Now().Unix()
	task.Type = 0
	task.LocalPath = localPath
	if service.MyService.Download().GetDownloadListByPath(task) > 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_EXIST_DOWNLOAD, Message: oasis_err2.GetMsg(oasis_err2.PERSON_EXIST_DOWNLOAD)})
		return
	}
	service.MyService.Download().AddDownloadTask(task)

	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONDOWNLOAD
	m.UUId = uuid
	go service.Dial(m, false)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary delete download file records
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param  uuid path string true "download uuid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/file/{uuid} [delete]
func DeletePersonDownloadFile(c *gin.Context) {

	id := c.Param("uuid")
	_, err := uuid.FromString(id)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	task := service.MyService.Download().GetDownloadById(id)
	if task.State == types.DOWNLOADING {
		m := model.MessageModel{}
		m.Data = ""
		m.From = config.ServerInfo.Token
		m.To = task.From
		m.Type = types.PERSONCANCEL
		m.UUId = task.UUID
		service.CancelList[task.UUID] = task.UUID
		service.Dial(m, false)
	}
	service.MyService.Download().DelDownload(id)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Get file download list
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param  state query int false "wait:0,downloading:1,pause:2,finish:3,error:4,finished:5" Enums(0,1,2,3,4,5)
// @Security ApiKeyAuth
// @Success 200 {object} []model2.PersonDownloadDBModel
// @Router /person/list [get]
func GetPersonDownloadList(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	list := service.MyService.Download().GetDownloadListByState(state)
	//if it is  downloading, it need to add 'already'
	for i := 0; i < len(list); i++ {
		if list[i].State == types.DOWNLOADING {
			tempDir := config.AppInfo.RootPath + "/temp" + "/" + list[i].UUID
			files, err := ioutil.ReadDir(tempDir)
			if err == nil {
				list[i].Already = len(files)
			}
		}
		list[i].Duration = time.Now().Unix() - list[i].Created
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
}

// @Summary edit friend's remarks
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param remarks formData string true "remarks name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/remarks/{shareid} [put]
func PutPersonRemarks(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	mark := c.PostForm("remarks")
	if err != nil || len(mark) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token
	friend.Mark = mark
	service.MyService.Friend().EditFriendMark(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary get my friend list
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {object}  []model2.FriendModel
// @Router /person/users [get]
func GetPersonFriend(c *gin.Context) {
	list := service.MyService.Friend().GetFriendList()
	for i := 0; i < len(list); i++ {
		if v, ok := service.UDPAddressMap[list[i].Token]; ok && len(v) > 0 {
			list[i].OnLine = true
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
}

// @Summary add friend
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/user/{shareids} [post]
func PostAddPersonFriend(c *gin.Context) {
	token := c.Param("shareids")
	tokenList := strings.Split(token, ",")

	for _, v := range tokenList {
		_, err := uuid.FromString(v)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
			return
		}

		if v == config.ServerInfo.Token {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_MYSELF, Message: oasis_err2.GetMsg(oasis_err2.PERSON_MYSELF)})
			return
		}

		udb := service.MyService.Friend().GetFriendById(model2.FriendModel{Token: v})
		if !reflect.DeepEqual(udb, model2.FriendModel{Token: v}) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_EXIST_FRIEND, Message: oasis_err2.GetMsg(oasis_err2.PERSON_EXIST_FRIEND)})
			return
		}

		user := service.MyService.Casa().GetUserInfoByShareId(v)
		if reflect.DeepEqual(user, model.UserInfo{}) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_NOT_EXIST_USER, Message: oasis_err2.GetMsg(oasis_err2.PERSON_NOT_EXIST_USER)})
			return
		}

		message := model.MessageModel{}
		message.Type = types.PERSONCONNECTION
		message.Data = v
		message.From = config.ServerInfo.Token
		message.To = v
		message.UUId = uuid.NewV4().String()

		go service.Dial(message, true)

		friend := model2.FriendModel{}
		friend.Token = v
		friend.Avatar = user.Avatar
		friend.Block = false
		friend.NickName = user.NickName
		friend.Profile = user.Desc
		friend.Version = user.Version
		service.MyService.Friend().AddFriend(friend)
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Get a list of directories
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param  share_id query string true "Opponent share_id"
// @Param  path query string true "dir path"
// @Security ApiKeyAuth
// @Success 200 {object}  []model.Path
// @Router /person/directory [get]
func GetPersonDirectory(c *gin.Context) {
	path := c.Query("path")
	token := c.Query("share_id")
	_, err := uuid.FromString(token)
	if len(path) == 0 || err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSON_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSON_REMOTE_ERROR)})
		return
	}
	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONDIRECTORY
	m.UUId = uuid
	result, err := service.Dial(m, false)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
		return
	}
	dataModel := []model.Path{}
	if uuid == m.UUId {
		dataModelByte, _ := json.Marshal(result.Data)
		err := json.Unmarshal(dataModelByte, &dataModel)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR), Data: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: dataModel})
}

// @Summary Modify the download storage directory
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags person
// @Security ApiKeyAuth
// @Param path formData string true "path"
// @Success 200 {string} string "ok"
// @Router /person/down/dir [post]
func PostPersonDownDir(c *gin.Context) {

	downPath := c.PostForm("path")

	if len(downPath) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if file.CheckNotExist(downPath) {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.DIR_NOT_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.DIR_NOT_EXISTS)})
		return
	}
	config.Cfg.Section("file").Key("DownloadDir").SetValue(downPath)
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	config.FileSettingInfo.DownloadDir = downPath
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Get the download storage directory
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/down/dir [get]
func GetPersonDownDir(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: config.FileSettingInfo.DownloadDir})
}

// @Summary Modify the shared directory
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags person
// @Security ApiKeyAuth
// @Param share formData string true "share"
// @Success 200 {string} string "ok"
// @Router /person/share [post]
func PostPersonShare(c *gin.Context) {

	share := c.PostForm("share")

	if len(share) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	var list []string
	json.Unmarshal([]byte(share), &list)

	if len(list) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	for _, v := range list {
		if !file.Exists(v) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.FILE_ALREADY_EXISTS, Message: oasis_err2.GetMsg(oasis_err2.FILE_ALREADY_EXISTS)})
			return
		}
	}

	config.Cfg.Section("file").Key("ShareDir").SetValue(strings.Join(list, "|"))
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	config.FileSettingInfo.ShareDir = list
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Get the shared directory
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/share [get]
func GetPersonShare(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: config.FileSettingInfo.ShareDir})
}

// @Summary Modify disabled status
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param block formData bool false "Disable or not,Default:false "
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/block/{shareid} [put]
func PutPersonBlock(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	block, _ := strconv.ParseBool(c.PostForm("block"))
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token
	friend.Block = block
	service.MyService.Friend().EditFriendBlock(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Delete my friend
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/user/{shareid} [delete]
func DeletePersonFriend(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token

	service.MyService.Friend().DeleteFriend(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary Get public person
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/public [delete]
func GetPersonPublic(c *gin.Context) {
	list := service.MyService.Casa().GetPersonPublic()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
}
