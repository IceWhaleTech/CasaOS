package service

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
)

type UserService interface {
	SetUser(username, pwd, token, email, desc string) error
	UpLoadFile(file multipart.File, name string) error
}

type user struct {
}

//设置用户名密码
func (c *user) SetUser(username, pwd, token, email, desc string) error {
	if len(username) > 0 {
		config.Cfg.Section("user").Key("UserName").SetValue(username)
		config.UserInfo.UserName = username
		config.Cfg.Section("user").Key("Initialized").SetValue("true")
		config.UserInfo.Initialized = true
	}
	if len(pwd) > 0 {
		config.Cfg.Section("user").Key("PWD").SetValue(pwd)
		config.UserInfo.PWD = pwd
	}
	if len(token) > 0 {
		config.Cfg.Section("user").Key("Token").SetValue(token)
		config.UserInfo.Token = token
	}
	if len(email) > 0 {
		config.Cfg.Section("user").Key("Email").SetValue(email)
		config.UserInfo.Email = email
	}
	if len(desc) > 0 {
		config.Cfg.Section("user").Key("Description").SetValue(desc)
		config.UserInfo.Description = desc
	}
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	return nil
}

//上传文件
func (c *user) UpLoadFile(file multipart.File, url string) error {
	out, _ := os.OpenFile(url, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer out.Close()
	io.Copy(out, file)
	return nil
}

//获取用户Service
func NewUserService() UserService {
	return &user{}
}
