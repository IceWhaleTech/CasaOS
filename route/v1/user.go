package v1

import (
	json2 "encoding/json"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/model/system_model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/jwt"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/tidwall/gjson"

	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary register user
// @Router /user/register/ [post]
func PostUserRegister(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)

	// @tiger - user_name 改成 username
	username := json["user_name"]
	pwd := json["password"]

	key := c.Param("key")
	if _, ok := service.UserRegisterHash[key]; !ok {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.KEY_NOT_EXIST, Message: common_err.GetMsg(common_err.KEY_NOT_EXIST)})
		return
	}

	if len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if len(pwd) < 6 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.PWD_IS_TOO_SIMPLE, Message: common_err.GetMsg(common_err.PWD_IS_TOO_SIMPLE)})
		return
	}
	oldUser := service.MyService.User().GetUserInfoByUserName(username)
	if oldUser.Id > 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_EXIST, Message: common_err.GetMsg(common_err.USER_EXIST)})
		return
	}

	user := model2.UserDBModel{}
	user.UserName = username
	user.Password = encryption.GetMD5ByStr(config.UserInfo.PWD)
	user.Role = "admin"

	user = service.MyService.User().CreateUser(user)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR)})
		return
	}
	file.MkDir(config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id))
	delete(service.UserRegisterHash, key)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})

}

// @Summary login
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Param user_name query string true "User name"
// @Param pwd  query string true "password"
// @Success 200 {string} string "ok"
// @Router /user/login [post]
func PostUserLogin(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)

	username := json["username"]

	// @tiger - 字段命名要一直，在注册的时候如果用 password，这里也要用 password
	pwd := json["pwd"]
	//check params is empty
	if len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK,
			model.Result{
				Success: common_err.ERROR,
				Message: common_err.GetMsg(common_err.INVALID_PARAMS),
			})
		return
	}
	user := service.MyService.User().GetUserAllInfoByName(username)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	if user.Password != encryption.GetMD5ByStr(pwd) {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.PWD_INVALID, Message: common_err.GetMsg(common_err.PWD_INVALID)})
		return
	}
	token := system_model.VerifyInformation{}
	token.AccessToken = jwt.GetAccessToken(user.UserName, user.Password, user.Id)
	token.RefreshToken = jwt.GetRefreshToken(user.UserName, user.Password, user.Id)
	token.ExpiresAt = time.Now().Add(3 * time.Hour * time.Duration(1)).Unix()
	data := make(map[string]interface{}, 2)
	user.Password = ""
	data["token"] = token

	// @tiger - 不建议直接透传数据库对象，而是适配到用于 API 输出的 model 对象
	data["user"] = user

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    data,
		})
}

// @Summary edit user head
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param file formData file true "用户头像"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/avatar [put]
func PutUserAvatar(c *gin.Context) {
	id := c.GetHeader("user_id")
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}
	if len(user.Avatar) > 0 {
		os.RemoveAll(config.AppInfo.UserDataPath + "/" + id + "/" + user.Avatar)
	}
	ext := filepath.Ext(f.Filename)
	avatarPath := config.AppInfo.UserDataPath + "/" + id + "/avatar" + ext
	c.SaveUploadedFile(f, avatarPath)
	user.Avatar = avatarPath
	service.MyService.User().UpdateUser(user)
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    user,
		})
}

/**
 * @description: get user avatar by user id
 * @param {query} id string user id
 * @method: GET
 */
func GetUserAvatar(c *gin.Context) {
	id := c.Param("id")
	user := service.MyService.User().GetUserInfoById(id)

	path := "default.png"
	if user.Id > 0 {
		path = user.Avatar
	}

	// @tiger - RESTful 规范下不应该返回文件本身内容，而是返回文件的静态URL，由前端去解析
	c.File(path)
}

// @Summary edit user name
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Param old_name  query string true "Old user name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/name/:id [put]
func PutUserName(c *gin.Context) {
	id := c.GetHeader("user_id")
	json := make(map[string]string)
	c.BindJSON(&json)
	userName := json["user_name"]
	if len(userName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR)})
		return
	}
	user := service.MyService.User().GetUserInfoById(id)

	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	user.UserName = userName
	service.MyService.User().UpdateUser(user)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: user})
}

// @Summary edit user password
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/password/:id [put]
func PutUserPwd(c *gin.Context) {
	//id := c.GetHeader("user_id")
	json := make(map[string]string)
	c.BindJSON(&json)
	oldPwd := json["old_pwd"]
	pwd := json["pwd"]
	id := json["user_id"]
	if len(oldPwd) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserAllInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	if user.Password != encryption.GetMD5ByStr(oldPwd) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.PWD_INVALID_OLD, Message: common_err.GetMsg(common_err.PWD_INVALID_OLD)})
		return
	}
	user.Password = encryption.GetMD5ByStr(pwd)
	service.MyService.User().UpdateUserPassword(user)
	user.Password = ""
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: user})
}

