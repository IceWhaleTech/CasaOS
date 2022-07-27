package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/middleware"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	jwt2 "github.com/IceWhaleTech/CasaOS/pkg/utils/jwt"
	v1 "github.com/IceWhaleTech/CasaOS/route/v1"
	"github.com/IceWhaleTech/CasaOS/web"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var OnlineDemo bool = false

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.Use(middleware.Cors())
	r.Use(middleware.WriteLog())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	gin.SetMode(config.ServerInfo.RunMode)

	r.StaticFS("/ui", http.FS(web.Static))
	r.GET("/", WebUIHome)
	// r.StaticFS("/assets", http.Dir("./static/assets"))
	// r.StaticFile("/favicon.ico", "./static/favicon.ico")
	//r.GET("/", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "ui/")
	//})

	r.POST("/v1/users/register", v1.PostUserRegister)
	r.POST("/v1/users/login", v1.PostUserLogin)
	r.GET("/v1/users/name", v1.GetUserAllUsername) //all/name
	r.POST("/v1/users/refresh", v1.PostUserRefreshToken)
	// No short-term modifications
	r.GET("/v1/users/image", v1.GetUserImage)

	r.GET("/v1/users/status", v1.GetUserStatus) //init/check
	//r.GET("/v1/guide/check", v1.GetGuideCheck)         // /v1/sys/guide_check
	r.GET("/v1/sys/debug", v1.GetSystemConfigDebug) // //debug

	r.GET("/v1/sys/socket-port", v1.GetSystemSocketPort) //sys/socket_port

	v1Group := r.Group("/v1")

	v1Group.Use(jwt2.JWT())
	{
		v1UsersGroup := v1Group.Group("/users")
		v1UsersGroup.Use()
		{
			v1UsersGroup.GET("/current", v1.GetUserInfo)
			v1UsersGroup.PUT("/current", v1.PutUserInfo)
			v1UsersGroup.PUT("/current/password", v1.PutUserPassword)

			v1UsersGroup.GET("/current/custom/:key", v1.GetUserCustomConf)
			v1UsersGroup.POST("/current/custom/:key", v1.PostUserCustomConf)
			v1UsersGroup.DELETE("/current/custom/:key", v1.DeleteUserCustomConf)

			v1UsersGroup.POST("/current/image/:key", v1.PostUserUploadImage)
			v1UsersGroup.PUT("/current/image/:key", v1.PutUserImage)
			//v1UserGroup.POST("/file/image/:key", v1.PostUserFileImage)
			v1UsersGroup.DELETE("/current/image", v1.DeleteUserImage)

			//v1UserGroup.PUT("/avatar", v1.PutUserAvatar)
			//v1UserGroup.GET("/avatar", v1.GetUserAvatar)
			v1UsersGroup.DELETE("/:id", v1.DeleteUser)
			v1UsersGroup.GET("/:username", v1.GetUserInfoByUsername)
			v1UsersGroup.DELETE("", v1.DeleteUserAll)
		}

		v1AppsGroup := v1Group.Group("/apps")
		v1AppsGroup.Use()
		{
			v1AppsGroup.GET("", v1.AppList) //list
			v1AppsGroup.GET("/:id", v1.AppInfo)
		}
		v1ContainerGroup := v1Group.Group("/container")
		v1ContainerGroup.Use()
		{
			v1ContainerGroup.GET("", v1.MyAppList) ///my/list
			v1ContainerGroup.GET("/usage", v1.AppUsageList)
			v1ContainerGroup.GET("/:id", v1.ContainerUpdateInfo)    ///update/:id/info
			v1ContainerGroup.GET("/:id/logs", v1.ContainerLog)      // /app/logs/:id
			v1ContainerGroup.GET("/networks", v1.GetDockerNetworks) //app/install/config

			v1ContainerGroup.GET("/:id/state", v1.GetContainerState) //app/state/:id ?state=install_progress
			// there are problems, temporarily do not deal with
			v1ContainerGroup.GET("/:id/terminal", v1.DockerTerminal) //app/terminal/:id
			v1ContainerGroup.POST("", v1.InstallApp)                 //app/install
			//v1ContainerGroup.GET("/:id", v1.ContainerInfo) // /app/info/:id

			v1ContainerGroup.PUT("/:id", v1.UpdateSetting) ///update/:id/setting

			v1ContainerGroup.PUT("/:id/state", v1.ChangAppState) // /app/state/:id
			v1ContainerGroup.DELETE("/:id", v1.UnInstallApp)     //app/uninstall/:id
			//Not used
			v1ContainerGroup.PUT("/:id/latest", v1.PutAppUpdate)
			//Not used
			v1ContainerGroup.POST("/share", v1.ShareAppFile)

		}
		v1AppCategoriesGroup := v1Group.Group("/app-categories")
		v1AppCategoriesGroup.Use()
		{
			v1AppCategoriesGroup.GET("", v1.CategoryList)
		}

		v1SysGroup := v1Group.Group("/sys")
		v1SysGroup.Use()
		{
			v1SysGroup.GET("/version", v1.GetSystemCheckVersion) //version/check
			v1SysGroup.POST("/update", v1.SystemUpdate)

			v1SysGroup.GET("/hardware", v1.GetSystemHardwareInfo) //hardware/info

			v1SysGroup.GET("/wsssh", v1.WsSsh)
			v1SysGroup.POST("/ssh-login", v1.PostSshLogin)
			//v1SysGroup.GET("/config", v1.GetSystemConfig) //delete
			//v1SysGroup.POST("/config", v1.PostSetSystemConfig)
			v1SysGroup.GET("/logs", v1.GetCasaOSErrorLogs) //error/logs
			//v1SysGroup.GET("/widget/config", v1.GetWidgetConfig)//delete
			//v1SysGroup.POST("/widget/config", v1.PostSetWidgetConfig)//delete

			v1SysGroup.POST("/stop", v1.PostKillCasaOS)

			v1SysGroup.GET("/utilization", v1.GetSystemUtilization)
			// v1SysGroup.GET("/cpu", v1.GetSystemCupInfo)
			// v1SysGroup.GET("/mem", v1.GetSystemMemInfo)
			// v1SysGroup.GET("/disk", v1.GetSystemDiskInfo)
			// v1SysGroup.GET("/network", v1.GetSystemNetInfo)

			v1SysGroup.PUT("/usb-auto-mount", v1.PutSystemUSBAutoMount) ///sys/usb/:status
			v1SysGroup.GET("/usb-auto-mount", v1.GetSystemUSBAutoMount) ///sys/usb/status

			v1SysGroup.GET("/server-info", nil)
			v1SysGroup.PUT("/server-info", nil)
			v1SysGroup.GET("/apps-state", v1.GetSystemAppsStatus)
			v1SysGroup.GET("/port", v1.GetCasaOSPort)
			v1SysGroup.PUT("/port", v1.PutCasaOSPort)
		}
		v1PortGroup := v1Group.Group("/port")
		v1PortGroup.Use()
		{
			v1PortGroup.GET("/", v1.GetPort)              //app/port
			v1PortGroup.GET("/state/:port", v1.PortCheck) //app/check/:port
		}

		v1FileGroup := v1Group.Group("/file")
		v1FileGroup.Use()
		{
			v1FileGroup.GET("", v1.GetDownloadSingleFile) //download/:path
			v1FileGroup.POST("", v1.PostCreateFile)
			v1FileGroup.PUT("", v1.PutFileContent)
			v1FileGroup.PUT("/name", v1.RenamePath)
			//file/rename
			v1FileGroup.GET("/content", v1.GetFilerContent) //file/read

			//File uploads need to be handled separately, and will not be modified here
			v1FileGroup.POST("/upload", v1.PostFileUpload)
			v1FileGroup.GET("/upload", v1.GetFileUpload)
			//v1FileGroup.GET("/download", v1.UserFileDownloadCommonService)
		}
		v1FolderGroup := v1Group.Group("/folder")
		v1FolderGroup.Use()
		{
			v1FolderGroup.PUT("/name", v1.RenamePath)
			v1FolderGroup.GET("", v1.DirPath)   ///file/dirpath
			v1FolderGroup.POST("", v1.MkdirAll) ///file/mkdir
		}
		v1BatchGroup := v1Group.Group("/batch")
		v1BatchGroup.Use()
		{

			v1BatchGroup.DELETE("", v1.DeleteFile) //file/delete
			v1BatchGroup.DELETE("/:id/task", v1.DeleteOperateFileOrDir)
			v1BatchGroup.POST("/task", v1.PostOperateFileOrDir) //file/operate
			v1BatchGroup.GET("", v1.GetDownloadFile)
		}
		v1ImageGroup := v1Group.Group("/image")
		v1ImageGroup.Use()
		{
			v1ImageGroup.GET("", v1.GetFileImage)
		}

		v1DisksGroup := v1Group.Group("/disks")
		v1DisksGroup.Use()
		{
			//v1DiskGroup.GET("/check", v1.GetDiskCheck) //delete
			//v1DisksGroup.GET("", v1.GetDiskInfo)

			//v1DisksGroup.POST("", v1.PostMountDisk)
			v1DisksGroup.GET("", v1.GetDiskList)
			// //format storage
			// v1DiskGroup.POST("/format", v1.PostDiskFormat)

			// //mount SATA disk
			// v1DiskGroup.POST("/mount", v1.PostMountDisk)

			// //umount sata disk
			// v1DiskGroup.POST("/umount", v1.PostDiskUmount)

			//v1DiskGroup.GET("/type", v1.FormatDiskType)//delete

			v1DisksGroup.DELETE("/part", v1.RemovePartition) //disk/delpart
		}

		v1StorageGroup := v1Group.Group("/storage")
		v1StorageGroup.Use()
		{
			v1StorageGroup.POST("", v1.PostDiskAddPartition)

			v1StorageGroup.PUT("", v1.PostDiskFormat)

			v1StorageGroup.DELETE("", v1.PostDiskUmount)
		}
		v1SambaGroup := v1Group.Group("/samba")
		v1SambaGroup.Use()
		{
			v1ConnectionsGroup := v1SambaGroup.Group("/connections")
			v1ConnectionsGroup.Use()
			{
				v1ConnectionsGroup.GET("", v1.GetSambaConnectionsList)
				v1ConnectionsGroup.POST("", v1.PostSambaConnectionsCreate)
				v1ConnectionsGroup.DELETE("/:id", v1.DeleteSambaConnections)
			}
			v1SharesGroup := v1SambaGroup.Group("/shares")
			v1SharesGroup.Use()
			{
				v1SharesGroup.GET("", v1.GetSambaSharesList)
				v1SharesGroup.POST("", v1.PostSambaSharesCreate)
				v1SharesGroup.DELETE("/:id", v1.DeleteSambaShares)
			}
		}
	}
	return r
}
