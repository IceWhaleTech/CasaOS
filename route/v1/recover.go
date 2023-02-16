package v1

import (
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/dropbox"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetRecoverStorage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	t := c.Param("type")
	currentTime := time.Now().UTC()
	currentDate := time.Now().UTC().Format("2006-01-02")
	notify := make(map[string]interface{})
	if t == "GoogleDrive" {
		add := google_drive.Addition{}
		add.Code = c.Query("code")
		if len(add.Code) == 0 {
			c.String(200, `<p>Code cannot be empty</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Code cannot be empty"
			logger.Error("Then code is empty: ", zap.String("code", add.Code), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}

		add.RootFolderID = "root"
		add.ClientID = google_drive.CLIENTID
		add.ClientSecret = google_drive.CLIENTSECRET

		var google_drive google_drive.GoogleDrive
		google_drive.Addition = add
		err := google_drive.Init(c)
		if err != nil {
			c.String(200, `<p>Initialization failure:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Initialization failure"
			logger.Error("Then init error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}

		username, err := google_drive.GetUserInfo(c)
		if err != nil {
			c.String(200, `<p>Failed to get user information:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Failed to get user information"
			logger.Error("Then get user info error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		dmap := make(map[string]string)
		dmap["username"] = username
		configs, err := service.MyService.Storage().GetConfig()
		if err != nil {
			c.String(200, `<p>Failed to get rclone config:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Failed to get rclone config"
			logger.Error("Then get config error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		for _, v := range configs.Remotes {
			cf, err := service.MyService.Storage().GetConfigByName(v)
			if err != nil {
				logger.Error("then get config by name error: ", zap.Error(err), zap.Any("name", v))
				continue
			}
			if cf["type"] == "drive" && cf["username"] == dmap["username"] {
				c.String(200, `<p>The same configuration has been added</p><script>window.close()</script>`)
				err := service.MyService.Storage().CheckAndMountByName(v)
				if err != nil {
					logger.Error("check and mount by name error: ", zap.Error(err), zap.Any("name", cf["username"]))
				}
				notify["status"] = "warn"
				notify["message"] = "The same configuration has been added"
				service.MyService.Notify().SendNotify("casaos:file:recover", notify)
				return
			}
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}

		//username = fileutil.NameAccumulation(username, "/mnt")
		username += "_google_drive_" + strconv.FormatInt(time.Now().Unix(), 10)

		dmap["client_id"] = add.ClientID
		dmap["client_secret"] = add.ClientSecret
		dmap["scope"] = "drive"
		dmap["mount_point"] = "/mnt/" + username
		dmap["token"] = `{"access_token":"` + google_drive.AccessToken + `","token_type":"Bearer","refresh_token":"` + google_drive.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*1).Add(time.Minute*50).Format("15:04:05") + `Z"}`
		service.MyService.Storage().CreateConfig(dmap, username, "drive")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")
		notify := make(map[string]interface{})
		notify["status"] = "success"
		notify["message"] = "Success"
		service.MyService.Notify().SendNotify("casaos:file:recover", notify)
	} else if t == "Dropbox" {
		add := dropbox.Addition{}
		add.Code = c.Query("code")
		if len(add.Code) == 0 {
			c.String(200, `<p>Code cannot be empty</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Code cannot be empty"
			logger.Error("Then code is empty error: ", zap.String("code", add.Code), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		add.RootFolderID = ""
		add.AppKey = dropbox.APPKEY
		add.AppSecret = dropbox.APPSECRET
		var dropbox dropbox.Dropbox
		dropbox.Addition = add
		err := dropbox.Init(c)
		if err != nil {
			c.String(200, `<p>Initialization failure:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Initialization failure"
			logger.Error("Then init error: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		username, err := dropbox.GetUserInfo(c)
		if err != nil {
			c.String(200, `<p>Failed to get user information:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Failed to get user information"
			logger.Error("Then get user information: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		dmap := make(map[string]string)
		dmap["username"] = username

		configs, err := service.MyService.Storage().GetConfig()
		if err != nil {
			c.String(200, `<p>Failed to get rclone config:`+err.Error()+`</p><script>window.close()</script>`)
			notify["status"] = "fail"
			notify["message"] = "Failed to get rclone config"
			logger.Error("Then get config error: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return
		}
		for _, v := range configs.Remotes {
			cf, err := service.MyService.Storage().GetConfigByName(v)
			if err != nil {
				logger.Error("then get config by name error: ", zap.Error(err), zap.Any("name", v))
				continue
			}
			if cf["type"] == "dropbox" && cf["username"] == dmap["username"] {
				c.String(200, `<p>The same configuration has been added</p><script>window.close()</script>`)
				err := service.MyService.Storage().CheckAndMountByName(v)
				if err != nil {
					logger.Error("check and mount by name error: ", zap.Error(err), zap.Any("name", cf["username"]))
				}

				notify["status"] = "warn"
				notify["message"] = "The same configuration has been added"
				service.MyService.Notify().SendNotify("casaos:file:recover", notify)
				return
			}
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}
		username += "_dropbox_" + strconv.FormatInt(time.Now().Unix(), 10)

		dmap["client_id"] = add.AppKey
		dmap["client_secret"] = add.AppSecret
		dmap["token"] = `{"access_token":"` + dropbox.AccessToken + `","token_type":"bearer","refresh_token":"` + dropbox.Addition.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*3).Add(time.Minute*50).Format("15:04:05") + `.780385354Z"}`
		dmap["mount_point"] = "/mnt/" + username
		// data.SetValue(username, "type", "dropbox")
		// data.SetValue(username, "client_id", add.AppKey)
		// data.SetValue(username, "client_secret", add.AppSecret)
		// data.SetValue(username, "mount_point", "/mnt/"+username)

		// data.SetValue(username, "token", `{"access_token":"`+dropbox.AccessToken+`","token_type":"bearer","refresh_token":"`+dropbox.Addition.RefreshToken+`","expiry":"`+currentDate+`T`+currentTime.Add(time.Hour*3).Format("15:04:05")+`.780385354Z"}`)
		// e = data.Save()
		// if e != nil {
		// 	c.String(200, `<p>保存配置失败:`+e.Error()+`</p>`)

		// 	return
		// }
		service.MyService.Storage().CreateConfig(dmap, username, "dropbox")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")

		notify["status"] = "success"
		notify["message"] = "Success"
		service.MyService.Notify().SendNotify("casaos:file:recover", notify)
	}

	c.String(200, `<p>Just close the page</p><script>window.close()</script>`)
}
