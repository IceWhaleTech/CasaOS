package v1

import (
	json2 "encoding/json"
	"net/http"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

// @Summary 登录zerotier获取token
// @Produce  application/json
// @Accept multipart/form-data
// @Tags zerotier
// @Param username formData string true "User name"
// @Param pwd  formData string true "password"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/login [post]
func ZeroTierGetToken(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("pwd")
	if len(username) == 0 || len(pwd) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	errInfo := service.MyService.ZeroTier().GetToken(username, pwd)

	if len(errInfo) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: oasis_err2.GetMsg(oasis_err2.GET_TOKEN_ERROR)})
	} else {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
	}
}

// @Summary 注册zerotier
// @Produce  application/json
// @Accept multipart/form-data
// @Tags zerotier
// @Param firstName formData string true "first name"
// @Param pwd  formData string true "password"
// @Param email  formData string true "email"
// @Param lastName  formData string true "last name"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/register [post]
func ZeroTierRegister(c *gin.Context) {
	firstName := c.PostForm("firstName")
	pwd := c.PostForm("pwd")
	email := c.PostForm("email")
	lastName := c.PostForm("lastName")
	if len(firstName) == 0 || len(pwd) == 0 || len(email) == 0 || len(lastName) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	errInfo := service.MyService.ZeroTier().ZeroTierRegister(email, lastName, firstName, pwd)
	if len(errInfo) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
	} else {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.ERROR, Message: errInfo})
	}
}

// @Summary 是否需要登录zerotier
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Success 200 {string} string "false:需要登录，true:不需要登录"
// @Router /zerotier/islogin [get]
func ZeroTierIsNeedLogin(c *gin.Context) {
	if len(config.ZeroTierInfo.Token) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: false})
	} else {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: true})
	}
}

// @Summary 获取zerotier网络列表
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/list [get]
func ZeroTierGetNetworkList(c *gin.Context) {
	jsonList, joined := service.MyService.ZeroTier().ZeroTierNetworkList(config.ZeroTierInfo.Token)
	rdata := make(map[string]interface{})
	rdata["network_list"] = jsonList
	rdata["joined"] = joined
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: rdata})
}

// @Summary 获取zerotier网络详情
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Success 200 {string} string "ok"
// @Router /zerotier/info/{id} [get]
func ZeroTierGetNetworkGetInfo(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info, joined := service.MyService.ZeroTier().ZeroTierGetInfo(config.ZeroTierInfo.Token, id)
	rdata := make(map[string]interface{})
	rdata["info"] = info
	rdata["joined"] = joined
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: rdata})
}

//// @Summary 获取zerotier网络状态
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Success 200 {string} string "ok"
//// @Router /zerotier/status [get]
//func ZeroTierGetNetworkGetStatus(c *gin.Context) {
//	status := service.MyService.ZeroTier().ZeroTierGetStatus(config.ZeroTierInfo.Token)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: status})
//}

//// @Summary 修改网络类型
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param type formData string true "Private true/false"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/type/{id} [put]
//func ZeroTierEditType(c *gin.Context) {
//	id := c.Param("id")
//	t := c.PostForm("type")
//	if len(id) == 0 || len(t) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"config":{"private":` + t + `}}`
//	info := service.MyService.ZeroTier().EditNetwork(config.ZeroTierInfo.Token, postData, id)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

//// @Summary 修改名称
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param name formData string true "需要过滤特殊字符串"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/name/{id} [put]
//func ZeroTierEditName(c *gin.Context) {
//	id := c.Param("id")
//	name := c.PostForm("name")
//	if len(id) == 0 || len(name) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"config":{"name":"` + name + `"}}`
//	info := service.MyService.ZeroTier().EditNetwork(config.ZeroTierInfo.Token, postData, id)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

//// @Summary V6Assign (注意三个属性需要一起传过来,不传的会被zerotier设置成false)
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param v6plan formData string false "true/false"
//// @Param rfc formData string false "true/false"
//// @Param auto formData string false "true/false"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/v6assign/{id} [put]
//func ZeroTierEditV6Assign(c *gin.Context) {
//	id := c.Param("id")
//	v6plan := c.PostForm("v6plan")
//	rfc := c.PostForm("rfc")
//	auto := c.PostForm("auto")
//	if len(id) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	var spicing string
//	if len(v6plan) > 0 {
//		spicing = `"6plane":` + v6plan
//	}
//	if len(rfc) > 0 {
//		if len(spicing) > 0 {
//			spicing += ","
//		}
//		spicing += `"rfc4193":` + rfc
//	}
//
//	if len(auto) > 0 {
//		if len(spicing) > 0 {
//			spicing += ","
//		}
//		spicing += `"zt":` + auto
//	}
//	postData := `{"config":{"v6AssignMode":{` + spicing + `}}}`
//	info := service.MyService.ZeroTier().EditNetwork(config.ZeroTierInfo.Token, postData, id)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

//// @Summary Broadcast
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param broadcast formData string true "true/false"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/broadcast/{id} [put]
//func ZeroTierEditBroadcast(c *gin.Context) {
//	id := c.Param("id")
//	broadcast := c.PostForm("broadcast")
//	if len(id) == 0 || len(broadcast) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"config":{"enableBroadcast":` + broadcast + `}}`
//	info := service.MyService.ZeroTier().EditNetwork(config.ZeroTierInfo.Token, postData, id)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

