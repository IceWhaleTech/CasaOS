package v1

import (
	"fmt"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/labstack/echo/v4"
)

func GetSearchResult(ctx echo.Context) error {
	json := make(map[string]string)
	ctx.Bind(&json)
	url := json["url"]

	if url == "" {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "key is empty"})
	}
	// data, err := service.MyService.Other().Search(key)
	data, err := service.MyService.Other().AgentSearch(url)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	}

	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}
