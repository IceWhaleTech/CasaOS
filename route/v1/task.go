/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 14:30:44
 * @FilePath: /CasaOS/route/v1/task.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

func PutTaskUpdate(c *gin.Context) {
	service.MyService.Task().SyncTaskService()

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
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
		c.JSON(http.StatusOK, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
		return
	}
	var m model2.TaskDBModel
	m.Id = uint(id)
	m.State = types.TASK_STATE_COMPLETED
	service.MyService.Task().Update(&m)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

func PostTaskAdd(c *gin.Context) {
	var m model2.TaskDBModel
	c.BindJSON(&m)
	service.MyService.Task().Add(&m)
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}
