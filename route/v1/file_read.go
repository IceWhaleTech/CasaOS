package v1

import (
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/service"

	"github.com/gin-gonic/gin"
)

type ListReq struct {
	model.PageReq
	Path    string `json:"path" form:"path"`
	Refresh bool   `json:"refresh"`
}
type ObjResp struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"is_dir"`
	Modified time.Time `json:"modified"`
	Sign     string    `json:"sign"`
	Thumb    string    `json:"thumb"`
	Type     int       `json:"type"`
	Path     string    `json:"path"`
}
type FsListResp struct {
	Content  []ObjResp `json:"content"`
	Total    int64     `json:"total"`
	Readme   string    `json:"readme"`
	Write    bool      `json:"write"`
	Provider string    `json:"provider"`
}

func FsList(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
		return
	}
	req.Validate()
	objs, err := service.MyService.FsService().FList(c, req.Path, req.Refresh)
	if err != nil {
		c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
		return
	}
	total, objs := pagination(objs, &req.PageReq)
	provider := "unknown"
	storage, err := service.MyService.FsService().GetStorage(req.Path)
	if err == nil {
		provider = storage.GetStorage().Driver
	}
	c.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: FsListResp{
		Content:  toObjsResp(objs, req.Path, false),
		Total:    int64(total),
		Readme:   "",
		Write:    false,
		Provider: provider,
	}})
}
func pagination(objs []model.Obj, req *model.PageReq) (int, []model.Obj) {
	pageIndex, pageSize := req.Page, req.PerPage
	total := len(objs)
	start := (pageIndex - 1) * pageSize
	if start > total {
		return total, []model.Obj{}
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return total, objs[start:end]
}

func toObjsResp(objs []model.Obj, parent string, encrypt bool) []ObjResp {
	var resp []ObjResp
	for _, obj := range objs {
		thumb, _ := model.GetThumb(obj)
		resp = append(resp, ObjResp{
			Name:     obj.GetName(),
			Size:     obj.GetSize(),
			IsDir:    obj.IsDir(),
			Modified: obj.ModTime(),
			Sign:     "",
			Path:     filepath.Join(parent, obj.GetName()),
			Thumb:    thumb,
			Type:     0,
		})
	}
	return resp
}
