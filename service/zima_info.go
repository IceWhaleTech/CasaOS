package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

//系统信息
type ZiMaService interface {
	GetCpuPercent() float64
	GetCpuCoreNum() int
	GetMemInfo() *mem.VirtualMemoryStat
	GetDiskInfo() *disk.UsageStat
	GetNetInfo() []net.IOCountersStat
	GetNet(physics bool) []string
	GetNetState(name string) string
	GetSysInfo() host.InfoStat
	GetDirPath(path string) []model.Path
	MkdirAll(path string) (int, error)
	CreateFile(path string) (int, error)
	RenameFile(oldF, newF string) (int, error)
	GetCpuInfo() []cpu.InfoStat
}

type zima struct {
}

//获取cpu占用率
func (c *zima) GetCpuPercent() float64 {
	percent, _ := cpu.Percent(0, false)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", percent[0]), 64)
	return value
}

//获取物理核心数
func (c *zima) GetCpuCoreNum() int {
	count, _ := cpu.Counts(false)
	return count
}

//cpu详情
func (c *zima) GetCpuInfo() []cpu.InfoStat {
	info, _ := cpu.Info()
	return info
}

//获取内存详情
func (c *zima) GetMemInfo() *mem.VirtualMemoryStat {
	memInfo, _ := mem.VirtualMemory()
	memInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memInfo.UsedPercent), 64)
	return memInfo
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

	ls, _ := ioutil.ReadDir(path)
	dirs := []model.Path{}

	if strings.Count(path, "/") > 1 {
		for _, l := range ls {
			pathTemp := path + l.Name()
			if l.IsDir() {
				pathTemp += "/"
			}

			dirs = append(dirs, model.Path{Name: l.Name(), Path: pathTemp, IsDir: l.IsDir()})
		}
	} else {
		dirs = append(dirs, model.Path{Name: "DATA", Path: "/DATA/", IsDir: true})
	}
	return dirs
}

//获取系统信息
func (c *zima) GetSysInfo() host.InfoStat {
	info, _ := host.Info()
	return *info
}

//shell脚本参数 {1:虚拟网卡  2:物理网卡}
func (c *zima) GetNet(physics bool) []string {
	t := "1"
	if physics {
		t = "2"
	}
	return command2.ExecResultStrArray("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetNetCard " + t)
}

//shell脚本参数 { 网卡名称 }
func (c *zima) GetNetState(name string) string {
	return command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;CatNetCardState " + name)
}

//网络信息
func (c *zima) GetNetInfo() []net.IOCountersStat {
	//loger.Error("输出个内容试试")
	parts, _ := net.IOCounters(true)
	//fmt.Println(net.ConntrackStatsWithContext(true))
	return parts
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
