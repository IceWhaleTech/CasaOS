package v1

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/dropbox"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/drivers/onedrive"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func GetRecoverStorage(ctx echo.Context) error {
	t := strings.TrimSuffix(ctx.Param("type"), "/")
	currentTime := time.Now().UTC()
	currentDate := time.Now().UTC().Format("2006-01-02")
	notify := make(map[string]interface{})
	event := "casaos:file:recover"
	if t == "GoogleDrive" {
		google_drive := google_drive.GetConfig()
		google_drive.Code = ctx.QueryParam("code")
		if len(google_drive.Code) == 0 {
			notify["status"] = "fail"
			notify["message"] = "Code cannot be empty"
			logger.Error("Then code is empty: ", zap.String("code", google_drive.Code), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify("casaos:file:recover", notify)
			return ctx.HTML(http.StatusOK, `<p>Code cannot be empty</p><script>window.close()</script>`)
		}
		err := google_drive.Init(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Initialization failure"
			logger.Error("Then init error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Initialization failure:`+err.Error()+`</p><script>window.close()</script>`)
		}

		username, err := google_drive.GetUserInfo(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get user information"
			logger.Error("Then get user info error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get user information:`+err.Error()+`</p><script>window.close()</script>`)
		}
		dmap := make(map[string]string)
		dmap["username"] = username
		configs, err := service.MyService.Storage().GetConfig()
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get rclone config"
			logger.Error("Then get config error: ", zap.Error(err), zap.Any("name", "google_drive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get rclone config:`+err.Error()+`</p><script>window.close()</script>`)
		}
		for _, v := range configs.Remotes {
			cf, err := service.MyService.Storage().GetConfigByName(v)
			if err != nil {
				logger.Error("then get config by name error: ", zap.Error(err), zap.Any("name", v))
				continue
			}
			if cf["type"] == "drive" && cf["username"] == dmap["username"] {
				err := service.MyService.Storage().CheckAndMountByName(v)
				if err != nil {
					logger.Error("check and mount by name error: ", zap.Error(err), zap.Any("name", cf["username"]))
				}
				notify["status"] = "warn"
				notify["message"] = "The same configuration has been added"
				service.MyService.Notify().SendNotify(event, notify)
				return ctx.HTML(http.StatusOK, `<p>The same configuration has been added</p><script>window.close()</script>`)
			}
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}

		username += "_google_drive_" + strconv.FormatInt(time.Now().Unix(), 10)

		dmap["client_id"] = google_drive.ClientID
		dmap["client_secret"] = google_drive.ClientSecret
		dmap["scope"] = "drive"
		dmap["mount_point"] = "/mnt/" + username
		dmap["token"] = `{"access_token":"` + google_drive.AccessToken + `","token_type":"Bearer","refresh_token":"` + google_drive.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*1).Add(time.Minute*50).Format("15:04:05") + `Z"}`
		service.MyService.Storage().CreateConfig(dmap, username, "drive")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")
		notify := make(map[string]interface{})
		notify["status"] = "success"
		notify["message"] = "Success"
		notify["driver"] = "GoogleDrive"
		service.MyService.Notify().SendNotify(event, notify)
	} else if t == "Dropbox" {
		dropbox := dropbox.GetConfig()

		dropbox.Code = ctx.QueryParam("code")
		if len(dropbox.Code) == 0 {

			notify["status"] = "fail"
			notify["message"] = "Code cannot be empty"
			logger.Error("Then code is empty error: ", zap.String("code", dropbox.Code), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Code cannot be empty</p><script>window.close()</script>`)
		}
		err := dropbox.Init(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Initialization failure"
			logger.Error("Then init error: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Initialization failure:`+err.Error()+`</p><script>window.close()</script>`)
		}
		username, err := dropbox.GetUserInfo(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get user information"
			logger.Error("Then get user information: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get user information:`+err.Error()+`</p><script>window.close()</script>`)
		}
		dmap := make(map[string]string)
		dmap["username"] = username

		configs, err := service.MyService.Storage().GetConfig()
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get rclone config"
			logger.Error("Then get config error: ", zap.Error(err), zap.Any("name", "dropbox"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get rclone config:`+err.Error()+`</p><script>window.close()</script>`)
		}
		for _, v := range configs.Remotes {
			cf, err := service.MyService.Storage().GetConfigByName(v)
			if err != nil {
				logger.Error("then get config by name error: ", zap.Error(err), zap.Any("name", v))
				continue
			}
			if cf["type"] == "dropbox" && cf["username"] == dmap["username"] {
				if err := service.MyService.Storage().CheckAndMountByName(v); err != nil {
					logger.Error("check and mount by name error: ", zap.Error(err), zap.Any("name", cf["username"]))
				}

				notify["status"] = "warn"
				notify["message"] = "The same configuration has been added"
				service.MyService.Notify().SendNotify(event, notify)
				return ctx.HTML(http.StatusOK, `<p>The same configuration has been added</p><script>window.close()</script>`)
			}
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}
		username += "_dropbox_" + strconv.FormatInt(time.Now().Unix(), 10)

		dmap["client_id"] = dropbox.AppKey
		dmap["client_secret"] = dropbox.AppSecret
		dmap["token"] = `{"access_token":"` + dropbox.AccessToken + `","token_type":"bearer","refresh_token":"` + dropbox.Addition.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*3).Add(time.Minute*50).Format("15:04:05") + `.780385354Z"}`
		dmap["mount_point"] = "/mnt/" + username
		service.MyService.Storage().CreateConfig(dmap, username, "dropbox")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")

		notify["status"] = "success"
		notify["message"] = "Success"
		notify["driver"] = "Dropbox"
		service.MyService.Notify().SendNotify(event, notify)
	} else if t == "Onedrive" {
		onedrive := onedrive.GetConfig()
		onedrive.Code = ctx.QueryParam("code")
		if len(onedrive.Code) == 0 {
			notify["status"] = "fail"
			notify["message"] = "Code cannot be empty"
			logger.Error("Then code is empty error: ", zap.String("code", onedrive.Code), zap.Any("name", "onedrive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Code cannot be empty</p><script>window.close()</script>`)
		}
		err := onedrive.Init(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Initialization failure"
			logger.Error("Then init error: ", zap.Error(err), zap.Any("name", "onedrive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Initialization failure:`+err.Error()+`</p><script>window.close()</script>`)
		}
		username, driveId, driveType, err := onedrive.GetInfo(context.Background())
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get user information"
			logger.Error("Then get user information: ", zap.Error(err), zap.Any("name", "onedrive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get user information:`+err.Error()+`</p><script>window.close()</script>`)
		}
		dmap := make(map[string]string)
		dmap["username"] = username

		configs, err := service.MyService.Storage().GetConfig()
		if err != nil {
			notify["status"] = "fail"
			notify["message"] = "Failed to get rclone config"
			logger.Error("Then get config error: ", zap.Error(err), zap.Any("name", "onedrive"))
			service.MyService.Notify().SendNotify(event, notify)
			return ctx.HTML(http.StatusOK, `<p>Failed to get rclone config:`+err.Error()+`</p><script>window.close()</script>`)
		}
		for _, v := range configs.Remotes {
			cf, err := service.MyService.Storage().GetConfigByName(v)
			if err != nil {
				logger.Error("then get config by name error: ", zap.Error(err), zap.Any("name", v))
				continue
			}
			if cf["type"] == "onedrive" && cf["username"] == dmap["username"] {
				if err := service.MyService.Storage().CheckAndMountByName(v); err != nil {
					logger.Error("check and mount by name error: ", zap.Error(err), zap.Any("name", cf["username"]))
				}

				notify["status"] = "warn"
				notify["message"] = "The same configuration has been added"
				service.MyService.Notify().SendNotify(event, notify)
				return ctx.HTML(http.StatusOK, `<p>The same configuration has been added</p><script>window.close()</script>`)
			}
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}
		username += "_onedrive_" + strconv.FormatInt(time.Now().Unix(), 10)

		dmap["client_id"] = onedrive.ClientID
		dmap["client_secret"] = onedrive.ClientSecret
		dmap["token"] = `{"access_token":"` + onedrive.AccessToken + `","token_type":"bearer","refresh_token":"` + onedrive.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*3).Add(time.Minute*50).Format("15:04:05") + `.780385354Z"}`
		dmap["mount_point"] = "/mnt/" + username
		dmap["drive_id"] = driveId
		dmap["drive_type"] = driveType
		service.MyService.Storage().CreateConfig(dmap, username, "onedrive")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")

		notify["status"] = "success"
		notify["message"] = "Success"
		notify["driver"] = "Onedrive"
		service.MyService.Notify().SendNotify(event, notify)
	}

	return ctx.HTML(200, `<p>Just close the page</p><script>window.close()</script>`)
}
