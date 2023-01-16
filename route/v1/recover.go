package v1

import (
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/internal/op"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func GetRecoverStorage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	t := c.Param("type")
	if t == "GoogleDrive" {

		mountPath := "google"

		mountPath += time.Now().Format("20060102150405")

		gd := op.GetDriverInfoMap()[t]
		var req model.Storage
		req.Driver = t
		req.MountPath = mountPath

		req.CacheExpiration = 5
		add := google_drive.Addition{}
		add.Code = c.Query("code")
		if len(add.Code) == 0 {
			c.String(200, `<p>code不可为空</p>`)
			return
		}
		add.RootFolderID = "root"
		for _, v := range gd.Additional {
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

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		addStr, err := json.Marshal(add)
		if err != nil {
			c.String(200, `<p>addition序列化失败</p>`)
			return
		}
		req.Addition = string(addStr)
		logger.Info("GetRecoverStorage", zap.Any("req", req))
		if _, err := service.MyService.Storages().CreateStorage(c, req); err != nil {
			c.String(200, `<p>添加失败:`+err.Error()+`</p>`)
			return
		}
		data := make(map[string]interface{})
		data["status"] = "success"
		service.MyService.Notify().SendNotify("recover_status", data)
	}

	c.String(200, `<p>关闭该页面即可</p><script>window.close()</script>`)
}
