package route

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/samba"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/encryption"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	"github.com/IceWhaleTech/CasaOS/service"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func InitFunction() {
	ShellInit()
	CheckSerialDiskMount()
	CheckToken2_11()
	ImportApplications()
	// Soon to be removed
	ChangeAPIUrl()
	MoveUserToDB()
	go InitNetworkMount()

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
				//if len(h.MountPoint) == 0 && len(v.Children) == 1 && h.FsType == "ext4" {
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
				//}
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

	if len(config.UserInfo.Description) == 0 {
		config.Cfg.Section("user").Key("Description").SetValue("nothing")
		config.UserInfo.Description = "nothing"
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
		user.Username = config.UserInfo.UserName
		user.Email = config.UserInfo.Email
		user.Nickname = config.UserInfo.NickName
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

func InitNetworkMount() {
	time.Sleep(time.Second * 10)
	connections := service.MyService.Connections().GetConnectionsList()
	for _, v := range connections {
		connection := service.MyService.Connections().GetConnectionByID(fmt.Sprint(v.ID))
		directories, err := samba.GetSambaSharesList(connection.Host, connection.Port, connection.Username, connection.Password)
		if err != nil {
			service.MyService.Connections().DeleteConnection(fmt.Sprint(connection.ID))
			loger.Error("mount samba err", zap.Any("err", err), zap.Any("info", connection))
			continue
		}
		baseHostPath := "/mnt/" + connection.Host

		mountPointList := service.MyService.System().GetDirPath(baseHostPath)
		for _, v := range mountPointList {
			service.MyService.Connections().UnmountSmaba(v.Path)
		}

		os.RemoveAll(baseHostPath)

		file.IsNotExistMkDir(baseHostPath)
		for _, v := range directories {
			mountPoint := baseHostPath + "/" + v
			file.IsNotExistMkDir(mountPoint)
			service.MyService.Connections().MountSmaba(connection.Username, connection.Host, v, connection.Port, mountPoint, connection.Password)
		}
		connection.Directories = strings.Join(directories, ",")
		service.MyService.Connections().UpdateConnection(&connection)
	}
}
