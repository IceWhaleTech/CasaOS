package service

import (
	"fmt"
	"io/ioutil"
	net2 "net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemService interface {
	UpdateSystemVersion(version string)
	GetSystemConfigDebug() []string
	GetCasaOSLogs(lineNumber int) string
	UpdateAssist()
	UpSystemPort(port string)
	GetTimeZone() string
	UpAppOrderFile(str, id string)
	GetAppOrderFile(id string) []byte
	GetNet(physics bool) []string
	GetNetInfo() []net.IOCountersStat
	GetCpuCoreNum() int
	GetCpuPercent() float64
	GetMemInfo() map[string]interface{}
	GetCpuInfo() []cpu.InfoStat
	GetDirPath(path string) []model.Path
	GetDirPathOne(path string) (m model.Path)
	GetNetState(name string) string
	GetDiskInfo() *disk.UsageStat
	GetSysInfo() host.InfoStat
	GetDeviceTree() string
	CreateFile(path string) (int, error)
	RenameFile(oldF, newF string) (int, error)
	MkdirAll(path string) (int, error)
	IsServiceRunning(name string) bool
	GetCPUTemperature() int
	GetCPUPower() map[string]string
	GetMacAddress() (string, error)
	SystemReboot() error
	SystemShutdown() error
}
type systemService struct{}

func (c *systemService) GetMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	inter := interfaces[0]
	return inter.HardwareAddr, nil
}

func (c *systemService) MkdirAll(path string) (int, error) {
	_, err := os.Stat(path)
	if err == nil {
		return common_err.DIR_ALREADY_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
			return common_err.SUCCESS, nil
		} else if strings.Contains(err.Error(), ": not a directory") {
			return common_err.FILE_OR_DIR_EXISTS, err
		}
	}
	return common_err.SERVICE_ERROR, err
}

func (c *systemService) RenameFile(oldF, newF string) (int, error) {
	_, err := os.Stat(newF)
	if err == nil {
		return common_err.DIR_ALREADY_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			err := os.Rename(oldF, newF)
			if err != nil {
				return common_err.SERVICE_ERROR, err
			}
			return common_err.SUCCESS, nil
		}
	}
	return common_err.SERVICE_ERROR, err
}

func (c *systemService) CreateFile(path string) (int, error) {
	_, err := os.Stat(path)
	if err == nil {
		return common_err.FILE_OR_DIR_EXISTS, nil
	} else {
		if os.IsNotExist(err) {
			file.CreateFile(path)
			return common_err.SUCCESS, nil
		}
	}
	return common_err.SERVICE_ERROR, err
}

func (c *systemService) GetDeviceTree() string {
	return command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;GetDeviceTree")
}

func (c *systemService) GetSysInfo() host.InfoStat {
	info, _ := host.Info()
	return *info
}

func (c *systemService) GetDiskInfo() *disk.UsageStat {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}
	diskInfo, _ := disk.Usage(path)
	diskInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.UsedPercent), 64)
	diskInfo.InodesUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskInfo.InodesUsedPercent), 64)
	return diskInfo
}

func (c *systemService) GetNetState(name string) string {
	return command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;CatNetCardState " + name)
}

func (c *systemService) GetDirPathOne(path string) (m model.Path) {
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

func (c *systemService) GetDirPath(path string) []model.Path {
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

func (c *systemService) GetCpuInfo() []cpu.InfoStat {
	info, _ := cpu.Info()
	return info
}

func (c *systemService) GetMemInfo() map[string]interface{} {
	memInfo, _ := mem.VirtualMemory()
	memInfo.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memInfo.UsedPercent), 64)
	memData := make(map[string]interface{})
	memData["total"] = memInfo.Total
	memData["available"] = memInfo.Available
	memData["used"] = memInfo.Used
	memData["free"] = memInfo.Free
	memData["usedPercent"] = memInfo.UsedPercent
	return memData
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
	return parts
}

func (c *systemService) GetNet(physics bool) []string {
	t := "1"
	if physics {
		t = "2"
	}
	return command2.ExecResultStrArray("source " + config.AppInfo.ShellPath + "/helper.sh ;GetNetCard " + t)
}

