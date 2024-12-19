package onedrive

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"

	"github.com/IceWhaleTech/CasaOS/internal/driver"
	"github.com/IceWhaleTech/CasaOS/model"
	"go.uber.org/zap"
)

type Onedrive struct {
	model.StorageA
	Addition
	AccessToken string
}

func (d *Onedrive) Config() driver.Config {
	return config
}

func (d *Onedrive) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Onedrive) Init(ctx context.Context) error {
	if d.ChunkSize < 1 {
		d.ChunkSize = 5
	}
	if len(d.RefreshToken) == 0 {
		return d.getRefreshToken()
	}
	return d.refreshToken()
}

func (d *Onedrive) GetUserInfo(ctx context.Context) (string, error) {
	return "", nil
}

func (d *Onedrive) GetInfo(ctx context.Context) (string, string, string, error) {
	url := d.GetMetaUrl(false, "/")
	user := Info{}
	_, err := d.Request(url, http.MethodGet, nil, &user)
	if err != nil {
		return "", "", "", err
	}

	return user.CreatedBy.User.Email, user.ParentReference.DriveID, user.ParentReference.DriveType, nil
}

func (d *Onedrive) GetSpaceSize(ctx context.Context) (used string, total string, err error) {
	host := onedriveHostMap[d.Region]
	url := fmt.Sprintf("%s/v1.0/me/drive/quota", host.Api)
	size := About{}
	resp, err := d.Request(url, http.MethodGet, nil, &size)
	if err != nil {
		return used, total, err
	}
	logger.Info("resp", zap.Any("resp", resp))
	used = strconv.Itoa(size.Used)
	total = strconv.Itoa(size.Total)
	return
}

func (d *Onedrive) Drop(ctx context.Context) error {
	return nil
}

var _ driver.Driver = (*Onedrive)(nil)
