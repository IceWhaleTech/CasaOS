package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oasis/model"
	"oasis/pkg/utils/oasis_err"
	"oasis/service"
	model2 "oasis/service/model"
	"strconv"
)

// @Summary 获取列表
// @Produce  application/json
// @Accept application/json
// @Tags share
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /share/list [get]
func GetShareDirList(c *gin.Context) {
	list := service.MyService.ShareDirectory().List(true)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: list})
}

// @Summary 添加文件共享
// @Produce  application/json
// @Accept multipart/form-data
// @Tags share
// @Security ApiKeyAuth
// @Param  path formData string true "要分享的文件路径"
// @Param  name formData string true "名称"
// @Param  comment formData string true "描述"
// @Param  read_only formData bool false "是否只读"
// @Param  writeable formData bool false "是否可写"
// @Param  browseable formData bool false "是否可浏览"
// @Param  user formData string false "用户"
// @Success 200 {string} string "ok"
// @Router /share/add [post]
func PostShareDirAdd(c *gin.Context) {

	name := c.PostForm("name")
	comment := c.PostForm("comment")
	path := c.PostForm("path")
	readOnly, _ := strconv.ParseBool(c.DefaultPostForm("read_only", "false"))
	writeable, _ := strconv.ParseBool(c.DefaultPostForm("writeable", "true"))
	browse, _ := strconv.ParseBool(c.DefaultPostForm("browseable", "true"))
	user := c.PostForm("user")

	if len(name) == 0 || len(comment) == 0 || len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}

	var m model2.ShareDirDBModel
	m.Name = name
	m.Comment = comment
	m.Path = path
	m.ReadOnly = readOnly
	m.Writeable = writeable
	m.Browseable = browse
	m.ValidUsers = user

	service.MyService.ShareDirectory().Add(&m)
	service.MyService.ShareDirectory().UpConfig()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary 删除分享
// @Produce  application/json
// @Accept application/json
// @Tags share
// @Security ApiKeyAuth
// @Param id path string true "id"
// @Success 200 {string} string "ok"
// @Router /share/del/{id} [delete]
func DeleteShareDirDel(c *gin.Context) {
	id := c.Param("id")
	service.MyService.ShareDirectory().Delete(id)
	service.MyService.ShareDirectory().UpConfig()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary 分享详情
// @Produce  application/json
// @Accept application/json
// @Tags share
// @Security ApiKeyAuth
// @Param id path string true "id"
// @Success 200 {string} string "ok"
// @Router /share/info/{id} [get]
func GetShareDirInfo(c *gin.Context) {
	id := c.Param("id")
	info := service.MyService.ShareDirectory().Info(id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: info})
}

// @Summary 更新分享详情
// @Produce  application/json
// @Accept application/json
// @Tags share
// @Security ApiKeyAuth
// @Param id path string true "id"
// @Param  path formData string true "要分享的文件路径"
// @Param  name formData string true "名称"
// @Param  comment formData string true "描述"
// @Param  read_only formData bool false "是否只读"
// @Param  writeable formData bool false "是否可写"
// @Param  browseable formData bool false "是否可浏览"
// @Param  user formData string false "用户"
// @Success 200 {string} string "ok"
// @Router /share/update/{id} [put]
func PutShareDirEdit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id == 0 {
		//todo 报错出去
	}

	name := c.PostForm("name")
	comment := c.PostForm("comment")
	path := c.PostForm("path")
	readOnly, _ := strconv.ParseBool(c.DefaultPostForm("read_only", "false"))
	writeable, _ := strconv.ParseBool(c.DefaultPostForm("writeable", "true"))
	browse, _ := strconv.ParseBool(c.DefaultPostForm("browseable", "true"))
	user := c.PostForm("user")

	if len(name) == 0 || len(comment) == 0 || len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}

	var m model2.ShareDirDBModel
	m.Id = uint(id)
	m.Name = name
	m.Comment = comment
	m.Path = path
	m.ReadOnly = readOnly
	m.Writeable = writeable
	m.Browseable = browse
	m.ValidUsers = user
	service.MyService.ShareDirectory().Update(&m)
	service.MyService.ShareDirectory().UpConfig()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}
