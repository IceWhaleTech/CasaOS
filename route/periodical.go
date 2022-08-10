/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-01 15:11:36
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-03 14:49:15
 * @FilePath: /CasaOS/route/periodical.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/service"
)

func SendNetINfoBySocket() {
	netList := service.MyService.System().GetNetInfo()
	newNet := []model.IOCountersStat{}
	nets := service.MyService.System().GetNet(true)
	for _, n := range netList {
		for _, netCardName := range nets {
			if n.Name == netCardName {
				item := *(*model.IOCountersStat)(unsafe.Pointer(&n))
				item.State = strings.TrimSpace(service.MyService.System().GetNetState(n.Name))
				item.Time = time.Now().Unix()
				newNet = append(newNet, item)
				break
			}
		}
	}
	service.MyService.Notify().SendNetInfoBySocket(newNet)
}

func SendCPUBySocket() {
	cpu := service.MyService.System().GetCpuPercent()
	num := service.MyService.System().GetCpuCoreNum()
	cpuData := make(map[string]interface{})
	cpuData["percent"] = cpu
	cpuData["num"] = num
	service.MyService.Notify().SendCPUInfoBySocket(cpuData)
}

func SendMemBySocket() {
	service.MyService.Notify().SendMemInfoBySocket(service.MyService.System().GetMemInfo())
}

func SendDiskBySocket() {
	list := service.MyService.Disk().LSBLK(true)

	summary := model.Summary{}
	healthy := true
	findSystem := 0

	for i := 0; i < len(list); i++ {
		if len(list[i].Children) > 0 && findSystem == 0 {

			for j := 0; j < len(list[i].Children); j++ {

				if len(list[i].Children[j].Children) > 0 {
					for _, v := range list[i].Children[j].Children {
						if v.MountPoint == "/" {
							s, _ := strconv.ParseUint(v.FSSize, 10, 64)
							a, _ := strconv.ParseUint(v.FSAvail, 10, 64)
							u, _ := strconv.ParseUint(v.FSUsed, 10, 64)
							summary.Size += s
							summary.Avail += a
							summary.Used += u
							findSystem = 1
							break
						}
					}
				} else {
					if list[i].Children[j].MountPoint == "/" {
						s, _ := strconv.ParseUint(list[i].Children[j].FSSize, 10, 64)
						a, _ := strconv.ParseUint(list[i].Children[j].FSAvail, 10, 64)
						u, _ := strconv.ParseUint(list[i].Children[j].FSUsed, 10, 64)
						summary.Size += s
						summary.Avail += a
						summary.Used += u
						findSystem = 1
						break
					}
				}
			}

		}
		if findSystem == 1 {
			findSystem += 1
			continue
		}
		if list[i].Tran == "sata" || list[i].Tran == "nvme" || list[i].Tran == "spi" || list[i].Tran == "sas" || strings.Contains(list[i].SubSystems, "virtio") || (list[i].Tran == "ata" && list[i].Type == "disk") {
			temp := service.MyService.Disk().SmartCTL(list[i].Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				healthy = true
			} else {
				healthy = temp.SmartStatus.Passed
			}

			//list[i].Temperature = temp.Temperature.Current

			if len(list[i].Children) > 0 {
				for _, v := range list[i].Children {
					s, _ := strconv.ParseUint(v.FSSize, 10, 64)
					a, _ := strconv.ParseUint(v.FSAvail, 10, 64)
					u, _ := strconv.ParseUint(v.FSUsed, 10, 64)
					summary.Size += s
					summary.Avail += a
					summary.Used += u
				}
			}

		}
	}

	summary.Health = healthy
	service.MyService.Notify().SendDiskInfoBySocket(summary)
}

func SendUSBBySocket() {
	usbList := service.MyService.Disk().LSBLK(false)
	usb := []model.DriveUSB{}
	for _, v := range usbList {
		if v.Tran == "usb" {
			temp := model.DriveUSB{}
			temp.Model = v.Model
			temp.Name = v.Name
			temp.Size = v.Size
			for _, child := range v.Children {
				if len(child.MountPoint) > 0 {
					avail, _ := strconv.ParseUint(child.FSAvail, 10, 64)
					temp.Avail += avail

				}
			}

			usb = append(usb, temp)
		}
	}
	service.MyService.Notify().SendUSBInfoBySocket(usb)
}

