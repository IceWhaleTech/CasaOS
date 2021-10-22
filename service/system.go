package service

import (
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
)

type SystemService interface {
	UpSystemConfig(str string, widget string)
	UpdateSystemVersion(version string)
	GetSystemConfigDebug() []string
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
	config.Cfg.SaveTo("conf/conf.ini")
}
func NewSystemService(log loger.OLog) SystemService {
	return &systemService{log: log}
}
