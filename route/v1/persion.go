package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func PersonTest(c *gin.Context) {

	//service.MyService.Person().GetPersionInfo("fb2333a1-72b2-4cb4-9e31-61ccaffa55b9")

	m := model.ConnectState{}
	m.CreatedAt = time.Now()
	m.From = config.ServerInfo.Token
	m.To = "fb2333a1-72b2-4cb4-9e31-61ccaffa55b9"
	m.Type = ""
	m.UUId = uuid.NewV4().String()

	//service.MyService.Person().Handshake(m)
	msg := model.MessageModel{}
	msg.Type = "connection"
	msg.Data = "fb2333a1-72b2-4cb4-9e31-61ccaffa55b9"
	msg.From = config.ServerInfo.Token
	msg.UUId = "1234567890"
	b, _ := json.Marshal(msg)
	err := service.WebSocketConn.WriteMessage(websocket.TextMessage, b)
	if err == nil {
		return
	}
}

//get other persion file
func GetPersionFile(c *gin.Context) {
	path := c.Query("path")
	persion := c.Query("persion")
	if len(path) == 0 && len(persion) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	//任务标识
	uuid := uuid.NewV4().String()

	//1.通知对方需要下载
	service.MyService.Person().GetFileDetail(uuid, path, persion)

	//2.添加数据库

	task := model2.PersionDownloadDBModel{}
	task.UUID = uuid
	task.Name = ""
	task.Length = 0
	task.Size = 0
	task.State = types.DOWNLOADAWAIT
	task.TempPath = ""
	task.Type = 0
	service.MyService.Person().AddDownloadTask(task)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}
func GetPersionDownloadList(c *gin.Context) {
	path := c.Query("path")
	persion := c.Query("persion")
	if len(path) == 0 && len(persion) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	//任务标识
	uuid := uuid.NewV4().String()

	//1.通知对方需要下载
	service.MyService.Person().GetFileDetail(uuid, path, persion)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary add friend
// @Produce  application/json
// @Accept application/json
// @Tags persion
// @Param  token formData int true "Opponent token"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /persion/edit [put]
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

// @Summary add friend
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
	//step:远程验证token是否存在
	msg := model.MessageModel{}
	msg.Type = types.PERSONADDFRIEND
	msg.To = token
	msg.Data = token
	msg.From = config.ServerInfo.Token
	msg.UUId = uuid.NewV4().String()
	b, _ := json.Marshal(msg)
	err := service.WebSocketConn.WriteMessage(websocket.TextMessage, b)
	fmt.Println(err)

	friend := model2.FriendModel{}
	friend.Token = token
	service.MyService.Friend().AddFriend(friend)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

func GetPersionDirectory(c *gin.Context) {
	path := c.Query("path")
	persion := c.Query("persion")
	if len(path) == 0 && len(persion) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	//任务标识
	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = persion
	m.Type = "directory"
	m.UUId = uuid
	result, err := service.Dial("192.168.2.225:9902", m)
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
