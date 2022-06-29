package v1

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	natType "github.com/Curtis-Milo/nat-type-identifier-go"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/ip_helper"

	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}

	task := service.MyService.Download().GetDownloadById(uid)
	if reflect.DeepEqual(task, model2.PersonDownloadDBModel{}) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_REMOTE_ERROR, Message: common_err.GetMsg(common_err.PERSON_REMOTE_ERROR)})
		return
	}
	token := task.From
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_REMOTE_ERROR, Message: common_err.GetMsg(common_err.PERSON_REMOTE_ERROR)})
		return
	}

	m := model.MessageModel{}
	m.Data = task.Path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONDOWNLOAD
	m.UUId = uid
	go service.Dial(m, false)

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if file.CheckNotExist(localPath) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DIR_NOT_EXISTS, Message: common_err.GetMsg(common_err.DIR_NOT_EXISTS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_REMOTE_ERROR, Message: common_err.GetMsg(common_err.PERSON_REMOTE_ERROR)})
		return
	}

	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_REMOTE_ERROR, Message: common_err.GetMsg(common_err.PERSON_REMOTE_ERROR)})
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
	task.Type = types.PERSONFILEDOWNLOAD
	task.LocalPath = localPath
	if service.MyService.Download().GetDownloadListByPath(task) > 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_EXIST_DOWNLOAD, Message: common_err.GetMsg(common_err.PERSON_EXIST_DOWNLOAD)})
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

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
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

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
	list := service.MyService.Download().GetDownloadListByState(state, types.PERSONFILEDOWNLOAD)
	//if it is  downloading, it need to add 'already'
	for i := 0; i < len(list); i++ {
		if list[i].State == types.DOWNLOADING {
			tempDir := config.AppInfo.TempPath + "/" + list[i].UUID
			files, err := ioutil.ReadDir(tempDir)
			if err == nil {
				list[i].Already = len(files)
			}
		}
		list[i].Duration = time.Now().Unix() - list[i].Created
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token
	friend.Mark = mark
	service.MyService.Friend().EditFriendMark(friend)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary edit friend's
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param write formData bool true "write"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/write/{shareid} [put]
func PutPersonWrite(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	write, _ := strconv.ParseBool(c.PostForm("write"))
	friend := model2.FriendModel{}
	friend.Token = token
	friend.Write = write
	service.MyService.Friend().EditFriendMark(friend)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary image thumbnail
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Param write formData bool true "write"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/image/thumbnail/{shareid} [get]
func GetPersonImageThumbnail(c *gin.Context) {
	token := c.Param("shareid")
	path := c.Query("path")
	_, err := uuid.FromString(token)
	if err != nil || len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONIMAGETHUMBNAIL
	m.UUId = uuid

	img, err := service.Dial(m, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}

	//	var buf bytes.Buffer
	//err = gob.NewEncoder(&buf).Encode(img.Data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(img.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}

	imageBuffer, _ := base64.StdEncoding.DecodeString(img.Data.(string))
	c.Writer.WriteString(string(imageBuffer))
	//	c.String(http.StatusOK, "data:image/"+filesuffix+";base64,"+img.Data.(string))
	//c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: img.Data.(string)})
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
			if ip_helper.HasLocalIP(net.ParseIP(strings.Split(v, ":")[0])) {
				list[i].LocalIP = strings.Split(v, ":")[0]
			}
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

// @Summary network type detection
// @Produce application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/detection [get]
func GetPersonDetection(c *gin.Context) {
	// - Blocked
	// - Open Internet
	// - Full Cone
	// - Symmetric UDP Firewall
	// - Restric NAT
	// - Restric Port NAT
	// - Symmetric NAT

	result, err := natType.GetDeterminedNatType(true, 5, "stun.l.google.com")
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}
	//result := service.MyService.Person().GetPersionNetWorkTypeDetection()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: result})
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
			c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
			return
		}

		if v == config.ServerInfo.Token {
			c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_MYSELF, Message: common_err.GetMsg(common_err.PERSON_MYSELF)})
			return
		}

		udb := service.MyService.Friend().GetFriendById(model2.FriendModel{Token: v})
		if !reflect.DeepEqual(udb, model2.FriendModel{Token: v}) {
			c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_EXIST_FRIEND, Message: common_err.GetMsg(common_err.PERSON_EXIST_FRIEND)})
			return
		}

		user := service.MyService.Casa().GetUserInfoByShareId(v)
		if reflect.DeepEqual(user, model.UserInfo{}) {
			c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_NOT_EXIST_USER, Message: common_err.GetMsg(common_err.PERSON_NOT_EXIST_USER)})
			return
		}

		message := model.MessageModel{}
		message.Type = types.PERSONCONNECTION
		message.Data = v
		message.From = config.ServerInfo.Token
		message.To = v
		message.UUId = uuid.NewV4().String()

		go service.Dial(message, true)

		msg := model.MessageModel{}
		msg.Type = types.PERSONGETIP
		msg.Data = ""
		msg.From = config.ServerInfo.Token
		msg.To = v
		msg.UUId = uuid.NewV4().String()

		service.Dial(msg, true)

		friend := model2.FriendModel{}
		friend.Token = v
		friend.Avatar = user.Avatar
		friend.Block = false
		friend.State = types.FRIENDSTATEWAIT
		friend.NickName = user.NickName
		friend.Profile = user.Desc
		friend.Version = user.Version
		service.MyService.Friend().AddFriend(friend)
	}

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if _, ok := service.UDPAddressMap[token]; !ok {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_REMOTE_ERROR, Message: common_err.GetMsg(common_err.PERSON_REMOTE_ERROR)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}
	dataModel := []model.Path{}
	if uuid == m.UUId {
		dataModelByte, _ := json.Marshal(result.Data)
		err := json.Unmarshal(dataModelByte, &dataModel)
		if err != nil {
			c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: dataModel})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if file.CheckNotExist(downPath) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.DIR_NOT_EXISTS, Message: common_err.GetMsg(common_err.DIR_NOT_EXISTS)})
		return
	}
	config.Cfg.Section("file").Key("DownloadDir").SetValue(downPath)
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	config.FileSettingInfo.DownloadDir = downPath
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary Get the download storage directory
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/down/dir [get]
func GetPersonDownDir(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: config.FileSettingInfo.DownloadDir})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}

	var list []string
	json.Unmarshal([]byte(share), &list)

	if len(list) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	for _, v := range list {
		if !file.Exists(v) {
			c.JSON(http.StatusOK, model.Result{Success: common_err.FILE_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
			return
		}
	}

	config.Cfg.Section("file").Key("ShareDir").SetValue(strings.Join(list, "|"))
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	config.FileSettingInfo.ShareDir = list
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary Get the shared directory
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/share [get]
func GetPersonShare(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: config.FileSettingInfo.ShareDir})
}

