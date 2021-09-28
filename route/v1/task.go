package v1

import (
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 获取task列表
// @Produce  application/json
// @Accept application/json
// @Tags task
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /task/list [get]
func GetTaskList(c *gin.Context) {
	//list := service.MyService.Task().List(true)
	list := service.MyService.Task().GetServerTasks()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS), Data: list})
}

func PutTaskUpdate(c *gin.Context) {
	service.MyService.Task().SyncTaskService()

	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

// @Summary 标记task已完成
// @Produce  application/json
// @Accept application/json
// @Tags task
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /task/completion/{id} [put]
func PutTaskMarkerCompletion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err.INVALID_PARAMS, Message: oasis_err.GetMsg(oasis_err.INVALID_PARAMS)})
		return
	}
	var m model2.TaskDBModel
	m.Id = uint(id)
	m.State = types.TASK_STATE_COMPLETED
	service.MyService.Task().Update(&m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}

func PostTaskAdd(c *gin.Context) {
	var m model2.TaskDBModel
	c.BindJSON(&m)
	service.MyService.Task().Add(&m)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err.SUCCESS, Message: oasis_err.GetMsg(oasis_err.SUCCESS)})
}
