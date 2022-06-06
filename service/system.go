package service

import (
	"fmt"
	"io/ioutil"
	net2 "net"
	"os"
	"strconv"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemService interface {
	UpSystemConfig(str string, widget string)
	UpdateSystemVersion(version string)
	GetSystemConfigDebug() []string
	GetCasaOSLogs(lineNumber int) string
	UpdateAssist()
	UpSystemPort(port string)
	GetTimeZone() string
	UpdateUSBAutoMount(state string)
	ExecUSBAutoMountShell(state string)
	UpAppOrderFile(str string)
	GetAppOrderFile() []byte
	GetNet(physics bool) []string
	GetNetInfo() []net.IOCountersStat
	GetCpuCoreNum() int
	GetCpuPercent() float64
	GetMemInfo() *mem.VirtualMemoryStat
}
type systemService struct {
	log loger.OLog
}

func (c *systemService) GetMemInfo() *mem.VirtualMemoryStat {
	memInfo, _ := mem.VirtualMemory()
	memInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memInfo.UsedPercent), 64)
	return memInfo
}

func (c *systemService) GetCpuPercent() float64 {
	percent, _ := cpu.Percent(0, false)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", percent[0]), 64)
	return value
}

func (c *systemService) GetCpuCoreNum() int {
	count, _ := cpu.Counts(false)
	return count
}

func (c *systemService) GetNetInfo() []net.IOCountersStat {
	parts, _ := net.IOCounters(true)
	//fmt.Println(net.ConntrackStatsWithContext(true))
	return parts
}
func (c *systemService) GetNet(physics bool) []string {
	t := "1"
	if physics {
		t = "2"
	}
	return command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetNetCard " + t)
}

func (s *systemService) UpdateSystemVersion(version string) {
	//command2.OnlyExec(config.AppInfo.ProjectPath + "/shell/tool.sh -r " + version)
	//s.log.Error(config.AppInfo.ProjectPath + "/shell/tool.sh -r " + version)
	s.log.Error(command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/tools.sh ;update " + version))
	//s.log.Error(command2.ExecResultStr(config.AppInfo.ProjectPath + "/shell/tool.sh -r " + version))
}
func (s *systemService) UpdateAssist() {
	s.log.Error(command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/assist.sh"))
}

func (s *systemService) GetTimeZone() string {
	return command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetTimeZone")
}

func (s *systemService) ExecUSBAutoMountShell(state string) {
	if state == "False" {
		command2.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;USB_Remove_File")
	} else {
		command2.OnlyExec("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;USB_Move_File")
	}

}

func (s *systemService) GetSystemConfigDebug() []string {
	return command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetSysInfo")
}
func (s *systemService) UpSystemConfig(str string, widget string) {
	if len(str) > 0 && str != config.SystemConfigInfo.ConfigStr {
		config.Cfg.Section("system").Key("ConfigStr").SetValue(str)
		config.SystemConfigInfo.ConfigStr = str
	}
	if len(widget) > 0 && widget != config.SystemConfigInfo.WidgetList {
		config.Cfg.Section("system").Key("WidgetList").SetValue(widget)
		config.SystemConfigInfo.WidgetList = widget
	}
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
}
func (s *systemService) UpAppOrderFile(str string) {
	file.WriteToPath([]byte(str), config.AppInfo.ProjectPath+"/conf", "app_order.json")
}
func (s *systemService) GetAppOrderFile() []byte {
	return file.ReadFullFile(config.AppInfo.ProjectPath + "/conf/app_order.json")
}
func (s *systemService) UpdateUSBAutoMount(state string) {
	config.ServerInfo.USBAutoMount = state
	config.Cfg.Section("server").Key("USBAutoMount").SetValue(state)
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
}
func (s *systemService) UpSystemPort(port string) {
	if len(port) > 0 && port != config.ServerInfo.HttpPort {
		config.Cfg.Section("server").Key("HttpPort").SetValue(port)
		config.ServerInfo.HttpPort = port
	}
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
}
func (s *systemService) GetCasaOSLogs(lineNumber int) string {
	file, err := os.Open(s.log.Path())
	if err != nil {
		return err.Error()
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err.Error()
	}

	return string(content)
}

func GetDeviceAllIP() []string {
	var address []string
	addrs, err := net2.InterfaceAddrs()
	if err != nil {
		return address
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net2.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To16() != nil {
				address = append(address, ipNet.IP.String())
			}
		}
	}
	return address
}
func NewSystemService() SystemService {
	return &systemService{}
}