// @Summary Get the shareid
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/shareid [get]
func GetPersonShareId(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: config.ServerInfo.Token})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token
	friend.Block = block
	service.MyService.Friend().EditFriendBlock(friend)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	friend := model2.FriendModel{}
	friend.Token = token

	service.MyService.Friend().DeleteFriend(friend)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary Get public person
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/public [get]
func GetPersonPublic(c *gin.Context) {
	list := service.MyService.Casa().GetPersonPublic()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

// @Summary upload file to friend
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Param path formData string true "Destination path"
// @Param local_path formData string true "Full path of the file to be uploaded"
// @Success 200 {string} string "ok"
// @Router /person/file/{shareid} [post]
func PostPersonFile(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	path := c.PostForm("path")
	localPath := c.PostForm("local_path")
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if !file.Exists(localPath) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.FILE_DOES_NOT_EXIST, Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST)})
		return
	}
	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = path
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONUPLOAD
	m.UUId = uuid
	go service.UDPSendData(m, localPath)

	f, _ := os.Stat(localPath)

	task := model2.PersonDownloadDBModel{}
	task.UUID = uuid
	task.Name = f.Name()
	task.Length = 0
	task.From = token
	task.Path = path
	task.Size = f.Size()
	task.State = types.DOWNLOADFINISHED
	task.Created = time.Now().Unix()
	task.Type = types.PERSONFILEUPLOAD
	task.LocalPath = localPath
	if service.MyService.Download().GetDownloadListByPath(task) > 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PERSON_EXIST_DOWNLOAD, Message: common_err.GetMsg(common_err.PERSON_EXIST_DOWNLOAD)})
		return
	}
	service.MyService.Download().AddDownloadTask(task)

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// @Summary agree add friend
// @Produce  application/json
// @Accept application/json
// @Tags person
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /person/friend/{shareid} [put]
func PutPersonAgreeFriend(c *gin.Context) {
	token := c.Param("shareid")
	_, err := uuid.FromString(token)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}

	user := service.MyService.Friend().GetFriendById(model2.FriendModel{Token: token})

	if user.State != types.FRIENDSTATEREQUEST {
		c.JSON(http.StatusOK, model.Result{Success: common_err.COMMAND_ERROR_INVALID_OPERATION, Message: common_err.GetMsg(common_err.COMMAND_ERROR_INVALID_OPERATION)})
		return
	}
	service.MyService.Friend().AgreeFrined(user.Token)

	uuid := uuid.NewV4().String()
	m := model.MessageModel{}
	m.Data = ""
	m.From = config.ServerInfo.Token
	m.To = token
	m.Type = types.PERSONAGREEFRIEND
	m.UUId = uuid
	go service.Dial(m, true)

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

