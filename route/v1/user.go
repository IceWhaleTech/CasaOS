package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/model/system_model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/jwt"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"

	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary register user
// @Router /user/register/ [post]
func PostUserRegister(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)
	username := json["user_name"]
	pwd := json["password"]
	key := c.GetHeader("key")
	//TODO:检查hash
	fmt.Println(key)

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
	//TODO:创建文件夹
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
func Login(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)

	username := json["user_name"]
	pwd := json["password"]
	//check params is empty
	if len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK,
			model.Result{
				Success: common_err.ERROR,
				Message: common_err.GetMsg(common_err.INVALID_PARAMS),
			})
		return
	}
	user := service.MyService.User().GetUserInfoByUserName(username)
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
	user.Password = ""
	token := system_model.VerifyInformation{}
	token.AccessToken = jwt.GetAccessToken(user.UserName, user.Password)
	token.RefreshToken = jwt.GetRefreshToken(user.UserName, user.Password)
	token.ExpiresAt = time.Now().Add(3 * time.Hour * time.Duration(1)).Format("2006-01-02 15:04:05")
	data := make(map[string]interface{}, 2)
	data["token"] = token
	data["user"] = user

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    data,
		})
}

// // @Summary edit user head
// // @Produce  application/json
// // @Accept multipart/form-data
// // @Tags user
// // @Param file formData file true "用户头像"
// // @Security ApiKeyAuth
// // @Success 200 {string} string "ok"
// // @Router /user/changhead [post]
// func PostUserHead(c *gin.Context) {
// 	file, _, _ := c.Request.FormFile("file")
// 	user_service.UpLoadFile(file, config.UserInfo.Head)
// 	c.JSON(http.StatusOK,
// 		model.Result{
// 			Success: common_err.SUCCESS,
// 			Message: common_err.GetMsg(common_err.SUCCESS),
// 			Data:    config.UserInfo.Head,
// 		})
// }

// @Summary edit user name
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Param old_name  query string true "Old user name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/name/:id [put]
func PutUserName(c *gin.Context) {
	id := c.Param("id")
	json := make(map[string]string)
	c.BindJSON(&json)
	userName := json["user_name"]
	if len(userName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
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
	id := c.Param("id")
	json := make(map[string]string)
	c.BindJSON(&json)
	oldPwd := json["old_password"]
	pwd := json["password"]
	if len(oldPwd) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	user := service.MyService.User().GetUserInfoById(id)
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
	service.MyService.User().UpdateUser(user)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: user})
}

// // @Summary edit user info
// // @Produce  application/json
// // @Accept multipart/form-data
// // @Tags user
// // @Param user_name formData string false "User Name"
// // @Param email formData string false "Email"
// // @Param description formData string false "Description"
// // @Param pwd formData string false "Password"
// // @Param old_pwd  formData string false "Old password"
// // @Param nick_name formData string false "nick name"
// // @Security ApiKeyAuth
// // @Success 200 {string} string "ok"
// // @Router /user/info [post]
// func PostUserChangeInfo(c *gin.Context) {
// 	username := c.PostForm("user_name")
// 	email := c.PostForm("email")
// 	description := c.PostForm("description")
// 	nickName := c.PostForm("nick_name")
// 	oldpwd := c.PostForm("old_pwd")
// 	pwd := c.PostForm("pwd")
// 	if len(pwd) > 0 && config.UserInfo.PWD != oldpwd {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.PWD_INVALID, Message: common_err.GetMsg(common_err.PWD_INVALID)})
// 		return
// 	}
// 	user_service.SetUser(username, pwd, "", email, description, nickName)
// 	data := make(map[string]string, 4)

// 	data["token"] = jwt2.GetToken(username, pwd)
// 	data["user_name"] = username
// 	data["head"] = config.UserInfo.Head
// 	data["nick_name"] = config.UserInfo.NickName
// 	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
// }

// @Summary edit user nick
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param nick_name formData string false "nick name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/nick [put]

func PutUserNick(c *gin.Context) {

	id := c.Param("id")
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

	go service.MyService.Casa().PushUserInfo()
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
	id := c.Param("id")
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

// // @Summary Modify user person information (Initialization use)
// // @Produce  application/json
// // @Accept multipart/form-data
// // @Tags user
// // @Param nick_name formData string false "user nick name"
// // @Param description formData string false "Description"
// // @Security ApiKeyAuth
// // @Success 200 {string} string "ok"
// // @Router /user/person/info [post]
// func PostUserPersonInfo(c *gin.Context) {
// 	desc := c.PostForm("description")
// 	nickName := c.PostForm("nick_name")
// 	if len(desc) == 0 || len(nickName) == 0 {
// 		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
// 		return
// 	}
// 	user_service.SetUser("", "", "", "", desc, nickName)
// 	data := make(map[string]string, 2)
// 	data["description"] = config.UserInfo.Description
// 	data["nick_name"] = config.UserInfo.NickName
// 	go service.MyService.Casa().PushUserInfo()
// 	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
// }

// @Summary get user info
// @Produce  application/json
// @Accept  application/json
// @Tags user
// @Success 200 {string} string "ok"
// @Router /user/info/:id [get]
func GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	user := service.MyService.User().GetUserInfoById(id)
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    user,
		})
}

// @Summary Get my shareId
// @Produce  application/json
// @Accept application/json
// @Tags user
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/shareid [get]
func GetUserShareID(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: config.ServerInfo.Token})
}

// @Summary get user info
// @Produce  application/json
// @Accept  application/json
// @Tags user
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
