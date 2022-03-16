package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	msg.Type = "hello"
	msg.Data = ""
	msg.From = config.ServerInfo.Token
	msg.To = token
	msg.UUId = uuid.NewV4().String()

	dd, err := service.Dial("", msg)
	if err == nil {
		fmt.Println(err)
	}
	fmt.Println(dd)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary add friend
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/file/{id} [delete]
func GetPersionFile(c *gin.Context) {

	path := c.Query("path")
	token := c.Query("token")
	if len(path) == 0 && len(token) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	//任务标识
	uuid := uuid.NewV4().String()

	//2.添加数据库

	task := model2.PersionDownloadDBModel{}
	task.UUID = uuid
	task.Name = ""
	task.Length = 0
	task.Size = 0
	task.State = types.DOWNLOADAWAIT
	task.Type = 0
	service.MyService.Download().AddDownloadTask(task)

	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = "file_data"
	m.UUId = uuid
	_, err := service.Dial("192.168.2.224:9902", m)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary delete download file records
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/file/{id} [delete]
func DeletePersionDownloadFile(c *gin.Context) {

	id := c.Param("id")
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
// @Param  state query int true "wait:1,loading:1,pause:2,finish:3,error:4" Enums(0,1,2,4)
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/list [get]
func GetPersionDownloadList(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	list := service.MyService.Download().GetDownloadListByState(state)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: list})
}

// @Summary add friend
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/edit/{token} [put]
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

// @Summary get friends list
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/users [get]
func GetPersionFriend(c *gin.Context) {
	list := service.MyService.Friend().GetFriendList()
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
	msg.Type = "connection"
	msg.Data = token
	msg.From = config.ServerInfo.Token
	msg.To = token
	msg.UUId = uuid.NewV4().String()

	_, err := service.Dial("", msg)

	fmt.Println(err)

	friend := model2.FriendModel{}
	friend.Token = token
	service.MyService.Friend().AddFriend(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary get directory list
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token query int true "Opponent token"
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
	//任务标识
	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = "directory"
	m.UUId = uuid
	result, err := service.Dial(service.UDPAddressMap[token], m)
	if err != nil {
		fmt.Println(err)
	}
	dataModel := []model.Path{}
	if m.UUId == m.UUId {
		dataModelByte, _ := json.Marshal(result.Data)
		err := json.Unmarshal(dataModelByte, &dataModel)
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: dataModel})
}
