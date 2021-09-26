package service

import (
	"fmt"
	"oasis/model"
	"oasis/pkg/config"
	command2 "oasis/pkg/utils/command"
	"strconv"
)

type SystemService interface {
	UpSystemConfig(systemConfig model.SystemConfig)
	UpdateSystemVersion(version string)
	GetSystemConfigDebug() []string
}
type systemService struct {
}

func (s *systemService) UpdateSystemVersion(version string) {
	fmt.Println("source  " + config.AppInfo.ProjectPath + "/shell/tools.sh ;update " + version)
	command2.OnlyExec("source  " + config.AppInfo.ProjectPath + "/shell/tools.sh ;update " + version)
}
func (s *systemService) GetSystemConfigDebug() []string {
	return command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetSysInfo")
}
func (s *systemService) UpSystemConfig(systemConfig model.SystemConfig) {
	if systemConfig.AutoUpdate != config.SystemConfigInfo.AutoUpdate {
		config.Cfg.Section("system").Key("AutoUpdate").SetValue(strconv.FormatBool(systemConfig.AutoUpdate))
		config.SystemConfigInfo.AutoUpdate = systemConfig.AutoUpdate
	}
	if systemConfig.SearchSwitch != config.SystemConfigInfo.SearchSwitch {
		config.Cfg.Section("system").Key("SearchSwitch").SetValue(strconv.FormatBool(systemConfig.SearchSwitch))
		config.SystemConfigInfo.SearchSwitch = systemConfig.SearchSwitch
	}
	if systemConfig.WidgetsSwitch != config.SystemConfigInfo.WidgetsSwitch {
		config.Cfg.Section("system").Key("WidgetsSwitch").SetValue(strconv.FormatBool(systemConfig.WidgetsSwitch))
		config.SystemConfigInfo.WidgetsSwitch = systemConfig.WidgetsSwitch
	}
	if systemConfig.ShortcutsSwitch != config.SystemConfigInfo.ShortcutsSwitch {
		config.Cfg.Section("system").Key("ShortcutsSwitch").SetValue(strconv.FormatBool(systemConfig.ShortcutsSwitch))
		config.SystemConfigInfo.ShortcutsSwitch = systemConfig.ShortcutsSwitch
	}
	if len(systemConfig.SearchEngine) > 0 && systemConfig.SearchEngine != config.SystemConfigInfo.SearchEngine {
		config.Cfg.Section("system").Key("SearchEngine").SetValue(systemConfig.SearchEngine)
		config.SystemConfigInfo.SearchEngine = systemConfig.SearchEngine
	}
	//	if len(systemConfig.Version) > 0 && systemConfig.Version != config.SystemConfigInfo.Version {
	//	config.Cfg.Section("system").Key("Version").SetValue(systemConfig.Version)
	//	config.SystemConfigInfo.Version = systemConfig.Version
	//}
	if len(systemConfig.Background) > 0 && systemConfig.Background != config.SystemConfigInfo.Background {
		config.Cfg.Section("system").Key("Background").SetValue(systemConfig.Background)
		config.SystemConfigInfo.Background = systemConfig.Background
	}
	if len(systemConfig.BackgroundType) > 0 && systemConfig.BackgroundType != config.SystemConfigInfo.BackgroundType {
		config.Cfg.Section("system").Key("BackgroundType").SetValue(systemConfig.BackgroundType)
		config.SystemConfigInfo.BackgroundType = systemConfig.BackgroundType
	}
	config.Cfg.SaveTo("conf/conf.ini")
}
func NewSystemService() SystemService {
	return &systemService{}
}
