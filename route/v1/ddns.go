package v1

import (
	"fmt"
	"github.com/forease/gotld"
	"github.com/gin-gonic/gin"
	"net/http"
	"oasis/model"
	ip_helper2 "oasis/pkg/utils/ip_helper"
	oasis_err2 "oasis/pkg/utils/oasis_err"
	"oasis/service"
	model2 "oasis/service/model"
	"os/exec"
	"strconv"
	"strings"
)

// @Summary 获取可以设置的ddns列表
// @Produce  application/json
// @Accept application/json
// @Tags ddns
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/getlist [get]
func DDNSGetDomainList(c *gin.Context) {

	host, domain, tld := gotld.GetSubdomain("bbb.aaa.liru-05.com.cn", 3)
	fmt.Println(strings.Replace(host, "."+domain, "", 1))
	fmt.Println(domain)
	fmt.Println(tld)

	data := make(map[string]interface{}, 2)
	t, api := service.MyService.DDNS().GetType("godaddy")
	data["godaddy"] = &model.GoDaddyModel{Type: t, ApiHost: api}
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err2.SUCCESS,
			Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
			Data:    data,
		})
	return
}

// @Summary 添加新的ddns（按给定模型返回内容）
// @Produce  application/json
// @Accept multipart/form-data
// @Tags ddns
// @Param type formData string true "类型"
// @Param host formData string true "host"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/set [post]
func DDNSAddConfig(c *gin.Context) {
	t, _ := strconv.Atoi(c.PostForm("type"))
	host := c.PostForm("host")
	_, domain, _ := gotld.GetSubdomain("host", 3)
	sub := strings.ReplaceAll(host, "."+domain, "")

	if service.MyService.DDNS().IsExis(t, domain, sub) {
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err2.ERROR,
				Message: "Repeat add",
			})
		return
	}
	var m model2.DDNSUpdataDBModel
	c.Bind(&m)
	if err := service.MyService.DDNS().SaveConfig(m); err != nil {
		c.JSON(http.StatusOK,
			model.Result{
				Success: oasis_err2.ERROR,
				Message: oasis_err2.GetMsg(oasis_err2.ERROR),
				Data:    err.Error(),
			})
		return
	}
	c.JSON(http.StatusOK,
		model.Result{
			Success: oasis_err2.SUCCESS,
			Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
		})
}

// @Summary 获取ip,仅做展示使用
// @Produce  application/json
// @Accept application/json
// @Tags ddns
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/ip [get]
func DDNSGetIP(c *gin.Context) {
	ipv4, ipv6 := service.MyService.DDNS().GetExternalIP()
	var ipjson = make(map[string]string, 2)
	ipjson["ipv4"] = ipv4
	ipjson["ipv6"] = ipv6
	c.JSON(http.StatusOK, &model.Result{
		Success: oasis_err2.SUCCESS,
		Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
		Data:    ipjson,
	})
}

// @Summary 测试网址是否可以ping通
// @Produce  application/json
// @Accept application/json
// @Tags ddns
// @Param  api_host path int true "api地址"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/ping/{api_host} [get]
func DDNSPing(c *gin.Context) {
	url := c.Param("api_host")
	url = strings.ReplaceAll(url, "https://", "")
	url = strings.ReplaceAll(url, "http://", "")
	cmd := exec.Command("ping", url, "-c", "1", "-W", "5")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusOK, &model.Result{
			Success: oasis_err2.ERROR,
			Message: err.Error(),
			Data:    false,
		})
	} else {
		c.JSON(http.StatusOK, &model.Result{
			Success: oasis_err2.SUCCESS,
			Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
			Data:    true,
		})
	}
}

// @Summary 获取已设置的列表
// @Produce  application/json
// @Accept application/json
// @Tags ddns
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/list [get]
func DDNSConfigList(c *gin.Context) {
	j := service.MyService.DDNS().GetConfigList()
	ip4 := ip_helper2.GetExternalIPV4()
	ip6 := ip_helper2.GetExternalIPV6()
	for i := 0; i < len(*j); i++ {
		(*j)[i].IPV6 = ip6
		(*j)[i].IPV4 = ip4
		cmd := exec.Command("ping", (*j)[i].Host+"."+(*j)[i].Domain, "-c", "1", "-W", "3")
		err := cmd.Run()
		if err != nil {
			(*j)[i].State = false
		} else {
			(*j)[i].State = true
		}
	}
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: j})
}

// @Summary 删除ddns
// @Produce  application/json
// @Accept application/json
// @Tags ddns
// @Param id path int true "ID"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /ddns/delete/{id} [delete]
func DDNSDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	isok := service.MyService.DDNS().DeleteConfig(uint(id))
	c.JSON(http.StatusOK, &model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: isok})
}
