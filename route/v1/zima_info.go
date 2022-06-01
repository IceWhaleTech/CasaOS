/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-27 18:07:13
 * @FilePath: /CasaOS/route/v1/zima_info.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/IceWhaleTech/CasaOS/model"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取cpu信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/getcpuinfo [get]
func CupInfo(c *gin.Context) {
	//检查参数是否正确
	cpu := service.MyService.System().GetCpuPercent()
	num := service.MyService.System().GetCpuCoreNum()
	data := make(map[string]interface{})
	data["percent"] = cpu
	data["num"] = num
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: data})

}

// @Summary 获取内存信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/getmeminfo [get]
func MemInfo(c *gin.Context) {

	//检查参数是否正确
	mem := service.MyService.System().GetMemInfo()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: mem})

}

// @Summary 获取硬盘信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/getdiskinfo [get]
func DiskInfo(c *gin.Context) {
	disk := service.MyService.ZiMa().GetDiskInfo()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: disk})
}

// @Summary 获取网络信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/getnetinfo [get]
func NetInfo(c *gin.Context) {
	netList := service.MyService.System().GetNetInfo()

	newNet := []model.IOCountersStat{}
	for _, n := range netList {
		for _, netCardName := range service.MyService.System().GetNet(true) {
			if n.Name == netCardName {
				item := *(*model.IOCountersStat)(unsafe.Pointer(&n))
				item.State = strings.TrimSpace(service.MyService.ZiMa().GetNetState(n.Name))
				item.Time = time.Now().Unix()
				newNet = append(newNet, item)
				break
			}
		}
	}

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: newNet})
}

// @Summary 获取信息系统信息
// @Produce  application/json
// @Accept application/json
// @Tags zima
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zima/sysinfo [get]
func SysInfo(c *gin.Context) {
	info := service.MyService.ZiMa().GetSysInfo()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}
