package service

import (
	"io/ioutil"
	"os"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
)

type SystemService interface {
	UpSystemConfig(str string, widget string)
	UpdateSystemVersion(version string)
	GetSystemConfigDebug() []string
	GetCasaOSLogs(lineNumber int) string
	UpdateAssist()
	UpSystemPort(port string)
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

func NewSystemService(log loger.OLog) SystemService {
	return &systemService{log: log}
}
