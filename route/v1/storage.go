/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-11 16:02:29
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-04 11:27:25
 * @FilePath: /CasaOS/route/v1/storage.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package v1

import (
	"reflect"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

func GetStorageList(c *gin.Context) {
	system := c.Query("system")
	storages := []model.Storages{}
	disks := service.MyService.Disk().LSBLK(false)
	diskNumber := 1
	children := 1
	findSystem := 0
	for _, d := range disks {
		if d.Tran != "usb" {
			tempSystemDisk := false
			children = 1
			tempDisk := model.Storages{
				DiskName: d.Model,
				Path:     d.Path,
				Size:     d.Size,
			}

			storageArr := []model.Storage{}
			temp := service.MyService.Disk().SmartCTL(d.Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				temp.SmartStatus.Passed = true
			}
			for _, v := range d.Children {
				if v.MountPoint != "" {
					if findSystem == 0 {
						if v.MountPoint == "/" {
							tempDisk.DiskName = "System"
							findSystem = 1
							tempSystemDisk = true
						}
						if len(v.Children) > 0 {
							for _, c := range v.Children {
								if c.MountPoint == "/" {
									tempDisk.DiskName = "System"
									findSystem = 1
									tempSystemDisk = true
									break
								}
							}
						}
					}

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
				if tempSystemDisk && len(system) > 0 {
					tempStorageArr := []model.Storage{}
					for i := 0; i < len(storageArr); i++ {
						if storageArr[i].MountPoint != "/boot/efi" && storageArr[i].Type != "swap" {
							tempStorageArr = append(tempStorageArr, storageArr[i])
						}
					}
					tempDisk.Children = tempStorageArr
					storages = append(storages, tempDisk)
					diskNumber += 1
				} else if !tempSystemDisk {
					tempDisk.Children = storageArr
					storages = append(storages, tempDisk)
					diskNumber += 1
				}

			}
		}
	}

	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: storages})
}
