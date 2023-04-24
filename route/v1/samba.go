/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-26 11:08:48
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-17 18:25:42
 * @FilePath: /CasaOS/route/v1/samba.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
	"go.uber.org/zap"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/samba"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/gin-gonic/gin"
)

// service

func GetSambaStatus(c *gin.Context) {
	if status, err := systemctl.IsServiceRunning("smbd"); err != nil || !status {
		c.JSON(http.StatusInternalServerError, model.Result{
			Success: common_err.SERVICE_NOT_RUNNING,
			Message: common_err.GetMsg(common_err.SERVICE_NOT_RUNNING),
		})
		return
	}

	needInit := true
	if file.Exists("/etc/samba/smb.conf") {
		str := file.ReadLine(1, "/etc/samba/smb.conf")
		if strings.Contains(str, "# Copyright (c) 2021-2022 CasaOS Inc. All rights reserved.") {
			needInit = false
		}
	}
	data := make(map[string]string, 1)
	data["need_init"] = fmt.Sprintf("%v", needInit)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

func GetSambaSharesList(c *gin.Context) {
	shares := service.MyService.Shares().GetSharesList()
	shareList := []model.Shares{}
	for _, v := range shares {
		shareList = append(shareList, model.Shares{
			Anonymous: v.Anonymous,
			Path:      v.Path,
			ID:        v.ID,
		})
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: shareList})
}

func PostSambaSharesCreate(c *gin.Context) {
	shares := []model.Shares{}
	c.ShouldBindJSON(&shares)
	for _, v := range shares {
		if v.Path == "" {
			c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INSUFFICIENT_PERMISSIONS, Message: common_err.GetMsg(common_err.INSUFFICIENT_PERMISSIONS)})
			return
		}
		if !file.Exists(v.Path) {
			c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.DIR_NOT_EXISTS, Message: common_err.GetMsg(common_err.DIR_NOT_EXISTS)})
			return
		}
		if len(service.MyService.Shares().GetSharesByPath(v.Path)) > 0 {
			c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.SHARE_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.SHARE_ALREADY_EXISTS)})
			return
		}
		if len(service.MyService.Shares().GetSharesByPath(filepath.Base(v.Path))) > 0 {
			c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.SHARE_NAME_ALREADY_EXISTS, Message: common_err.GetMsg(common_err.SHARE_NAME_ALREADY_EXISTS)})
			return
		}
	}
	for _, v := range shares {
		shareDBModel := model2.SharesDBModel{}
		shareDBModel.Anonymous = true
		shareDBModel.Path = v.Path
		shareDBModel.Name = filepath.Base(v.Path)
		os.Chmod(v.Path, 0o777)
		service.MyService.Shares().CreateShare(shareDBModel)
	}

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: shares})
}

func DeleteSambaShares(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INSUFFICIENT_PERMISSIONS, Message: common_err.GetMsg(common_err.INSUFFICIENT_PERMISSIONS)})
		return
	}
	service.MyService.Shares().DeleteShare(id)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: id})
}

// client
func GetSambaConnectionsList(c *gin.Context) {
	connections := service.MyService.Connections().GetConnectionsList()
	connectionList := []model.Connections{}
	for _, v := range connections {
		connectionList = append(connectionList, model.Connections{
			ID:         v.ID,
			Username:   v.Username,
			Port:       v.Port,
			Host:       v.Host,
			MountPoint: v.MountPoint,
		})
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: connectionList})
}

func PostSambaConnectionsCreate(c *gin.Context) {
	connection := model.Connections{}
	c.ShouldBindJSON(&connection)
	if connection.Port == "" {
		connection.Port = "445"
	}
	if connection.Username == "" || connection.Host == "" {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CHARACTER_LIMIT, Message: common_err.GetMsg(common_err.CHARACTER_LIMIT)})
		return
	}

	// if ok, _ := regexp.MatchString(`^[\w@#*.]{4,30}$`, connection.Password); !ok {
	// 	c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CHARACTER_LIMIT, Message: common_err.GetMsg(common_err.CHARACTER_LIMIT)})
	// 	return
	// }
	// if ok, _ := regexp.MatchString(`^[\w@#*.]{4,30}$`, connection.Username); !ok {
	// 	c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	// 	return
	// }
	// if !ip_helper.IsIPv4(connection.Host) && !ip_helper.IsIPv6(connection.Host) {
	// 	c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	// 	return
	// }
	// if ok, _ := regexp.MatchString("^[0-9]{1,6}$", connection.Port); !ok {
	// 	c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS)})
	// 	return
	// }

	connection.Host = strings.Split(connection.Host, "/")[0]
	// check is exists
	connections := service.MyService.Connections().GetConnectionByHost(connection.Host)
	if len(connections) > 0 {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.Record_ALREADY_EXIST, Message: common_err.GetMsg(common_err.Record_ALREADY_EXIST), Data: common_err.GetMsg(common_err.Record_ALREADY_EXIST)})
		return
	}
	// check connect is ok
	directories, err := samba.GetSambaSharesList(connection.Host, connection.Port, connection.Username, connection.Password)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}

	connectionDBModel := model2.ConnectionsDBModel{}
	connectionDBModel.Username = connection.Username
	connectionDBModel.Password = connection.Password
	connectionDBModel.Host = connection.Host
	connectionDBModel.Port = connection.Port
	connectionDBModel.Directories = strings.Join(directories, ",")
	baseHostPath := "/mnt/" + connection.Host
	connectionDBModel.MountPoint = baseHostPath
	connection.MountPoint = baseHostPath
	file.IsNotExistMkDir(baseHostPath)
	for _, v := range directories {
		mountPoint := baseHostPath + "/" + v
		file.IsNotExistMkDir(mountPoint)
		service.MyService.Connections().MountSmaba(connectionDBModel.Username, connectionDBModel.Host, v, connectionDBModel.Port, mountPoint, connectionDBModel.Password)
	}

	service.MyService.Connections().CreateConnection(&connectionDBModel)

	connection.ID = connectionDBModel.ID
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: connection})
}

func DeleteSambaConnections(c *gin.Context) {
	id := c.Param("id")
	connection := service.MyService.Connections().GetConnectionByID(id)
	if connection.Username == "" {
		c.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.Record_NOT_EXIST, Message: common_err.GetMsg(common_err.Record_NOT_EXIST)})
		return
	}
	mountPointList, err := samba.GetSambaSharesList(connection.Host, connection.Port, connection.Username, connection.Password)
	//mountPointList, err := service.MyService.System().GetDirPath(connection.MountPoint)
	if err != nil {
		c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	baseHostPath := "/mnt/" + connection.Host
	for _, v := range mountPointList {
		if service.IsMounted(baseHostPath + "/" + v) {
			err := service.MyService.Connections().UnmountSmaba(baseHostPath + "/" + v)
			if err != nil {
				logger.Error("unmount smaba error", zap.Error(err), zap.Any("path", baseHostPath+"/"+v))
				c.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
				return
			}
		}

	}
	dir, _ := ioutil.ReadDir(connection.MountPoint)
	if len(dir) == 0 {
		os.RemoveAll(connection.MountPoint)
	}
	service.MyService.Connections().DeleteConnection(id)
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: id})
}
