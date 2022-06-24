package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	port2 "github.com/IceWhaleTech/CasaOS/pkg/utils/port"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/version"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Summary check version
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/version/check [get]
func GetSystemCheckVersion(c *gin.Context) {
	need, version := version.IsNeedUpdate(service.MyService.Casa().GetCasaosVersion())
	if need {
		installLog := model2.AppNotify{}
		installLog.State = 0
		installLog.Message = "New version " + version.Version + " is ready, ready to upgrade"
		installLog.Type = types.NOTIFY_TYPE_NEED_CONFIRM
		installLog.CreatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)
		installLog.Name = "CasaOS System"
		service.MyService.Notify().AddLog(installLog)
	}
	data := make(map[string]interface{}, 3)
	data["is_need"] = need
	data["version"] = version
	data["current_version"] = types.CURRENTVERSION
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

// @Summary 系统信息
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/update [post]
func SystemUpdate(c *gin.Context) {
	need, version := version.IsNeedUpdate(service.MyService.Casa().GetCasaosVersion())
	if need {
		service.MyService.System().UpdateSystemVersion(version.Version)
	}
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

//Get system config
func GetSystemConfig(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: json.RawMessage(config.SystemConfigInfo.ConfigStr)})
}

// @Summary  get logs
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/error/logs [get]
func GetCasaOSErrorLogs(c *gin.Context) {
	line, _ := strconv.Atoi(c.DefaultQuery("line", "100"))
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: service.MyService.System().GetCasaOSLogs(line)})
}

// @Summary 修改配置文件
// @Produce  application/json
// @Accept multipart/form-data
// @Tags sys
// @Param config formData string true "config json string"
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/changhead [post]
func PostSetSystemConfig(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	service.MyService.System().UpSystemConfig(string(buf[0:n]), "")
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    json.RawMessage(config.SystemConfigInfo.ConfigStr),
		})
}

//系统配置
func GetSystemConfigDebug(c *gin.Context) {

	array := service.MyService.System().GetSystemConfigDebug()
	disk := service.MyService.System().GetDiskInfo()
	sys := service.MyService.System().GetSysInfo()
	//todo 准备sync需要显示的数据(镜像,容器)
	var systemAppStatus string
	images := service.MyService.Docker().IsExistImage("linuxserver/syncthing")
	systemAppStatus += "Sync img: " + strconv.FormatBool(images) + "\n\t"

	list := service.MyService.App().GetSystemAppList()
	for _, v := range list {
		systemAppStatus += v.Image + ",\n\t"
	}

	systemAppStatus += "Sync Key length: " + strconv.Itoa(len(config.SystemConfigInfo.SyncKey))

	var bugContent string = fmt.Sprintf(`
	 - OS: %s
	 - CasaOS Version: %s
	 - Disk Total: %v 
	 - Disk Used: %v 
	 - Sync State: %s
	 - System Info: %s
	 - Browser: $Browser$ 
	 - Version: $Version$
`, sys.OS, types.CURRENTVERSION, disk.Total>>20, disk.Used>>20, systemAppStatus, array)

	//	array = append(array, fmt.Sprintf("disk,total:%v,used:%v,UsedPercent:%v", disk.Total>>20, disk.Used>>20, disk.UsedPercent))

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: bugContent})
}

//widget配置
func GetWidgetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: json.RawMessage(config.SystemConfigInfo.WidgetList)})
}

// @Summary 修改组件配置文件
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/widget/config [post]
func PostSetWidgetConfig(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	service.MyService.System().UpSystemConfig("", string(buf[0:n]))
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    json.RawMessage(config.SystemConfigInfo.WidgetList),
		})
}

// @Summary get casaos server port
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/port [get]
func GetCasaOSPort(c *gin.Context) {
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    config.ServerInfo.HttpPort,
		})
}

// @Summary edit casaos server port
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Param port json string true "port"
// @Success 200 {string} string "ok"
// @Router /sys/port [put]
func PutCasaOSPort(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)
	portStr := json["port"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusOK,
			model.Result{
				Success: common_err.ERROR,
				Message: err.Error(),
			})
		return
	}

	isAvailable := port2.IsPortAvailable(port, "tcp")
	if !isAvailable {
		c.JSON(http.StatusOK,
			model.Result{
				Success: common_err.PORT_IS_OCCUPIED,
				Message: common_err.GetMsg(common_err.PORT_IS_OCCUPIED),
			})
		return
	}
	service.MyService.System().UpSystemPort(strconv.Itoa(port))
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
		})
}

// @Summary 检查是否进入引导状态
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/init/check [get]
func GetSystemInitCheck(c *gin.Context) {
	data := make(map[string]interface{}, 2)

	if service.MyService.User().GetUserCount() > 0 {
		data["initialized"] = true
		data["key"] = ""
	} else {
		key := uuid.NewV4().String()
		service.UserRegisterHash[key] = key
		data["key"] = key
		data["initialized"] = false
	}
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    data,
		})
}

