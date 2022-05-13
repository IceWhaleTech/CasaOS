package service

import (
	"io/ioutil"
	"net"
	"os"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
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
}
type systemService struct {
	log loger.OLog
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
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return address
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To16() != nil {
				address = append(address, ipNet.IP.String())
			}
		}
	}
	return address
}
func NewSystemService(log loger.OLog) SystemService {
	return &systemService{log: log}
}
