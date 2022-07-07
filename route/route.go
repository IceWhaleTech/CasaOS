// 这是一个用来反馈 API 设计的 PR，不要 merge

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

	// @tiger - 为了方便未来的模块化迭代，前端输出需要独立端口，不要和 API 端口公用。
	r.StaticFS("/ui", http.FS(web.Static))
	r.GET("/", WebUIHome)
	// r.StaticFS("/assets", http.Dir("./static/assets"))
	// r.StaticFile("/favicon.ico", "./static/favicon.ico")
	//r.GET("/", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "ui/")
	//})

	r.POST("/v1/user/register/:key", v1.PostUserRegister)
	r.POST("/v1/user/login", v1.PostUserLogin) //
	r.GET("/v1/user/all/name", v1.GetUserAllUserName)

	r.GET("/v1/sys/init/check", v1.GetSystemInitCheck)
	r.GET("/v1/guide/check", v1.GetGuideCheck)
	r.GET("/v1/debug", v1.GetSystemConfigDebug)

	r.GET("/v1/user/avatar/:id", v1.GetUserAvatar)
	r.GET("/v1/user/image", v1.GetUserImage)

	r.GET("/v1/sys/socket/port", v1.GetSystemSocketPort)
	//r.POST("/v1/user/refresh/token", v1.PostUserRefreshToken)
	v1Group := r.Group("/v1")

	v1Group.Use(jwt2.JWT())
	{
		v1UserGroup := v1Group.Group("/user")
		v1UserGroup.Use()
		{
			v1UserGroup.GET("/info", v1.GetUserInfo)
			v1UserGroup.PUT("/username", v1.PutUserName)
			v1UserGroup.PUT("/password", v1.PutUserPwd)
			v1UserGroup.PUT("/nick", v1.PutUserNick)
			v1UserGroup.PUT("/desc", v1.PutUserDesc)
			v1UserGroup.GET("/info", v1.GetUserInfoByUserName)
			v1UserGroup.GET("/custom/:key", v1.GetUserCustomConf)
			v1UserGroup.POST("/custom/:key", v1.PostUserCustomConf)
			v1UserGroup.DELETE("/custom/:key", v1.DeleteUserCustomConf)
			v1UserGroup.POST("/upload/image/:key", v1.PostUserUploadImage)
			v1UserGroup.POST("/file/image/:key", v1.PostUserFileImage)
			v1UserGroup.DELETE("/image", v1.DeleteUserImage)

			v1UserGroup.PUT("/avatar", v1.PutUserAvatar)
			v1UserGroup.GET("/avatar", v1.GetUserAvatar)
			v1UserGroup.DELETE("/delete/:id", v1.DeleteUser)

		}
		v1AppGroup := v1Group.Group("/app")
		v1AppGroup.Use()
		{
			//获取我的已安装的列表
			v1AppGroup.GET("/my/list", v1.MyAppList)
			//
			v1AppGroup.GET("/usage", v1.AppUsageList)
			//app详情
			v1AppGroup.GET("/appinfo/:id", v1.AppInfo)
			//获取未安装的列表
			v1AppGroup.GET("/list", v1.AppList)
			//获取端口
			v1AppGroup.GET("/port", v1.GetPort)
			//检查端口
			v1AppGroup.GET("/check/:port", v1.PortCheck)

			v1AppGroup.GET("/category", v1.CategoryList)

			v1AppGroup.GET("/terminal/:id", v1.DockerTerminal)
			//app容器详情
			v1AppGroup.GET("/info/:id", v1.ContainerInfo)
			//app容器日志
			v1AppGroup.GET("/logs/:id", v1.ContainerLog)
			//暂停或启动容器
			v1AppGroup.PUT("/state/:id", v1.ChangAppState)
			//安装app
			v1AppGroup.POST("/install", v1.InstallApp)
			//卸载app
			v1AppGroup.DELETE("/uninstall/:id", v1.UnInstallApp)
			//获取进度
			v1AppGroup.GET("/state/:id", v1.GetContainerState)
			//更新容器配置
			v1AppGroup.PUT("/update/:id/setting", v1.UpdateSetting)
			//获取可能新数据
			v1AppGroup.GET("/update/:id/info", v1.ContainerUpdateInfo)
			v1AppGroup.GET("/rely/:id/info", v1.ContainerRelyInfo)
			v1AppGroup.GET("/install/config", v1.GetDockerInstallConfig)
			v1AppGroup.PUT("/update/:id", v1.PutAppUpdate)
			v1AppGroup.POST("/share", v1.ShareAppFile)
		}

		v1SysGroup := v1Group.Group("/sys")
		v1SysGroup.Use()
		{
			v1SysGroup.GET("/version/check", v1.GetSystemCheckVersion)
			v1SysGroup.GET("/hardware/info", v1.GetSystemHardwareInfo)
			v1SysGroup.POST("/update", v1.SystemUpdate)
			v1SysGroup.GET("/wsssh", v1.WsSsh)
			v1SysGroup.GET("/config", v1.GetSystemConfig)
			//v1SysGroup.POST("/config", v1.PostSetSystemConfig)
			v1SysGroup.GET("/error/logs", v1.GetCasaOSErrorLogs)
			v1SysGroup.GET("/widget/config", v1.GetWidgetConfig)
			v1SysGroup.POST("/widget/config", v1.PostSetWidgetConfig)
			v1SysGroup.GET("/port", v1.GetCasaOSPort)
			v1SysGroup.PUT("/port", v1.PutCasaOSPort)
			v1SysGroup.POST("/stop", v1.PostKillCasaOS)
			v1SysGroup.GET("/utilization", v1.GetSystemUtilization)
			v1SysGroup.PUT("/usb/:status", v1.PutSystemUSBAutoMount)
			v1SysGroup.GET("/usb/status", v1.GetSystemUSBAutoMount)
			v1SysGroup.GET("/cpu", v1.GetSystemCupInfo)
			v1SysGroup.GET("/mem", v1.GetSystemMemInfo)
			v1SysGroup.GET("/disk", v1.GetSystemDiskInfo)
			v1SysGroup.GET("/network", v1.GetSystemNetInfo)
		}
		v1FileGroup := v1Group.Group("/file")
		v1FileGroup.Use()
		{
			v1FileGroup.PUT("/rename", v1.RenamePath)
			v1FileGroup.GET("/read", v1.GetFilerContent)
			v1FileGroup.POST("/upload", v1.PostFileUpload)
			v1FileGroup.GET("/upload", v1.GetFileUpload)
			v1FileGroup.GET("/dirpath", v1.DirPath)
			//create folder
			v1FileGroup.POST("/mkdir", v1.MkdirAll)
			v1FileGroup.POST("/create", v1.PostCreateFile)

			v1FileGroup.GET("/download", v1.GetDownloadFile)
			v1FileGroup.GET("/download/*path", v1.GetDownloadSingleFile)
			v1FileGroup.POST("/operate", v1.PostOperateFileOrDir)
			v1FileGroup.DELETE("/delete", v1.DeleteFile)
			v1FileGroup.PUT("/update", v1.PutFileContent)
			v1FileGroup.GET("/image", v1.GetFileImage)
			v1FileGroup.DELETE("/operate/:id", v1.DeleteOperateFileOrDir)
			//v1FileGroup.GET("/download", v1.UserFileDownloadCommonService)
		}
		v1DiskGroup := v1Group.Group("/disk")
		v1DiskGroup.Use()
		{
			v1DiskGroup.GET("/check", v1.GetDiskCheck)

			v1DiskGroup.GET("/list", v1.GetDiskList)

			//获取磁盘详情
			v1DiskGroup.GET("/info", v1.GetDiskInfo)

			//format storage
			v1DiskGroup.POST("/format", v1.PostDiskFormat)

			// add storage
			v1DiskGroup.POST("/storage", v1.PostDiskAddPartition)

			//mount SATA disk
			v1DiskGroup.POST("/mount", v1.PostMountDisk)

			//umount sata disk
			v1DiskGroup.POST("/umount", v1.PostDiskUmount)

			//获取可以格式化的内容
			v1DiskGroup.GET("/type", v1.FormatDiskType)

			//删除分区
			v1DiskGroup.DELETE("/delpart", v1.RemovePartition)
			v1DiskGroup.GET("/usb", v1.GetUSBList)

		}
		v1Group.GET("/sync/config", v1.GetSyncConfig)
	}
	return r
}
