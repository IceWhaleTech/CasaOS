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

var swagHandler gin.HandlerFunc
var OnlineDemo bool = false

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.Use(middleware.Cors())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	gin.SetMode(config.ServerInfo.RunMode)

	r.StaticFS("/ui", http.FS(web.Static))
	r.GET("/", WebUIHome)
	//r.GET("/", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "ui/")
	//})
	if swagHandler != nil {
		r.GET("/swagger/*any", swagHandler)
	}

	r.POST("/v1/user/login", v1.Login)

	r.GET("/v1/guide/check", v1.GetGuideCheck)

	r.GET("/v1/debug", v1.GetSystemConfigDebug)
	//set user
	r.POST("/v1/user/setusernamepwd", v1.Set_Name_Pwd)
	//get user info
	r.GET("/v1/user/info", v1.GetUserInfo)
	//get user info
	r.GET("/v1/person/shareid", v1.GetPersonShareId)
	v1Group := r.Group("/v1")

	v1Group.Use(jwt2.JWT(swagHandler))
	{
		v1UserGroup := v1Group.Group("/user")
		v1UserGroup.Use()
		{

			//chang head
			v1UserGroup.POST("/head", v1.PostUserHead)
			//chang user name
			v1UserGroup.PUT("/username", v1.PutUserName)
			//chang pwd
			v1UserGroup.PUT("/password", v1.PutUserPwd)
			//edit user info
			v1UserGroup.POST("/info", v1.PostUserChangeInfo)
			v1UserGroup.PUT("/nick", v1.PutUserChangeNick)
			v1UserGroup.PUT("/desc", v1.PutUserChangeDesc)
			v1UserGroup.POST("/person/info", v1.PostUserPersonInfo)

			v1UserGroup.GET("/shareid", v1.GetUserShareID)

		}

		v1ZiMaGroup := v1Group.Group("/zima")
		v1ZiMaGroup.Use()
		{
			//获取cpu信息
			v1ZiMaGroup.GET("/getcpuinfo", v1.CupInfo)
			//获取内存信息
			v1ZiMaGroup.GET("/getmeminfo", v1.MemInfo)
			//获取硬盘信息
			v1ZiMaGroup.GET("/getdiskinfo", v1.DiskInfo)

			//获取网络信息
			v1ZiMaGroup.GET("/getnetinfo", v1.NetInfo)

			//获取系统信息
			v1ZiMaGroup.GET("/sysinfo", v1.SysInfo)
		}
		v1DDNSGroup := v1Group.Group("/ddns")
		v1DDNSGroup.Use()
		{
			//获取ddns列表
			v1DDNSGroup.GET("/getlist", v1.DDNSGetDomainList)
			//测试连接性
			v1DDNSGroup.GET("/ping/:api_host", v1.DDNSPing)
			//获取ip
			v1DDNSGroup.GET("/ip", v1.DDNSGetIP)
			//设置ddns
			v1DDNSGroup.POST("/set", v1.DDNSAddConfig)
			//获取ddns
			v1DDNSGroup.GET("/list", v1.DDNSConfigList)
			//获取ddns
			v1DDNSGroup.DELETE("/delete/:id", v1.DDNSDelete)
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
			//分类
			v1AppGroup.GET("/category", v1.CategoryList)
			//容器相关
			v1AppGroup.GET("/terminal/:id", v1.DockerTerminal)
			v1AppGroup.GET("/order", v1.GetAppOrder)
			v1AppGroup.POST("/order", v1.PostAppOrder)
			//app容器详情
			v1AppGroup.GET("/info/:id", v1.ContainerInfo)
			//app容器日志
			v1AppGroup.GET("/logs/:id", v1.ContainerLog)
			//暂停或启动容器
			v1AppGroup.PUT("/state/:id", v1.ChangAppState)
			//安装app
			v1AppGroup.POST("/install/:id", v1.InstallApp)
			//卸载app
			v1AppGroup.DELETE("/uninstall/:id", v1.UnInstallApp)
			//获取安装进度
			v1AppGroup.GET("/speed/:id", v1.GetInstallSpeed)
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
			v1SysGroup.GET("/check", v1.CheckVersion)
			v1SysGroup.GET("/hardware/info", v1.GetSystemHardwareInfo)
			v1SysGroup.GET("/client/version", v1.GetClientVersion)
			v1SysGroup.POST("/update", v1.SystemUpdate)
			v1SysGroup.GET("/wsssh", v1.WsSsh)
			v1SysGroup.GET("/config", v1.GetSystemConfig)
			v1SysGroup.GET("/error/logs", v1.GetCasaOSErrorLogs)
			v1SysGroup.POST("/config", v1.PostSetSystemConfig)
			v1SysGroup.GET("/widget/config", v1.GetWidgetConfig)
			v1SysGroup.POST("/widget/config", v1.PostSetWidgetConfig)
			v1SysGroup.GET("/port", v1.GetCasaOSPort)
			v1SysGroup.PUT("/port", v1.PutCasaOSPort)
			v1SysGroup.POST("/kill", v1.PostKillCasaOS)
			v1SysGroup.GET("/info", v1.Info)
			v1SysGroup.PUT("/usb/off", v1.PutSystemOffUSBAutoMount)
			v1SysGroup.PUT("/usb/on", v1.PutSystemOnUSBAutoMount)
			v1SysGroup.GET("/usb", v1.GetSystemUSBAutoMount)
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
			v1FileGroup.GET("/new/download", v1.GetFileDownloadNew)
			v1FileGroup.POST("/operate", v1.PostOperateFileOrDir)
			v1FileGroup.DELETE("/delete", v1.DeleteFile)
			v1FileGroup.PUT("/update", v1.PutFileContent)
			v1FileGroup.GET("/image", v1.GetFileImage)

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
			v1DiskGroup.POST("/format", v1.FormatDisk)

			// add storage
			v1DiskGroup.POST("/storage", v1.AddPartition)

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
		v1ShareGroup := v1Group.Group("/share")
		v1ShareGroup.Use()
		{
			v1ShareGroup.POST("/add", v1.PostShareDirAdd)
			v1ShareGroup.DELETE("/del/:id", v1.DeleteShareDirDel)
			v1ShareGroup.GET("/list", v1.GetShareDirList)
			v1ShareGroup.GET("/info/:id", v1.GetShareDirInfo)
			v1ShareGroup.PUT("/update/:id", v1.PutShareDirEdit)
		}
		v1TaskGroup := v1Group.Group("/task")
		v1TaskGroup.Use()
		{
			v1TaskGroup.GET("/list", v1.GetTaskList)
			v1TaskGroup.PUT("/update", v1.PutTaskUpdate)
			v1TaskGroup.POST("/add", v1.PostTaskAdd)
			v1TaskGroup.PUT("/completion/:id", v1.PutTaskMarkerCompletion)

		}

		v1NotifyGroup := v1Group.Group("/notify")
		v1NotifyGroup.Use()
		{
			v1NotifyGroup.GET("/ws", v1.NotifyWS)
			v1NotifyGroup.PUT("/read/:id", v1.PutNotifyRead)
		}

		v1PersonGroup := v1Group.Group("/person")
		v1PersonGroup.Use()
		{
			v1PersonGroup.GET("/detection", v1.GetPersonDetection)
			v1PersonGroup.GET("/users", v1.GetPersonFriend)
			v1PersonGroup.POST("/user/:shareids", v1.PostAddPersonFriend)
			v1PersonGroup.DELETE("/user/:shareid", v1.DeletePersonFriend)
			v1PersonGroup.GET("/directory", v1.GetPersonDirectory)
			v1PersonGroup.GET("/file", v1.GetPersonFile)
			v1PersonGroup.GET("/refile/:uuid", v1.GetPersonReFile)
			v1PersonGroup.PUT("/remarks/:shareid", v1.PutPersonRemarks)
			v1PersonGroup.GET("/list", v1.GetPersonDownloadList)
			v1PersonGroup.DELETE("/file/:uuid", v1.DeletePersonDownloadFile)

			v1PersonGroup.POST("/share", v1.PostPersonShare)
			v1PersonGroup.POST("/file/:shareid", v1.PostPersonFile)
			v1PersonGroup.GET("/share", v1.GetPersonShare)
			v1PersonGroup.POST("/down/dir", v1.PostPersonDownDir)
			v1PersonGroup.GET("/down/dir", v1.GetPersonDownDir)
			v1PersonGroup.PUT("/block/:shareid", v1.PutPersonBlock)
			v1PersonGroup.GET("/public", v1.GetPersonPublic)
			v1PersonGroup.PUT("/friend/:shareid", v1.PutPersonAgreeFriend)
			v1PersonGroup.PUT("/write/:shareid", v1.PutPersonWrite)
			v1PersonGroup.GET("/image/thumbnail/:shareid", v1.GetPersonImageThumbnail)

		}
		v1AnalyseGroup := v1Group.Group("/analyse")
		v1AnalyseGroup.Use()
		{
			v1AnalyseGroup.POST("/app", v1.PostAppAnalyse)
		}
		v1Group.GET("/sync/config", v1.GetSyncConfig)
		v1Group.Any("/syncthing/*url", v1.SyncToSyncthing)

	}
	return r
}
