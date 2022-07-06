package route

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	uuid "github.com/satori/go.uuid"
)

func InitFunction() {

	ShellInit()
	CheckSerialDiskMount()

	CheckToken2_11()
	ImportApplications()
	ChangeAPIUrl()

	MoveUserToDB()
}

func CheckSerialDiskMount() {
	// check mount point
	dbList := service.MyService.Disk().GetSerialAll()

	list := service.MyService.Disk().LSBLK(true)
	mountPoint := make(map[string]string, len(dbList))
	//remount
	for _, v := range dbList {
		mountPoint[v.UUID] = v.MountPoint
	}
	for _, v := range list {
		command.ExecEnabledSMART(v.Path)
		if v.Children != nil {
			for _, h := range v.Children {
				if len(h.MountPoint) == 0 && len(v.Children) == 1 && h.FsType == "ext4" {
					if m, ok := mountPoint[h.UUID]; ok {
						//mount point check
						volume := m
						if !file.CheckNotExist(m) {
							for i := 0; file.CheckNotExist(volume); i++ {
								volume = m + strconv.Itoa(i+1)
							}
						}
						service.MyService.Disk().MountDisk(h.Path, volume)
						if volume != m {
							ms := model2.SerialDisk{}
							ms.UUID = v.UUID
							ms.MountPoint = volume
							service.MyService.Disk().UpdateMountPoint(ms)
						}

					}
				}
			}
		}
	}
	service.MyService.Disk().RemoveLSBLKCache()
	command.OnlyExec("source " + config.AppInfo.ShellPath + "/helper.sh ;AutoRemoveUnuseDir")
}
func ShellInit() {
	command.OnlyExec("curl -fsSL https://raw.githubusercontent.com/IceWhaleTech/get/main/assist.sh | bash")
	if !file.CheckNotExist("/casaOS") {
		command.OnlyExec("source /casaOS/server/shell/update.sh ;")
		command.OnlyExec("source " + config.AppInfo.ShellPath + "/delete-old-service.sh ;")
	}

}
func CheckToken2_11() {
	if len(config.ServerInfo.Token) == 0 {
		token := uuid.NewV4().String
		config.ServerInfo.Token = token()
		config.Cfg.Section("server").Key("Token").SetValue(token())
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}

	sysType := runtime.GOOS
	if len(config.FileSettingInfo.DownloadDir) == 0 {
		downloadPath := "/DATA/Downloads"
		if sysType == "windows" {
			downloadPath = "C:\\CasaOS\\DATA\\Downloads"
		}
		if sysType == "darwin" {
			downloadPath = "./CasaOS/DATA/Downloads"
		}
		config.Cfg.Section("file").Key("DownloadDir").SetValue(downloadPath)
		config.FileSettingInfo.DownloadDir = downloadPath
		file.IsNotExistMkDir(config.FileSettingInfo.DownloadDir)
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}

	if len(config.UserInfo.Description) == 0 {
		config.Cfg.Section("user").Key("Description").SetValue("nothing")
		config.UserInfo.Description = "nothing"
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}
	if len(config.ServerInfo.Handshake) == 0 {
		config.Cfg.Section("server").Key("Handshake").SetValue("socket.casaos.io")
		config.ServerInfo.Handshake = "socket.casaos.io"
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}

	if service.MyService.System().GetSysInfo().KernelArch == "aarch64" && config.ServerInfo.USBAutoMount != "True" && strings.Contains(service.MyService.System().GetDeviceTree(), "Raspberry Pi") {
		service.MyService.System().UpdateUSBAutoMount("False")
		service.MyService.System().ExecUSBAutoMountShell("False")
	}

	// str := []string{}
	// str = append(str, "ddd")
	// str = append(str, "aaa")
	// ddd := strings.Join(str, "|")
	// config.Cfg.Section("file").Key("ShareDir").SetValue(ddd)

	// config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)

}

func ImportApplications() {
	service.MyService.App().ImportApplications(true)
}

// 0.3.1
func ChangeAPIUrl() {

	newAPIUrl := "https://api.casaos.io/casaos-api"
	if config.ServerInfo.ServerApi == "https://api.casaos.zimaboard.com" {
		config.ServerInfo.ServerApi = newAPIUrl
		config.Cfg.Section("server").Key("ServerApi").SetValue(newAPIUrl)
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}

}

//0.3.3
//Transferring user data to the database
func MoveUserToDB() {

	if len(config.UserInfo.UserName) > 0 && service.MyService.User().GetUserInfoByUserName(config.UserInfo.UserName).Id == 0 {
		user := model2.UserDBModel{}
		user.UserName = config.UserInfo.UserName
		user.Email = config.UserInfo.Email
		user.NickName = config.UserInfo.NickName
		user.Password = encryption.GetMD5ByStr(config.UserInfo.PWD)
		user.Role = "admin"
		user = service.MyService.User().CreateUser(user)
		if user.Id > 0 {
			userPath := config.AppInfo.UserDataPath + "/" + strconv.Itoa(user.Id)
			file.MkDir(userPath)
			os.Rename("/casaOS/server/conf/app_order.json", userPath+"/app_order.json")
		}

	}
}