func SendAllHardwareStatusBySocket() {

	netList := service.MyService.System().GetNetInfo()
	newNet := []model.IOCountersStat{}
	nets := service.MyService.System().GetNet(true)
	for _, n := range netList {
		for _, netCardName := range nets {
			if n.Name == netCardName {
				item := *(*model.IOCountersStat)(unsafe.Pointer(&n))
				item.State = strings.TrimSpace(service.MyService.System().GetNetState(n.Name))
				item.Time = time.Now().Unix()
				newNet = append(newNet, item)
				break
			}
		}
	}
	cpu := service.MyService.System().GetCpuPercent()
	num := service.MyService.System().GetCpuCoreNum()
	cpuData := make(map[string]interface{})
	cpuData["percent"] = cpu
	cpuData["num"] = num

	list := service.MyService.Disk().LSBLK(true)

	summary := model.Summary{}
	healthy := true
	findSystem := 0

	for i := 0; i < len(list); i++ {
		if len(list[i].Children) > 0 && findSystem == 0 {

			for j := 0; j < len(list[i].Children); j++ {

				if len(list[i].Children[j].Children) > 0 {
					for _, v := range list[i].Children[j].Children {
						if v.MountPoint == "/" {
							s, _ := strconv.ParseUint(v.FSSize, 10, 64)
							a, _ := strconv.ParseUint(v.FSAvail, 10, 64)
							u, _ := strconv.ParseUint(v.FSUsed, 10, 64)
							summary.Size += s
							summary.Avail += a
							summary.Used += u
							findSystem = 1
							break
						}
					}
				} else {
					if list[i].Children[j].MountPoint == "/" {
						s, _ := strconv.ParseUint(list[i].Children[j].FSSize, 10, 64)
						a, _ := strconv.ParseUint(list[i].Children[j].FSAvail, 10, 64)
						u, _ := strconv.ParseUint(list[i].Children[j].FSUsed, 10, 64)
						summary.Size += s
						summary.Avail += a
						summary.Used += u
						findSystem = 1
						break
					}
				}
			}

		}
		if findSystem == 1 {
			findSystem += 1
			continue
		}
		if list[i].Tran == "sata" || list[i].Tran == "nvme" || list[i].Tran == "spi" || list[i].Tran == "sas" || strings.Contains(list[i].SubSystems, "virtio") || (list[i].Tran == "ata" && list[i].Type == "disk") {
			temp := service.MyService.Disk().SmartCTL(list[i].Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				healthy = true
			} else {
				healthy = temp.SmartStatus.Passed
			}
			if len(list[i].Children) > 0 {
				for _, v := range list[i].Children {
					s, _ := strconv.ParseUint(v.FSSize, 10, 64)
					a, _ := strconv.ParseUint(v.FSAvail, 10, 64)
					u, _ := strconv.ParseUint(v.FSUsed, 10, 64)
					summary.Size += s
					summary.Avail += a
					summary.Used += u
				}
			}

		}
	}

	summary.Health = healthy

	usbList := service.MyService.Disk().LSBLK(false)
	usb := []model.DriveUSB{}
	for _, v := range usbList {
		if v.Tran == "usb" {
			temp := model.DriveUSB{}
			temp.Model = v.Model
			temp.Name = v.Name
			temp.Size = v.Size
			for _, child := range v.Children {
				if len(child.MountPoint) > 0 {
					avail, _ := strconv.ParseUint(child.FSAvail, 10, 64)
					temp.Avail += avail
				}
			}
			usb = append(usb, temp)
		}
	}
	memInfo := service.MyService.System().GetMemInfo()

	service.MyService.Notify().SendAllHardwareStatusBySocket(summary, usb, memInfo, cpuData, newNet)

}
