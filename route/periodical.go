//go:build !darwin
// +build !darwin

/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-07-01 15:11:36
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-09-05 16:28:46
 * @FilePath: /CasaOS/route/periodical.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package route

import (
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

	var cpuModel = "arm"
	if cpu := service.MyService.System().GetCpuInfo(); len(cpu) > 0 {
		if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "intel") > 0 {
			cpuModel = "intel"
		} else if strings.Count(strings.ToLower(strings.TrimSpace(cpu[0].ModelName)), "amd") > 0 {
			cpuModel = "amd"
		}
	}

	num := service.MyService.System().GetCpuCoreNum()
	cpuData := make(map[string]interface{})
	cpuData["percent"] = cpu
	cpuData["num"] = num
	cpuData["temperature"] = service.MyService.System().GetCPUTemperature()
	cpuData["power"] = service.MyService.System().GetCPUPower()
	cpuData["model"] = cpuModel

	memInfo := service.MyService.System().GetMemInfo()

	service.MyService.Notify().SendAllHardwareStatusBySocket(memInfo, cpuData, newNet)

}

// func MonitoryUSB() {
// 	var matcher netlink.Matcher

// 	conn := new(netlink.UEventConn)
// 	if err := conn.Connect(netlink.UdevEvent); err != nil {
// 		loger.Error("udev err", zap.Any("Unable to connect to Netlink Kobject UEvent socket", err))
// 	}
// 	defer conn.Close()

// 	queue := make(chan netlink.UEvent)
// 	errors := make(chan error)
// 	quit := conn.Monitor(queue, errors, matcher)

// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
// 	go func() {
// 		<-signals
// 		close(quit)
// 		os.Exit(0)
// 	}()

// 	for {
// 		select {
// 		case uevent := <-queue:
// 			if uevent.Env["DEVTYPE"] == "disk" {
// 				time.Sleep(time.Microsecond * 500)
// 				SendUSBBySocket()
// 				continue
// 			}
// 		case err := <-errors:
// 			loger.Error("udev err", zap.Any("err", err))
// 		}
// 	}

// }
