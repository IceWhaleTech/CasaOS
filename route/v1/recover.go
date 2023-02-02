package v1

import (
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/drivers/dropbox"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/internal/op"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

func GetRecoverStorage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	t := c.Param("type")
	currentTime := time.Now().UTC()
	currentDate := time.Now().UTC().Format("2006-01-02")
	//	timeStr := time.Now().Format("20060102150405")
	if t == "GoogleDrive" {

		gd := op.GetDriverInfoMap()[t]

		add := google_drive.Addition{}
		add.Code = c.Query("code")
		if len(add.Code) == 0 {
			c.String(200, `<p>code cannot be empty</p>`)
			return
		}
		add.RootFolderID = "root"
		for _, v := range gd {
			if v.Name == "client_id" {
				add.ClientID = v.Default
			}
			if v.Name == "client_secret" {
				add.ClientSecret = v.Default
			}
			if v.Name == "chunk_size" {
				cs, err := strconv.ParseInt(v.Default, 10, 64)
				if err != nil {
					cs = 5
				}
				add.ChunkSize = cs
			}
		}

		var google_drive google_drive.GoogleDrive
		google_drive.Addition = add
		err := google_drive.Init(c)
		if err != nil {
			c.String(200, `<p>Initialization failure:`+err.Error()+`</p>`)
			return
		}

		username, err := google_drive.GetUserInfo(c)
		if err != nil {
			c.String(200, `<p>Failed to get user information:`+err.Error()+`</p>`)
			return
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}
		username += "_drive"
		dataMap, _ := service.MyService.Storage().GetConfigByName(username)
		if len(dataMap) > 0 {
			c.String(200, `<p>The same configuration has been added</p>`)
			service.MyService.Storage().CheckAndMountByName(username)
			return
		}
		dmap := make(map[string]string)
		dmap["client_id"] = add.ClientID
		dmap["client_secret"] = add.ClientSecret
		dmap["scope"] = "drive"
		dmap["mount_point"] = "/mnt/" + username
		dmap["token"] = `{"access_token":"` + google_drive.AccessToken + `","token_type":"Bearer","refresh_token":"` + google_drive.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*1).Format("15:04:05") + `Z"}`
		// data.SetValue(username, "type", "drive")
		// data.SetValue(username, "client_id", "865173455964-4ce3gdl73ak5s15kn1vkn73htc8tant2.apps.googleusercontent.com")
		// data.SetValue(username, "client_secret", "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")
		// data.SetValue(username, "scope", "drive")
		// data.SetValue(username, "mount_point", "/mnt/"+username)
		// data.SetValue(username, "token", `{"access_token":"`+google_drive.AccessToken+`","token_type":"Bearer","refresh_token":"`+google_drive.RefreshToken+`","expiry":"`+currentDate+`T`+currentTime.Add(time.Hour*1).Format("15:04:05")+`Z"}`)
		// e = data.Save()
		// if e != nil {
		// 	c.String(200, `<p>保存配置失败:`+e.Error()+`</p>`)
		// 	return
		// }
		service.MyService.Storage().CreateConfig(dmap, username, "drive")
		service.MyService.Storage().MountStorage("/mnt/"+username, username+":")
		notify := make(map[string]interface{})
		notify["status"] = "success"
		service.MyService.Notify().SendNotify("recover_status", notify)
	} else if t == "Dropbox" {

		//mountPath += timeStr

		db := op.GetDriverInfoMap()[t]

		add := dropbox.Addition{}
		add.Code = c.Query("code")
		if len(add.Code) == 0 {
			c.String(200, `<p>code cannot be empty</p>`)
			return
		}
		add.RootFolderID = ""
		for _, v := range db {
			if v.Name == "app_key" {
				add.AppKey = v.Default
			}
			if v.Name == "app_secret" {
				add.AppSecret = v.Default
			}
		}
		var dropbox dropbox.Dropbox
		dropbox.Addition = add
		err := dropbox.Init(c)
		if err != nil {
			c.String(200, `<p>Initialization failure:`+err.Error()+`</p>`)
			return
		}
		username, err := dropbox.GetUserInfo(c)
		if err != nil {
			c.String(200, `<p>Failed to get user information:`+err.Error()+`</p>`)
			return
		}
		if len(username) > 0 {
			a := strings.Split(username, "@")
			username = a[0]
		}
		username += "_dropbox"
		dataMap, _ := service.MyService.Storage().GetConfigByName(username)
		if len(dataMap) > 0 {
			c.String(200, `<p>The same configuration has been added</p>`)
			service.MyService.Storage().CheckAndMountByName(username)
			return
		}
		dmap := make(map[string]string)
		dmap["client_id"] = add.AppKey
		dmap["client_secret"] = add.AppSecret
		dmap["token"] = `{"access_token":"` + dropbox.AccessToken + `","token_type":"bearer","refresh_token":"` + dropbox.Addition.RefreshToken + `","expiry":"` + currentDate + `T` + currentTime.Add(time.Hour*3).Format("15:04:05") + `.780385354Z"}`
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
		notify := make(map[string]interface{})
		notify["status"] = "success"
		service.MyService.Notify().SendNotify("recover_status", notify)
	}

	c.String(200, `<p>Just close the page</p><script>window.close()</script>`)
}
