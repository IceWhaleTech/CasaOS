package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"oasis/model"
	"oasis/pkg/config"
	jwt2 "oasis/pkg/utils/jwt"
	oasis_err2 "oasis/pkg/utils/oasis_err"
	"oasis/service"
)

var user_service service.UserService

func init() {
	user_service = service.NewUserService()
}

// @Summary 设置用户名和密码
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param username formData string true "User name"
// @Param pwd  formData string true "password"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/setusernamepwd [post]
func Set_Name_Pwd(c *gin.Context) {
	//json := make(map[string]string)
	//c.BindJSON(&json)
	username := c.PostForm("username")
	pwd := c.PostForm("pwd")
	//老用户名是否存在即新用户名和密码的验证
	if len(config.UserInfo.UserName) > 0 || len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK,
			model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	//开始设置
	err := user_service.SetUser(username, pwd, "", "", "")
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: fmt.Sprintf("%v", err)})
		return
	} else {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
		return
	}
}

// @Summary 登录
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param username formData string true "User name"
// @Param pwd  formData string true "password"
// @Success 200 {string} string "ok"
// @Router /user/login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("pwd")
	//检查参数是否正确
	if len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err2.ERROR,
				Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
			})
		return
	}

	//if config.UserInfo.UserName == username && config.UserInfo.PWD == pwd {
	if username == "admin" && pwd == "admin" {
		token := jwt2.GetToken(username, pwd)
		//user_service.SetUser("", "", token, "", "")
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err2.SUCCESS,
				Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
				Data:    token,
			})
		return
	}
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err2.ERROR,
			Message: oasis_err2.GetMsg(oasis_err2.ERROR),
		})

}

// @Summary 修改头像
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param file formData file true "用户头像"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/changhead [post]
func Up_Load_Head(c *gin.Context) {
	file, _, _ := c.Request.FormFile("file")
	user_service.UpLoadFile(file, config.UserInfo.Head)
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err2.SUCCESS,
			Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
			Data:    config.UserInfo.Head,
		})
	return
}

// @Summary 修改用户名
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param username formData string true "User name"
// @Param oldname  formData string true "Old user name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/changusername [put]
func Chang_User_Name(c *gin.Context) {
	oldname := c.PostForm("oldname")
	username := c.PostForm("username")
	if len(username) == 0 || config.UserInfo.UserName != oldname {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
		return
	}
	user_service.SetUser(username, "", "", "", "")
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
	return
}

// @Summary 修改密码
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param pwd formData string true "Password"
// @Param oldpwd  formData string true "Old password"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/changuserpwd [put]
func Chang_User_Pwd(c *gin.Context) {
	oldpwd := c.PostForm("oldpwd")
	pwd := c.PostForm("pwd")
	if len(pwd) == 0 || config.UserInfo.PWD != oldpwd {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.ERROR)})
		return
	}
	user_service.SetUser("", pwd, "", "", "")
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
	return
}

// @Summary 修改用户信息
// @Produce  application/json
// @Accept multipart/form-data
// @Tags user
// @Param username formData string false "User Name"
// @Param email formData string false "Email"
// @Param description formData string false "Description"
// @Param pwd formData string false "Password"
// @Param oldpwd  formData string false "Old password"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/changuserinfo [post]
func Chang_User_Info(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	description := c.PostForm("description")
	oldpwd := c.PostForm("oldpwd")
	pwd := c.PostForm("pwd")
	if len(pwd) > 0 && config.UserInfo.PWD != oldpwd {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.PWD_INVALID, Message: oasis_err2.GetMsg(oasis_err2.PWD_INVALID)})
		return
	}
	user_service.SetUser(username, pwd, "", email, description)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
	return
}

// @Summary 获取用户详情
// @Produce  application/json
// @Accept mapplication/json
// @Tags user
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	var u = make(map[string]string, 2)
	u["user_name"] = config.UserInfo.UserName
	u["token"] = config.UserInfo.Token
	u["head"] = config.UserInfo.Head
	u["email"] = config.UserInfo.Email
	u["description"] = config.UserInfo.Description
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err2.SUCCESS,
			Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
			Data:    u,
		})
	return
}