// @Summary 网络列表
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Success 200 {string} string "ok"
// @Router /zerotier/member/{id} [get]
func ZeroTierMemberList(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.ZeroTier().MemberList(config.ZeroTierInfo.Token, id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary create new network
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/create [post]
func ZeroTierCreateNetwork(c *gin.Context) {
	info := service.MyService.ZeroTier().CreateNetwork(config.ZeroTierInfo.Token)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

//// @Summary 通过/拒绝客户端
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param mId path string true "member_id"
//// @Param auth formData string true "true/false"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/member/{id}/auth/{mId} [put]
//func ZeroTierMemberAuth(c *gin.Context) {
//	id := c.Param("id")
//	mId := c.Param("mId")
//	auth := c.PostForm("auth")
//	if len(id) == 0 || len(mId) == 0 || len(auth) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"config":{"authorized":` + auth + `}}`
//	info := service.MyService.ZeroTier().EditNetworkMember(config.ZeroTierInfo.Token, postData, id, mId)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

//// @Summary 修改名字
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param mId path string true "member_id"
//// @Param name formData string true "name"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/member/{id}/name/{mId} [put]
//func ZeroTierMemberName(c *gin.Context) {
//	id := c.Param("id")
//	mId := c.Param("mId")
//	name := c.PostForm("name")
//	if len(id) == 0 || len(mId) == 0 || len(name) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"name":"` + name + `"}`
//	info := service.MyService.ZeroTier().EditNetworkMember(config.ZeroTierInfo.Token, postData, id, mId)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

//// @Summary 修改桥接
//// @Produce  application/json
//// @Accept application/json
//// @Tags zerotier
//// @Security ApiKeyAuth
//// @Param id path string true "network id"
//// @Param mId path string true "member_id"
//// @Param bridge formData string true "true/false"
//// @Success 200 {string} string "ok"
//// @Router /zerotier/member/{id}/bridge/{mId} [put]
//func ZeroTierMemberBridge(c *gin.Context) {
//	id := c.Param("id")
//	mId := c.Param("mId")
//	bridge := c.PostForm("bridge")
//	if len(id) == 0 || len(mId) == 0 || len(bridge) == 0 {
//		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
//		return
//	}
//	postData := `{"config":{"activeBridge":` + bridge + `}}`
//	info := service.MyService.ZeroTier().EditNetworkMember(config.ZeroTierInfo.Token, postData, id, mId)
//	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
//}

// @Summary 修改网络
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Param json formData string true "json数据"
// @Success 200 {string} string "ok"
// @Router /zerotier/edit/{id} [put]
func ZeroTierEdit(c *gin.Context) {
	id := c.Param("id")
	json := c.PostForm("json")
	if len(id) == 0 || len(json) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.ZeroTier().EditNetwork(config.ZeroTierInfo.Token, json, id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 获取已加入的网络
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/joined/list [get]
func ZeroTierJoinedList(c *gin.Context) {
	info := service.MyService.ZeroTier().GetJoinNetworks()
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: json2.RawMessage(info)})
}

// @Summary 修改网络用户信息
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Param mId path string true "mId"
// @Param json formData string true "json数据"
// @Success 200 {string} string "ok"
// @Router /zerotier/member/{id}/edit/{mId} [put]
func ZeroTierMemberEdit(c *gin.Context) {
	id := c.Param("id")
	mId := c.Param("mId")
	json := c.PostForm("json")
	if len(id) == 0 || len(json) == 0 || len(mId) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.ZeroTier().EditNetworkMember(config.ZeroTierInfo.Token, json, id, mId)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 删除网络中的用户
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Param mId path string true "member_id"
// @Success 200 {string} string "ok"
// @Router /zerotier/member/{id}/del/{mId} [delete]
func ZeroTierMemberDelete(c *gin.Context) {
	id := c.Param("id")
	mId := c.Param("mId")
	if len(id) == 0 || len(mId) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.ZeroTier().DeleteMember(config.ZeroTierInfo.Token, id, mId)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 删除网络
// @Produce  application/json
// @Accept application/json
// @Tags zerotier
// @Security ApiKeyAuth
// @Param id path string true "network id"
// @Success 200 {string} string "ok"
// @Router /zerotier/network/{id}/del [delete]
func ZeroTierDeleteNetwork(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	info := service.MyService.ZeroTier().DeleteNetwork(config.ZeroTierInfo.Token, id)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 加入网络
// @Produce  application/json
// @Accept multipart/form-data
// @Tags zerotier
// @Param id path string true "network id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/join/{id} [post]
func ZeroTierJoinNetwork(c *gin.Context) {
	networkId := c.Param("id")
	if len(networkId) != 16 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	for _, v := range networkId {
		if !service.MyService.ZeroTier().NetworkIdFilter(v) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
			return
		}
	}
	service.MyService.ZeroTier().ZeroTierJoinNetwork(networkId)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary 获取zerotier网络列表
// @Produce  application/json
// @Accept multipart/form-data
// @Tags zerotier
// @Param id path string true "network id"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /zerotier/leave/{id} [post]
func ZeroTierLeaveNetwork(c *gin.Context) {
	networkId := c.Param("id")

	if len(networkId) != 16 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}

	for _, v := range networkId {
		if !service.MyService.ZeroTier().NetworkIdFilter(v) {
			c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
			return
		}
	}
	service.MyService.ZeroTier().ZeroTierLeaveNetwork(networkId)

	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}