// @Summary edit user nick
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Param nick_name query string false "nick name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/nick [put]
func PutUserNick(c *gin.Context) {

	id := c.GetHeader("user_id")
	json := make(map[string]string)
	c.BindJSON(&json)
	nickName := json["nick_name"]
	if len(nickName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	user.NickName = nickName
	service.MyService.User().UpdateUser(user)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: user})
}

// @Summary edit user description
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param description formData string false "Description"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/desc [put]
func PutUserDesc(c *gin.Context) {
	id := c.GetHeader("user_id")
	json := make(map[string]string)
	c.BindJSON(&json)
	desc := json["description"]
	if len(desc) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	user.Description = desc

	service.MyService.User().UpdateUser(user)

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: user})
}

// @Summary get user info
// @Produce  application/json
// @Accept  application/json
// @Tags user
// @Success 200 {string} string "ok"
// @Router /user/info/:id [get]
func GetUserInfo(c *gin.Context) {
	id := c.GetHeader("user_id")
	user := service.MyService.User().GetUserInfoById(id)

	//*****
	// @tiger - 应该和 PostUserLogin 中的 user 对象一致。而不是重构一系列字段。
	var u = make(map[string]string, 5)
	u["user_name"] = user.UserName // 改成 username
	u["head"] = user.Avatar        // 应该和 /v1/user/avatar/:id 一致，改成 avatar
	u["email"] = user.Email
	u["description"] = user.NickName
	u["nick_name"] = user.NickName  // 改成 nickname
	u["id"] = strconv.Itoa(user.Id) // (nice-to-have) 最佳实践是用随机字符来代表 ID。顺序数字有可预测性

	//**

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    u,
		})
}

// @Summary get user info
// @Produce  application/json
// @Accept  application/json
// @Tags user
// @Success 200 {string} string "ok"
// @Router /user/info [get]
func GetUserInfoByUserName(c *gin.Context) {
	json := make(map[string]string)

	// @tiger 当前这个设计的问题是：GET 不应该同时接收 request body。
	//        GET 方法应该只接收 URL 参数
	c.BindJSON(&json)
	userName := json["user_name"]
	if len(userName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserInfoByUserName(userName)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	//**

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    user,
		})
}

/**
 * @description: get all usernames
 * @method:GET
 * @router:/user/all/name
 */
func GetUserAllUserName(c *gin.Context) {
	users := service.MyService.User().GetAllUserName()
	names := []string{}
	for _, v := range users {
		names = append(names, v.UserName)
	}
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    names,
		})
}

/**
 * @description:get custom file by user
 * @param {path} name string "file name"
 * @method: GET
 * @router: /user/custom/:key
 */
func GetUserCustomConf(c *gin.Context) {
	name := c.Param("key")
	if len(name) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	id := c.GetHeader("user_id")

	user := service.MyService.User().GetUserInfoById(id)
	//	user := service.MyService.User().GetUserInfoByUserName(userName)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	filePath := config.AppInfo.UserDataPath + "/" + id + "/" + name + ".json"

	data := file.ReadFullFile(filePath)
	if !gjson.ValidBytes(data) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: string(data)})
		return
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: json2.RawMessage(string(data))})
}

/**
 * @description:create or update custom conf by user
 * @param {path} name string "file name"
 * @method:POST
 * @router:/user/custom/:key
 */
func PostUserCustomConf(c *gin.Context) {
	name := c.Param("key")
	if len(name) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	id := c.GetHeader("user_id")
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	data, _ := ioutil.ReadAll(c.Request.Body)
	filePath := config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id)

	file.WriteToPath(data, filePath, name+".json")
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: json2.RawMessage(string(data))})
}

/**
 * @description: delete user custom config
 * @param {path} key string
 * @method:delete
 * @router:/user/custom/:key
 */
func DeleteUserCustomConf(c *gin.Context) {
	name := c.Param("key")
	if len(name) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	id := c.GetHeader("user_id")
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	filePath := config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id) + "/" + name + ".json"
	os.Remove(filePath) // @tiger - 这里万一无法实际删除，后面仍然有可能返回成功
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

/**
 * @description:
 * @param {path} id string "user id"
 * @method:DELETE
 * @router:/user/delete/:id
 */
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	service.MyService.User().DeleteUserById(id)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: id})
}

/**
 * @description:update user image
 * @method:POST
 * @router:/user/file/image/:key
 */
