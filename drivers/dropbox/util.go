package dropbox

import (
	"fmt"
	"net/http"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/base"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	app_key    = "private build"
	app_secret = "private build"
)

func (d *Dropbox) getRefreshToken() error {
	url := "https://api.dropbox.com/oauth2/token"
	var resp base.TokenResp
	var e TokenError

	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).
		SetFormData(map[string]string{
			"code":         d.Code,
			"grant_type":   "authorization_code",
			"redirect_uri": "https://cloudoauth.files.casaos.app",
		}).SetBasicAuth(d.Addition.AppKey, d.Addition.AppSecret).SetHeader("Content-Type", "application/x-www-form-urlencoded").Post(url)
	if err != nil {
		return err
	}
	logger.Info("get refresh token", zap.String("res", res.String()))
	if e.Error != "" {
		return fmt.Errorf(e.Error)
	}
	d.RefreshToken = resp.RefreshToken
	return nil

}
func (d *Dropbox) refreshToken() error {
	url := "https://api.dropbox.com/oauth2/token"
	var resp base.TokenResp
	var e TokenError

	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).
		SetFormData(map[string]string{
			"refresh_token": d.RefreshToken,
			"grant_type":    "refresh_token",
		}).SetBasicAuth(d.Addition.AppKey, d.Addition.AppSecret).SetHeader("Content-Type", "application/x-www-form-urlencoded").Post(url)
	if err != nil {
		return err
	}
	logger.Info("get refresh token", zap.String("res", res.String()))
	if e.Error != "" {
		return fmt.Errorf(e.Error)
	}
	d.AccessToken = resp.AccessToken
	return nil

}
func (d *Dropbox) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeader("Authorization", "Bearer "+d.AccessToken)
	req.SetHeader("Content-Type", "application/json")
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	var e Error
	req.SetError(&e)
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	if e.Error.Code != 0 {
		if e.Error.Code == 401 {
			err = d.refreshToken()
			if err != nil {
				return nil, err
			}
			return d.request(url, method, callback, resp)
		}
		return nil, fmt.Errorf("%s: %v", e.Error.Message, e.Error.Errors)
	}
	return res.Body(), nil
}
func (d *Dropbox) getFiles(path string) ([]File, error) {

	res := make([]File, 0)
	var resp Files
	body := base.Json{
		"limit": 2000,
		"path":  path,
	}

	_, err := d.request("https://api.dropboxapi.com/2/files/list_folder", http.MethodPost, func(req *resty.Request) {
		req.SetBody(body)
	}, &resp)
	if err != nil {
		return nil, err
	}
	res = append(res, resp.Files...)

	return res, nil
}
func GetConfig() Dropbox {
	dp := Dropbox{}
	dp.RootFolderID = ""
	dp.AuthUrl = "https://www.dropbox.com/oauth2/authorize?client_id=" + app_key + "&redirect_uri=https://cloudoauth.files.casaos.app&response_type=code&token_access_type=offline&state=${HOST}%2Fv2%2Frecover%2FDropbox&&force_reapprove=true&force_reauthentication=true"
	dp.AppKey = app_key
	dp.AppSecret = app_secret
	dp.Icon = "./img/driver/Dropbox.svg"
	return dp
}
