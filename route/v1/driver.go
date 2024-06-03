package v1

import (
	"github.com/IceWhaleTech/CasaOS-Common/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/drivers/dropbox"
	"github.com/IceWhaleTech/CasaOS/drivers/google_drive"
	"github.com/IceWhaleTech/CasaOS/drivers/onedrive"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/labstack/echo/v4"
)

func ListDriverInfo(ctx echo.Context) error {
	list := []model.Drive{}

	google := google_drive.GetConfig()
	list = append(list, model.Drive{
		Name:    "Google Drive",
		Icon:    google.Icon,
		AuthUrl: google.AuthUrl,
	})
	dp := dropbox.GetConfig()
	list = append(list, model.Drive{
		Name:    "Dropbox",
		Icon:    dp.Icon,
		AuthUrl: dp.AuthUrl,
	})
	od := onedrive.GetConfig()
	list = append(list, model.Drive{
		Name:    "OneDrive",
		Icon:    od.Icon,
		AuthUrl: od.AuthUrl,
	})
	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}
