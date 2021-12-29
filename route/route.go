package route

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/middleware"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	jwt2 "github.com/IceWhaleTech/CasaOS/pkg/utils/jwt"
	v1 "github.com/IceWhaleTech/CasaOS/route/v1"
	"github.com/IceWhaleTech/CasaOS/web"
	"github.com/gin-gonic/gin"
)

var swagHandler gin.HandlerFunc
var OnlineDemo bool = false

func InitRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.Cors())
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

	r.GET("/debug", v1.GetSystemConfigDebug)
	//set user
	r.POST("/v1/user/setusernamepwd", v1.Set_Name_Pwd)
	//get user info
	r.GET("/v1/user/info", v1.UserInfo)

	v1Group := r.Group("/v1")

	v1Group.Use(jwt2.JWT(swagHandler))
	{
		v1UserGroup := v1Group.Group("/user")
		v1UserGroup.Use()
		{

			//chang head
			v1UserGroup.POST("/changhead", v1.Up_Load_Head)
			//chang user name
			v1UserGroup.PUT("/changusername", v1.Chang_User_Name)
			//chang pwd
			v1UserGroup.PUT("/changuserpwd", v1.Chang_User_Pwd)
			//edit user info
			v1UserGroup.POST("/changuserinfo", v1.Chang_User_Info)

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

		v1ZeroTierGroup := v1Group.Group("/zerotier")
		v1ZeroTierGroup.Use()
		{
			//获取zerotier token
			v1ZeroTierGroup.POST("/login", v1.ZeroTierGetToken)
			//注册zerotier
			v1ZeroTierGroup.POST("/register", v1.ZeroTierRegister)
			//是否需要登录
			v1ZeroTierGroup.GET("/islogin", v1.ZeroTierIsNeedLogin)
			//获取网络列表
			v1ZeroTierGroup.GET("/list", v1.ZeroTierGetNetworkList)
			//加入网络
			v1ZeroTierGroup.POST("/join/:id", v1.ZeroTierJoinNetwork)
			//离开网络
			v1ZeroTierGroup.POST("/leave/:id", v1.ZeroTierLeaveNetwork)
			//详情
			v1ZeroTierGroup.GET("/info/:id", v1.ZeroTierGetNetworkGetInfo)
			////网络状态
			//v1ZeroTierGroup.GET("/status", v1.ZeroTierGetNetworkGetStatus)
			//修改网络类型
			//v1ZeroTierGroup.PUT("/type/:id", v1.ZeroTierEditType)
			//修改网络类型
			//v1ZeroTierGroup.PUT("/name/:id", v1.ZeroTierEditName)
			//修改v6 assign
			//v1ZeroTierGroup.PUT("/v6assign/:id", v1.ZeroTierEditV6Assign)
			//修改 broadcast
			//v1ZeroTierGroup.PUT("/broadcast/:id", v1.ZeroTierEditBroadcast)
			//create new network
			v1ZeroTierGroup.POST("/create", v1.ZeroTierCreateNetwork)
			//获取用户列表
			v1ZeroTierGroup.GET("/member/:id", v1.ZeroTierMemberList)
			//修改用户信息
			//v1ZeroTierGroup.PUT("/members/:id/auth/:mId", v1.ZeroTierMemberAuth)
			//修改网络用户name
			//v1ZeroTierGroup.PUT("/members/:id/name/:mId", v1.ZeroTierMemberName)
			v1ZeroTierGroup.DELETE("/members/:id/del/:mId", v1.ZeroTierMemberDelete)
			v1ZeroTierGroup.DELETE("/network/:id/del", v1.ZeroTierDeleteNetwork)
			//修改网络用户bridge功能
			//v1ZeroTierGroup.PUT("/members/:id/bridge/:mId", v1.ZeroTierMemberBridge)
			v1ZeroTierGroup.PUT("/edit/:id", v1.ZeroTierEdit)
			v1ZeroTierGroup.GET("/joined/list", v1.ZeroTierJoinedList)
			v1ZeroTierGroup.PUT("/member/:id/edit/:mId", v1.ZeroTierMemberEdit)

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
			v1AppGroup.GET("/mylist", v1.MyAppList)
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
			//准备安装
			//v1AppGroup.GET("/ready/:id", v1.ReadyInstall)
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
			//v1AppGroup.POST("/custom/install", v1.CustomInstallApp)
			v1AppGroup.POST("/share", v1.ShareAppFile)
		}

		v1SysGroup := v1Group.Group("/sys")
		v1SysGroup.Use()
		{
			//获取检查版本是否需要升级
			v1SysGroup.GET("/check", v1.CheckVersion)
			v1SysGroup.POST("/update", v1.SystemUpdate)
			v1SysGroup.GET("/sys", v1.Sys)
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
		}
		v1FileGroup := v1Group.Group("/file")
		v1FileGroup.Use()
		{
			//修改文件名称/目录名称
			v1FileGroup.PUT("/rename", v1.RenamePath)
			v1FileGroup.GET("/read", v1.GetFilerContent)
			v1FileGroup.POST("/upload", v1.PostFileUpload)
			v1FileGroup.GET("/dirpath", v1.DirPath)
			//创建目录
			v1FileGroup.POST("/mkdir", v1.MkdirAll)
			v1FileGroup.POST("/create", v1.PostCreateFile)

			v1FileGroup.GET("/download", v1.GetDownloadFile)
			v1FileGroup.PUT("/move", v1.PutFileMove)
			//v1FileGroup.GET("/download", v1.UserFileDownloadCommonService)
		}
		v1DiskGroup := v1Group.Group("/disk")
		v1DiskGroup.Use()
		{
			v1DiskGroup.GET("/check", v1.GetDiskCheck)
			//获取磁盘列表
			v1DiskGroup.GET("/list", v1.GetPlugInDisk)

			//获取磁盘详情
			v1DiskGroup.GET("/info", v1.GetDiskInfo)

			//格式化磁盘
			v1DiskGroup.POST("/format", v1.FormatDisk)

			//添加分区
			v1DiskGroup.POST("/part", v1.AddPartition)

			//获取可以格式化的内容
			v1DiskGroup.GET("/type", v1.FormatDiskType)

			//删除分区
			v1DiskGroup.DELETE("/delpart", v1.RemovePartition)

			//mount SATA disk
			v1DiskGroup.POST("/mount", v1.PostMountDisk)

			//umount SATA disk
			v1DiskGroup.POST("/umount", v1.PostDiskUmount)
			v1DiskGroup.DELETE("/remove/:id", v1.DeleteDisk)
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
		v1ShortcutsGroup := v1Group.Group("/shortcuts")
		v1ShortcutsGroup.Use()
		{
			v1ShortcutsGroup.GET("/list", v1.GetShortcutsList)
			v1ShortcutsGroup.POST("/add", v1.PostShortcutsAdd)
			v1ShortcutsGroup.PUT("/edit", v1.PutShortcutsEdit)
			v1ShortcutsGroup.DELETE("/del/:id", v1.DeleteShortcutsDelete)
		}
		v1NotifyGroup := v1Group.Group("/notify")
		v1NotifyGroup.Use()
		{
			v1NotifyGroup.GET("/ws", v1.NotifyWS)
			v1NotifyGroup.PUT("/read/:id", v1.PutNotifyRead)
		}
		v1SearchGroup := v1Group.Group("/search")
		v1SearchGroup.Use()
		{
			v1SearchGroup.GET("/search", v1.GetSearchList)
		}
		v1Group.GET("/sync/config", v1.GetSyncConfig)
		v1Group.Any("/syncthing/*url", v1.SyncToSyncthing)

	}
	return r
}
