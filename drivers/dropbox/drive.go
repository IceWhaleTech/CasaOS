package dropbox

import (
	"context"
	"errors"
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Dropbox struct {
	model.StorageA
	Addition
	AccessToken string
}

func (d *Dropbox) Config() driver.Config {
	return config
}

func (d *Dropbox) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Dropbox) Init(ctx context.Context) error {
	if len(d.RefreshToken) == 0 {
		d.getRefreshToken()
	}
	return d.refreshToken()
}

func (d *Dropbox) Drop(ctx context.Context) error {

	return nil
}

func (d *Dropbox) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	files, err := d.getFiles(dir.GetID())
	if err != nil {
		return nil, err
	}
	return utils.SliceConvert(files, func(src File) (model.Obj, error) {
		return fileToObj(src), nil
	})
}

func (d *Dropbox) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	url := "https://content.dropboxapi.com/2/files/download"
	link := model.Link{
		URL:    url,
		Method: http.MethodPost,
		Header: http.Header{
			"Authorization":   []string{"Bearer " + d.AccessToken},
			"Dropbox-API-Arg": []string{`{"path": "` + file.GetPath() + `"}`},
		},
	}
	return &link, nil
}
func (d *Dropbox) GetUserInfo(ctx context.Context) (string, error) {
	url := "https://api.dropboxapi.com/2/users/get_current_account"
	user := UserInfo{}
	resp, err := d.request(url, http.MethodPost, func(req *resty.Request) {
		req.SetHeader("Content-Type", "")
	}, &user)
	if err != nil {
		return "", err
	}
	logger.Info("resp", zap.Any("resp", string(resp)))
	return user.Email, nil
}
func (d *Dropbox) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	return nil
}

func (d *Dropbox) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	return nil
}

func (d *Dropbox) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	return nil
}

func (d *Dropbox) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errors.New("not support")
}

func (d *Dropbox) Remove(ctx context.Context, obj model.Obj) error {
	return nil
}

func (d *Dropbox) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	return nil
}
func (d *Dropbox) GetInfo(ctx context.Context) (string, string, string, error) {
	return "", "", "", nil
}

var _ driver.Driver = (*Dropbox)(nil)
