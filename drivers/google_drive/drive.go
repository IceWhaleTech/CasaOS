package google_drive

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/base"
	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type GoogleDrive struct {
	model.StorageA
	Addition
	AccessToken string
}

func (d *GoogleDrive) Config() driver.Config {
	return config
}

func (d *GoogleDrive) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *GoogleDrive) Init(ctx context.Context) error {
	if d.ChunkSize == 0 {
		d.ChunkSize = 5
	}
	if len(d.RefreshToken) == 0 {
		d.getRefreshToken()
	}
	return d.refreshToken()
}

func (d *GoogleDrive) Drop(ctx context.Context) error {
	return nil
}

func (d *GoogleDrive) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	files, err := d.getFiles(dir.GetID())
	if err != nil {
		return nil, err
	}
	return utils.SliceConvert(files, func(src File) (model.Obj, error) {
		return fileToObj(src), nil
	})
}

func (d *GoogleDrive) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	url := fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?includeItemsFromAllDrives=true&supportsAllDrives=true", file.GetID())
	_, err := d.request(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	link := model.Link{
		Method: http.MethodGet,
		URL:    url + "&alt=media",
		Header: http.Header{
			"Authorization": []string{"Bearer " + d.AccessToken},
		},
	}
	return &link, nil
}
func (d *GoogleDrive) GetUserInfo(ctx context.Context) (string, error) {
	url := "https://content.googleapis.com/drive/v3/about?fields=user"
	user := UserInfo{}
	resp, err := d.request(url, http.MethodGet, nil, &user)
	if err != nil {
		return "", err
	}
	logger.Info("resp", zap.Any("resp", resp))
	return user.User.EmailAddress, nil
}

func (d *GoogleDrive) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	data := base.Json{
		"name":     dirName,
		"parents":  []string{parentDir.GetID()},
		"mimeType": "application/vnd.google-apps.folder",
	}
	_, err := d.request("https://www.googleapis.com/drive/v3/files", http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, nil)
	return err
}

func (d *GoogleDrive) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	query := map[string]string{
		"addParents":    dstDir.GetID(),
		"removeParents": "root",
	}
	url := "https://www.googleapis.com/drive/v3/files/" + srcObj.GetID()
	_, err := d.request(url, http.MethodPatch, func(req *resty.Request) {
		req.SetQueryParams(query)
	}, nil)
	return err
}

func (d *GoogleDrive) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	data := base.Json{
		"name": newName,
	}
	url := "https://www.googleapis.com/drive/v3/files/" + srcObj.GetID()
	_, err := d.request(url, http.MethodPatch, func(req *resty.Request) {
		req.SetBody(data)
	}, nil)
	return err
}

func (d *GoogleDrive) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return errors.New("not support")
}

func (d *GoogleDrive) Remove(ctx context.Context, obj model.Obj) error {
	url := "https://www.googleapis.com/drive/v3/files/" + obj.GetID()
	_, err := d.request(url, http.MethodDelete, nil, nil)
	return err
}

func (d *GoogleDrive) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	obj := stream.GetOld()
	var (
		e    Error
		url  string
		data base.Json
		res  *resty.Response
		err  error
	)
	if obj != nil {
		url = fmt.Sprintf("https://www.googleapis.com/upload/drive/v3/files/%s?uploadType=resumable&supportsAllDrives=true", obj.GetID())
		data = base.Json{}
	} else {
		data = base.Json{
			"name":    stream.GetName(),
			"parents": []string{dstDir.GetID()},
		}
		url = "https://www.googleapis.com/upload/drive/v3/files?uploadType=resumable&supportsAllDrives=true"
	}
	req := base.NoRedirectClient.R().
		SetHeaders(map[string]string{
			"Authorization":           "Bearer " + d.AccessToken,
			"X-Upload-Content-Type":   stream.GetMimetype(),
			"X-Upload-Content-Length": strconv.FormatInt(stream.GetSize(), 10),
		}).
		SetError(&e).SetBody(data).SetContext(ctx)
	if obj != nil {
		res, err = req.Patch(url)
	} else {
		res, err = req.Post(url)
	}
	if err != nil {
		return err
	}
	if e.Error.Code != 0 {
		if e.Error.Code == 401 {
			err = d.refreshToken()
			if err != nil {
				return err
			}
			return d.Put(ctx, dstDir, stream, up)
		}
		return fmt.Errorf("%s: %v", e.Error.Message, e.Error.Errors)
	}
	putUrl := res.Header().Get("location")
	if stream.GetSize() < d.ChunkSize*1024*1024 {
		_, err = d.request(putUrl, http.MethodPut, func(req *resty.Request) {
			req.SetHeader("Content-Length", strconv.FormatInt(stream.GetSize(), 10)).SetBody(stream.GetReadCloser())
		}, nil)
	} else {
		err = d.chunkUpload(ctx, stream, putUrl)
	}
	return err
}

var _ driver.Driver = (*GoogleDrive)(nil)