// // @Summary upload file
// // @Produce  application/json
// // @Accept  multipart/form-data
// // @Tags person
// // @Security ApiKeyAuth
// // @Param path formData string false "file path"
// // @Param file formData file true "file"
// // @Success 200 {string} string "ok"
// // @Router /person/upload/{shareid} [get]
// func GetPersonFileUpload(c *gin.Context) {

// 	token := c.Param("shareid")
// 	_, err := uuid.FromString(token)
// 	path := c.Query("path")
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
// 		return
// 	}

// 	relative := c.Query("relativePath")
// 	fileName := c.Query("filename")
// 	chunkNumber := c.Query("chunkNumber")
// 	totalChunks, _ := strconv.Atoi(c.DefaultQuery("totalChunks", "0"))
// 	dirPath := ""
// 	hash := file.GetHashByContent([]byte(fileName))
// 	tempDir := "/casaOS/temp/" + hash + strconv.Itoa(totalChunks) + "/"
// 	if fileName != relative {
// 		dirPath = strings.TrimSuffix(relative, fileName)
// 		tempDir += dirPath
// 		file.MkDir(path + "/" + dirPath)
// 	}
// 	tempDir += chunkNumber
// 	if !file.CheckNotExist(tempDir) {
// 		c.JSON(200, model.Result{Success: 200, Message: common_err.GetMsg(common_err.FILE_ALREADY_EXISTS)})
// 		return
// 	}

// 	c.JSON(204, model.Result{Success: 204, Message: common_err.GetMsg(common_err.SUCCESS)})
// }

// // @Summary upload file
// // @Produce  application/json
// // @Accept  multipart/form-data
// // @Tags person
// // @Security ApiKeyAuth
// // @Param path formData string false "file path"
// // @Param file formData file true "file"
// // @Success 200 {string} string "ok"
// // @Router /person/upload [post]
// func PostPersonFileUpload(c *gin.Context) {
// 	token := c.Param("shareid")
// 	_, err := uuid.FromString(token)
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
// 		return
// 	}
// 	f, _, _ := c.Request.FormFile("file")
// 	relative := c.PostForm("relativePath")
// 	fileName := c.PostForm("filename")
// 	totalChunks, _ := strconv.Atoi(c.DefaultPostForm("totalChunks", "0"))
// 	chunkNumber := c.PostForm("chunkNumber")
// 	dirPath := ""
// 	path := c.PostForm("path")

// 	hash := file.GetHashByContent([]byte(fileName))

// 	if len(path) == 0 {
// 		c.JSON(common_err.INVALID_PARAMS, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
// 		return
// 	}
// 	tempDir := "/casaOS/temp/" + hash + strconv.Itoa(totalChunks) + "/"

// 	if fileName != relative {
// 		dirPath = strings.TrimSuffix(relative, fileName)
// 		tempDir += dirPath
// 		file.MkDir(path + "/" + dirPath)
// 	}

// 	path += "/" + relative

// 	if !file.CheckNotExist(tempDir + chunkNumber) {
// 		file.RMDir(tempDir + chunkNumber)
// 	}

// 	if totalChunks > 1 {
// 		file.IsNotExistMkDir(tempDir)

// 		out, _ := os.OpenFile(tempDir+chunkNumber, os.O_WRONLY|os.O_CREATE, 0644)
// 		defer out.Close()
// 		_, err := io.Copy(out, f)
// 		if err != nil {
// 			c.JSON(common_err.ERROR, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
// 			return
// 		}
// 	} else {
// 		out, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
// 		defer out.Close()
// 		_, err := io.Copy(out, f)
// 		if err != nil {
// 			c.JSON(common_err.ERROR, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
// 		return
// 	}

// 	fileNum, err := ioutil.ReadDir(tempDir)
// 	if err != nil {
// 		c.JSON(common_err.ERROR, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
// 		return
// 	}
// 	if totalChunks == len(fileNum) {
// 		file.RMDir(tempDir)
// 	}

// 	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
// }