func PostUserFileImage(c *gin.Context) {
	id := c.GetHeader("user_id")
	json := make(map[string]string)
	c.BindJSON(&json)

	path := json["path"]
	key := c.Param("key")
	if len(path) == 0 || len(key) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.FILE_DOES_NOT_EXIST, Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST)})
		return
	}

	_, err := file.GetImageExt(path)

	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.NOT_IMAGE, Message: common_err.GetMsg(common_err.NOT_IMAGE)})
		return
	}

	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	fstat, _ := os.Stat(path)
	if fstat.Size() > 10<<20 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.IMAGE_TOO_LARGE, Message: common_err.GetMsg(common_err.IMAGE_TOO_LARGE)})
		return
	}
	ext := file.GetExt(path)
	filePath := config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id) + "/" + key + ext
	file.CopySingleFile(path, filePath, "overwrite")

	data := make(map[string]string, 3)
	data["path"] = filePath
	data["file_name"] = key + ext
	data["online_path"] = "/v1/user/image?path=" + filePath
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

/**
 * @description:create or update user's custom image
 * @param {formData} file file "a file to be uploaded"
 * @param {path} key string "file name"
 * @method:POST
 * @router:/user/upload/image/:key
 */
func PostUserUploadImage(c *gin.Context) {
	id := c.GetHeader("user_id")
	f, err := c.FormFile("file")
	key := c.Param("key")
	if len(key) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
		return
	}

	_, err = file.GetImageExtByName(f.Filename)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: common_err.NOT_IMAGE, Message: common_err.GetMsg(common_err.NOT_IMAGE)})
		return
	}
	ext := filepath.Ext(f.Filename)
	user := service.MyService.User().GetUserInfoById(id)

	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	path := config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id) + "/" + key + ext
	c.SaveUploadedFile(f, path)
	data := make(map[string]string, 3)
	data["path"] = path
	data["file_name"] = key + ext
	data["online_path"] = "/v1/user/image?path=" + path
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

/**
 * @description: get current user's image
 * @method:GET
 * @router:/user/image/:id
 */
func GetUserImage(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(http.StatusNotFound, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusNotFound, model.Result{Success: common_err.FILE_DOES_NOT_EXIST, Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST)})
		return
	}
	if !strings.Contains(filePath, config.AppInfo.UserDataPath) {
		c.JSON(http.StatusNotFound, model.Result{Success: common_err.INSUFFICIENT_PERMISSIONS, Message: common_err.GetMsg(common_err.INSUFFICIENT_PERMISSIONS)})
		return
	}

	fileTmp, _ := os.Open(filePath)
	defer fileTmp.Close()

	fileName := path.Base(filePath)

	// @tiger - RESTful 规范下不应该返回文件本身内容，而是返回文件的静态URL，由前端去解析
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url2.PathEscape(fileName))
	c.File(filePath)
}
func DeleteUserImage(c *gin.Context) {
	id := c.GetHeader("user_id")
	path := c.Query("path")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserInfoById(id)
	if user.Id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.USER_NOT_EXIST, Message: common_err.GetMsg(common_err.USER_NOT_EXIST)})
		return
	}
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.FILE_DOES_NOT_EXIST, Message: common_err.GetMsg(common_err.FILE_DOES_NOT_EXIST)})
		return
	}
	if !strings.Contains(path, config.AppInfo.UserDataPath+"/"+strconv.Itoa(user.Id)) {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INSUFFICIENT_PERMISSIONS, Message: common_err.GetMsg(common_err.INSUFFICIENT_PERMISSIONS)})
		return
	}
	os.Remove(path)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

////refresh token
// func PostUserRefreshToken(c *gin.Context) {
// 	json := make(map[string]string)
// 	c.BindJSON(&json)
// 	refresh := json["refresh_token"]
// 	claims, err := jwt.ParseToken(refresh)
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.VERIFICATION_FAILURE), Data: err.Error()})
// 		return
// 	}
// 	if claims.VerifyExpiresAt(time.Now(), true) || claims.VerifyIssuer("refresh", true) {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.VERIFICATION_FAILURE, Message: common_err.GetMsg(common_err.VERIFICATION_FAILURE)})
// 		return
// 	}
// 	newToken := jwt.GetAccessToken(claims.UserName, claims.PassWord, claims.Id)
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.ERROR, Message: common_err.GetMsg(common_err.ERROR), Data: err.Error()})
// 		return
// 	}
// 	verifyInfo := system_model.VerifyInformation{}
// 	verifyInfo.AccessToken = newToken
// 	verifyInfo.RefreshToken = jwt.GetRefreshToken(claims.UserName, claims.PassWord, claims.Id)
// 	verifyInfo.ExpiresAt = time.Now().Add(3 * time.Hour * time.Duration(1)).Unix()

// 	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: verifyInfo})

// }
