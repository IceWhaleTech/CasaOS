/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-11 16:02:29
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-29 14:14:17
 * @FilePath: /CasaOS/route/v1/storage.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

func GetStorageList(c *gin.Context) {

	storages := []model.Storages{}
	disks := service.MyService.Disk().LSBLK(false)
	diskNumber := 1
	children := 1
	for _, d := range disks {
		children = 1
		if d.Tran == "sata" || d.Tran == "nvme" || d.Tran == "spi" || d.Tran == "sas" || strings.Contains(d.SubSystems, "virtio") || (d.Tran == "ata" && d.Type == "disk") {
			storageArr := []model.Storage{}
			temp := service.MyService.Disk().SmartCTL(d.Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				temp.SmartStatus.Passed = true
			}
			for _, v := range d.Children {
				if v.MountPoint != "" {
					stor := model.Storage{}
					stor.MountPoint = v.MountPoint
					stor.Size = v.FSSize
					stor.Avail = v.FSAvail
					stor.Path = v.Path
					stor.Type = v.FsType
					stor.DriveName = "System"
					if len(v.Label) == 0 {
						stor.Label = "Storage" + strconv.Itoa(diskNumber) + "_" + strconv.Itoa(children)
						children += 1
					} else {
						stor.Label = v.Label
					}
					storageArr = append(storageArr, stor)
				}
			}
			if len(storageArr) > 0 {
				storages = append(storages, model.Storages{
					DiskName: d.Model,
					Path:     d.Path,
					Size:     d.Size,
					Children: storageArr,
				})
				diskNumber += 1
			}
		}
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: storages})
}