func (s *systemService) UpdateSystemVersion(version string) {
	if file.Exists(config.AppInfo.LogPath + "/upgrade.log") {
		os.Remove(config.AppInfo.LogPath + "/upgrade.log")
	}
	file.CreateFile(config.AppInfo.LogPath + "/upgrade.log")
	// go command2.OnlyExec("curl -fsSL https://raw.githubusercontent.com/LinkLeong/casaos-alpha/main/update.sh | bash")
	if len(config.ServerInfo.UpdateUrl) > 0 {
		go command2.OnlyExec("curl -fsSL " + config.ServerInfo.UpdateUrl + " | bash")
	} else {
		go command2.OnlyExec("curl -fsSL https://get.casaos.io/update | bash")
	}

	// s.log.Error(config.AppInfo.ProjectPath + "/shell/tool.sh -r " + version)
	// s.log.Error(command2.ExecResultStr(config.AppInfo.ProjectPath + "/shell/tool.sh -r " + version))
}

func (s *systemService) UpdateAssist() {
	command2.ExecResultStrArray("source " + config.AppInfo.ShellPath + "/assist.sh")
}

func (s *systemService) GetTimeZone() string {
	return command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;GetTimeZone")
}

func (s *systemService) GetSystemConfigDebug() []string {
	return command2.ExecResultStrArray("source " + config.AppInfo.ShellPath + "/helper.sh ;GetSysInfo")
}

func (s *systemService) UpAppOrderFile(str, id string) {
	file.WriteToPath([]byte(str), config.AppInfo.DBPath+"/"+id, "app_order.json")
}

func (s *systemService) GetAppOrderFile(id string) []byte {
	return file.ReadFullFile(config.AppInfo.UserDataPath + "/" + id + "/app_order.json")
}

func (s *systemService) UpSystemPort(port string) {
	if len(port) > 0 && port != config.ServerInfo.HttpPort {
		config.Cfg.Section("server").Key("HttpPort").SetValue(port)
		config.ServerInfo.HttpPort = port
	}
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
}

func (s *systemService) GetCasaOSLogs(lineNumber int) string {
	file, err := os.Open(filepath.Join(config.AppInfo.LogPath, fmt.Sprintf("%s.%s",
		config.AppInfo.LogSaveName,
		config.AppInfo.LogFileExt,
	)))
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

func (s *systemService) IsServiceRunning(name string) bool {
	status := command2.ExecResultStr("source " + config.AppInfo.ShellPath + "/helper.sh ;CheckServiceStatus smbd")
	return strings.TrimSpace(status) == "running"
}

// find thermal_zone of cpu.
// assertions:
//   - thermal_zone "type" and "temp" are required fields
//     (https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-class-thermal)
func GetCPUThermalZone() string {
	keyName := "cpu_thermal_zone"

	var path string
	if result, ok := Cache.Get(keyName); ok {
		path, ok = result.(string)
		if ok {
			return path
		}
	}

	var name string
	cpu_types := []string{"x86_pkg_temp", "cpu", "CPU", "soc"}
	stub := "/sys/devices/virtual/thermal/thermal_zone"
	for i := 0; i < 100; i++ {
		path = stub + strconv.Itoa(i)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			name = strings.TrimSuffix(string(file.ReadFullFile(path+"/type")), "\n")
			for _, s := range cpu_types {
				if strings.HasPrefix(name, s) {
					logger.Info(fmt.Sprintf("CPU thermal zone found: %s, path: %s.", name, path))
					Cache.SetDefault(keyName, path)
					return path
				}
			}
		} else {
			if len(name) > 0 { //proves at least one zone
				path = stub + "0"
			} else {
				path = ""
			}
			break
		}
	}
	Cache.SetDefault(keyName, path)
	return path
}

func (s *systemService) GetCPUTemperature() int {
	outPut := ""
	path := GetCPUThermalZone()
	if len(path) > 0 {
		outPut = string(file.ReadFullFile(path + "/temp"))
	} else {
		outPut = "0"
	}

	celsius, _ := strconv.Atoi(strings.TrimSpace(outPut))

	if celsius > 1000 {
		celsius = celsius / 1000
	}
	return celsius
}

func (s *systemService) GetCPUPower() map[string]string {
	data := make(map[string]string, 2)
	data["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	if file.Exists("/sys/class/powercap/intel-rapl/intel-rapl:0/energy_uj") {
		data["value"] = strings.TrimSpace(string(file.ReadFullFile("/sys/class/powercap/intel-rapl/intel-rapl:0/energy_uj")))
	} else {
		data["value"] = "0"
	}
	return data
}

func (s *systemService) SystemReboot() error {
	//cmd := exec.Command("/bin/bash", "-c", "reboot")
	arg := []string{"6"}
	cmd := exec.Command("init", arg...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
func (s *systemService) SystemShutdown() error {
	arg := []string{"0"}
	cmd := exec.Command("init", arg...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func NewSystemService() SystemService {

	return &systemService{}
}
