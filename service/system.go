package service

import (
	"bufio"
	"fmt"
	"io"
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

func (s *systemService) GetCasaOSLogs(lineNumber int) string {

	reader, err := file.NewReadLineFromEnd(s.log.Path())
	if err != nil {
		return ""
	}
	defer reader.Close()
	test, err := reader.ReadLine()

	fmt.Println(err)
	fmt.Println(test)
	return string(test)
	file, _ := os.Open(s.log.Path())
	fileScanner := bufio.NewReader(file)
	lineNumber = 5
	lineCount := 1
	var r string
	for i := lineCount; i < lineNumber; i++ {
		line, _, err := fileScanner.ReadLine()
		r += string(line)
		if err == io.EOF {
			return r
		}
		// 如下是某些业务逻辑操作
		// 如下代码打印每次读取的文件行内容

	}
	defer file.Close()
	return r
}

func NewSystemService(log loger.OLog) SystemService {
	return &systemService{log: log}
}
