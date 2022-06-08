package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
)

//系统信息
type ZiMaService interface {
	GetDiskInfo() *disk.UsageStat

	GetNetState(name string) string
	GetSysInfo() host.InfoStat
	GetDirPath(path string) []model.Path
	GetDirPathOne(path string) (m model.Path)
	MkdirAll(path string) (int, error)
	CreateFile(path string) (int, error)
	RenameFile(oldF, newF string) (int, error)
	GetCpuInfo() []cpu.InfoStat
	GetDeviceTree() string
}

var NetArray [][]model.IOCountersStat

type zima struct {
}

//cpu详情
func (c *zima) GetCpuInfo() []cpu.InfoStat {
	info, _ := cpu.Info()
	return info
}

//获取硬盘详情
func (c *zima) GetDiskInfo() *disk.UsageStat {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}
	diskInfo, _ := disk.Usage(path)
	diskInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.UsedPercent), 64)
	diskInfo.InodesUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.InodesUsedPercent), 64)
	return diskInfo
}

//获取硬盘目录
func (c *zima) GetDirPath(path string) []model.Path {
	if path == "/DATA" {
		sysType := runtime.GOOS
		if sysType == "windows" {
			path = "C:\\CasaOS\\DATA"
		}
		if sysType == "darwin" {
			path = "./CasaOS/DATA"
		}

	}

	ls, _ := ioutil.ReadDir(path)
	dirs := []model.Path{}
	if len(path) > 0 {
		for _, l := range ls {
			filePath := filepath.Join(path, l.Name())
			link, err := filepath.EvalSymlinks(filePath)
			if err != nil {
				link = filePath
			}
			temp := model.Path{Name: l.Name(), Path: filePath, IsDir: l.IsDir(), Date: l.ModTime(), Size: l.Size()}
			if filePath != link {
				file, _ := os.Stat(link)
				temp.IsDir = file.IsDir()
			}
			dirs = append(dirs, temp)
		}
	} else {
		dirs = append(dirs, model.Path{Name: "DATA", Path: "/DATA/", IsDir: true, Date: time.Now()})
	}
	return dirs
}

func (c *zima) GetDirPathOne(path string) (m model.Path) {

	f, err := os.Stat(path)

	if err != nil {
		return
	}
	m.IsDir = f.IsDir()
	m.Name = f.Name()
	m.Path = path
	m.Size = f.Size()
	m.Date = f.ModTime()
	return
}

//获取系统信息
func (c *zima) GetSysInfo() host.InfoStat {
	info, _ := host.Info()
	return *info
}

func (c *zima) GetDeviceTree() string {
	return command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetDeviceTree")
}

//shell脚本参数 { 网卡名称 }
func (c *zima) GetNetState(name string) string {
	return command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;CatNetCardState " + name)
}

//mkdir
func (c *zima) MkdirAll(path string) (int, error) {
	_, err := os.Stat(path)
	if err == nil {
		return oasis_err.DIR_ALREADY_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
			return oasis_err.SUCCESS, nil
		} else if strings.Contains(err.Error(), ": not a directory") {
			return oasis_err.FILE_OR_DIR_EXISTS, err
		}
	}
	return oasis_err.ERROR, err
}

//create
func (c *zima) CreateFile(path string) (int, error) {
	_, err := os.Stat(path)
	if err == nil {
		return oasis_err.FILE_OR_DIR_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			file.CreateFile(path)
			return oasis_err.SUCCESS, nil
		}
	}
	return oasis_err.ERROR, err
}

//修改文件
func (c *zima) RenameFile(oldF, newF string) (int, error) {

	_, err := os.Stat(newF)
	if err == nil {
		return oasis_err.DIR_ALREADY_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			err := os.Rename(oldF, newF)
			if err != nil {
				return oasis_err.ERROR, err
			}
			return oasis_err.SUCCESS, nil
		}
	}
	return oasis_err.ERROR, err
}

//获取zima服务
func NewZiMaService() ZiMaService {
	return &zima{}
}

// func LoopNet() {
// 	netList := MyService.ZiMa().GetNetInfo()

// 	nets := MyService.ZiMa().GetNet(true)
// 	num := 0
// 	for i := 0; i < len(netList); i++ {

// 		for _, netCardName := range nets {

// 			if netList[i].Name == netCardName {
// 				var netArray []model.IOCountersStat
// 				if len(NetArray) < (num + 1) {
// 					netArray = []model.IOCountersStat{}
// 				} else {
// 					netArray = NetArray[num]
// 				}
// 				item := *(*model.IOCountersStat)(unsafe.Pointer(&netList[i]))
// 				item.State = strings.TrimSpace(MyService.ZiMa().GetNetState(netList[i].Name))
// 				item.Time = time.Now().Unix()

// 				if len(netArray) >= 60 {
// 					netArray = netArray[1:]
// 				}
// 				netArray = append(netArray, item)
// 				if len(NetArray) < (num + 1) {
// 					NetArray = append(NetArray, []model.IOCountersStat{})
// 				}

// 				NetArray[num] = netArray

// 				num++
// 				break
// 			}
// 		}

// 	}
// }