// @Summary active killing casaos
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/restart [post]
func PostKillCasaOS(c *gin.Context) {
	os.Exit(0)
}

// @Summary Turn off usb auto-mount
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/usb/off [put]
func PutSystemUSBAutoMount(c *gin.Context) {
	status := c.Param("status")
	if status == "on" {
		service.MyService.System().UpdateUSBAutoMount("True")
		service.MyService.System().ExecUSBAutoMountShell("True")
	} else {
		service.MyService.System().UpdateUSBAutoMount("False")
		service.MyService.System().ExecUSBAutoMountShell("False")
	}

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
		})
}

// @Summary Turn off usb auto-mount
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/usb [get]
func GetSystemUSBAutoMount(c *gin.Context) {
	state := "True"
	if config.ServerInfo.USBAutoMount == "False" {
		state = "False"
	}

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    state,
		})
}

// @Summary get system hardware info
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/hardware/info [get]
func GetSystemHardwareInfo(c *gin.Context) {

	data := make(map[string]string, 1)
	data["drive_model"] = service.MyService.System().GetDeviceTree()
	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    data,
		})
}

// @Summary system utilization
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/utilization [get]
func GetSystemUtilization(c *gin.Context) {
	var data = make(map[string]interface{}, 6)

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
		if list[i].Tran == "sata" || list[i].Tran == "nvme" || list[i].Tran == "spi" || list[i].Tran == "sas" {
			temp := service.MyService.Disk().SmartCTL(list[i].Path)
			if reflect.DeepEqual(temp, model.SmartctlA{}) {
				continue
			}

			//list[i].Temperature = temp.Temperature.Current
			if !temp.SmartStatus.Passed {
				healthy = false
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
	data["disk"] = summary
	usbList := service.MyService.Disk().LSBLK(false)
	usb := []model.DriveUSB{}
	for _, v := range usbList {
		if v.Tran == "usb" {
			temp := model.DriveUSB{}
			temp.Model = v.Model
			temp.Name = v.Name
			temp.Size = v.Size
			mountTemp := true
			if len(v.Children) == 0 {
				mountTemp = false
			}
			for _, child := range v.Children {
				if len(child.MountPoint) > 0 {
					avail, _ := strconv.ParseUint(child.FSAvail, 10, 64)
					temp.Avail += avail
					used, _ := strconv.ParseUint(child.FSUsed, 10, 64)
					temp.Used += used
				} else {
					mountTemp = false
				}
			}
			temp.Mount = mountTemp
			usb = append(usb, temp)
		}
	}
	data["usb"] = usb
	cpu := service.MyService.System().GetCpuPercent()
	num := service.MyService.System().GetCpuCoreNum()
	cpuData := make(map[string]interface{})
	cpuData["percent"] = cpu
	cpuData["num"] = num
	data["cpu"] = cpuData
	data["mem"] = service.MyService.System().GetMemInfo()

	//拼装网络信息
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

	data["net"] = newNet

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}

// @Summary Get notification port
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/socket/port [get]
func GetSystemSocketPort(c *gin.Context) {

	c.JSON(http.StatusOK,
		model.Result{
			Success: common_err.SUCCESS,
			Message: common_err.GetMsg(common_err.SUCCESS),
			Data:    config.ServerInfo.SocketPort,
		})
}

// @Summary get cpu info
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/cpu [get]
func GetSystemCupInfo(c *gin.Context) {
	cpu := service.MyService.System().GetCpuPercent()
	num := service.MyService.System().GetCpuCoreNum()
	data := make(map[string]interface{})
	data["percent"] = cpu
	data["num"] = num
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})

}

// @Summary get mem info
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/mem [get]
func GetSystemMemInfo(c *gin.Context) {
	mem := service.MyService.System().GetMemInfo()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: mem})

}

// @Summary get disk info
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/disk [get]
func GetSystemDiskInfo(c *gin.Context) {
	disk := service.MyService.System().GetDiskInfo()
	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: disk})
}

// @Summary get Net info
// @Produce  application/json
// @Accept application/json
// @Tags sys
// @Security ApiKeyAuth
// @Success 200 {string} string "ok"
// @Router /sys/net [get]
func GetSystemNetInfo(c *gin.Context) {
	netList := service.MyService.System().GetNetInfo()

	newNet := []model.IOCountersStat{}
	for _, n := range netList {
		for _, netCardName := range service.MyService.System().GetNet(true) {
			if n.Name == netCardName {
				item := *(*model.IOCountersStat)(unsafe.Pointer(&n))
				item.State = strings.TrimSpace(service.MyService.System().GetNetState(n.Name))
				item.Time = time.Now().Unix()
				newNet = append(newNet, item)
				break
			}
		}
	}

	c.JSON(http.StatusOK, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: newNet})
}
