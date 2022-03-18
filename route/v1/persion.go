package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func PersonTest(c *gin.Context) {

	token := c.Query("token")
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
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary retry download file
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  uui path string true "download uuid"
// @Param  path query string true "file path"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/refile/{uuid} [get]
func GetPersionReFile(c *gin.Context) {

	path := c.Query("path")
	uuid := c.Param("uuid")

	if len(path) == 0 && len(uuid) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	task := service.MyService.Download().GetDownloadById(uuid)
	if reflect.DeepEqual(task, model2.PersionDownloadDBModel{}) {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSION_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSION_REMOTE_ERROR)})
		return
	}
	token := task.From
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSION_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSION_REMOTE_ERROR)})
		return
	}

	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONDOWNLOAD
	m.UUId = uuid
	go service.Dial(m, false)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary download file
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token query string true "opponent token"
// @Param  path query string true "file path"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/file [get]
func GetPersionFile(c *gin.Context) {

	path := c.Query("path")
	token := c.Query("token")
	if len(path) == 0 && len(token) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSION_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSION_REMOTE_ERROR)})
		return
	}
	// task id
	uuid := uuid.NewV4().String()

	task := model2.PersionDownloadDBModel{}
	task.UUID = uuid
	task.Name = ""
	task.Length = 0
	task.From = token
	task.Size = 0
	task.State = types.DOWNLOADAWAIT
	task.Type = 0
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
// @Tags persion
// @Param  uuid path string true "download uuid"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/file/{uuid} [delete]
func DeletePersionDownloadFile(c *gin.Context) {

	id := c.Param("uuid")
	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	service.MyService.Download().DelDownload(id)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary get file download list
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  state query int true "wait:1,downloading:1,pause:2,finish:3,error:4" Enums(0,1,2,4)
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/list [get]
func GetPersionDownloadList(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	list := service.MyService.Download().GetDownloadListByState(state)
	//if it is  downloading, it need to add 'already'
	if state == strconv.Itoa(types.DOWNLOADING) {
		for i := 0; i < len(list); i++ {
			tempDir := config.AppInfo.RootPath + "/temp" + "/" + list[i].UUID
			files, err := ioutil.ReadDir(tempDir)
			if err == nil {
				list[i].Already = len(files)
			}
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
}

// @Summary edit friend's nick
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param token path string true "token"
// @Param nick formData string true "nick name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/nick/{token} [put]
func PutPersionNick(c *gin.Context) {
	token := c.Param("token")
	nick := c.PostForm("nick")
	if len(token) == 0 || len(nick) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token
	friend.NickName = nick
	service.MyService.Friend().EditFriendNick(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary get friend list
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/users [get]
func GetPersionFriend(c *gin.Context) {
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
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/user [post]
func PostAddPersionFriend(c *gin.Context) {
	token := c.PostForm("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	msg := model.MessageModel{}
	msg.Type = types.PERSONCONNECTION
	msg.Data = token
	msg.From = config.ServerInfo.Token
	msg.To = token
	msg.UUId = uuid.NewV4().String()

	go service.Dial(msg, true)

	friend := model2.FriendModel{}
	friend.Token = token
	service.MyService.Friend().AddFriend(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary get directory list
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token query string true "Opponent token"
// @Param  path query string true "dir path"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/directory [get]
func GetPersionDirectory(c *gin.Context) {
	path := c.Query("path")
	token := c.Query("token")
	if len(path) == 0 && len(token) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PERSION_REMOTE_ERROR, Message: oasis_err2.GetMsg(oasis_err2.PERSION_REMOTE_ERROR)})
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
