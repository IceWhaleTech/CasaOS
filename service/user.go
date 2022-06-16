/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-03-18 11:40:55
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-16 19:08:09
 * @FilePath: /CasaOS/service/user.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package service

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type UserService interface {
	SetUser(username, pwd, token, email, desc, nickName string) error
	UpLoadFile(file multipart.File, name string) error
	CreateUser(m model.UserDBModel) model.UserDBModel
	GetUserCount() (userCount int64)
	UpdateUser(m model.UserDBModel)
	GetUserInfoById(id string) (m model.UserDBModel)
	GetUserInfoByUserName(userName string) (m model.UserDBModel)
	GetAllUserName() (list []model.UserDBModel)
}

var UserRegisterHash = make(map[string]string)

type userService struct {
	db *gorm.DB
}

func (u *userService) GetAllUserName() (list []model.UserDBModel) {
	u.db.Select("user_name").Find(&list)
	return
}
func (u *userService) CreateUser(m model.UserDBModel) model.UserDBModel {
	u.db.Create(&m)
	return m
}

func (u *userService) GetUserCount() (userCount int64) {
	u.db.Find(&model.UserDBModel{}).Count(&userCount)
	return
}

func (u *userService) UpdateUser(m model.UserDBModel) {
	u.db.Save(&m)
}

func (u *userService) GetUserInfoById(id string) (m model.UserDBModel) {
	u.db.Where("id= ?", id).First(&m)
	return
}

func (u *userService) GetUserInfoByUserName(userName string) (m model.UserDBModel) {
	u.db.Where("user_name= ?", userName).First(&m)
	return
}

//设置用户名密码
func (u *userService) SetUser(username, pwd, token, email, desc, nickName string) error {
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
	if len(email) > 0 {
		config.Cfg.Section("user").Key("Email").SetValue(email)
		config.UserInfo.Email = email
	}
	if len(desc) > 0 {
		config.Cfg.Section("user").Key("Description").SetValue(desc)
		config.UserInfo.Description = desc
	}
	if len(nickName) > 0 {
		config.Cfg.Section("user").Key("NickName").SetValue(nickName)
		config.UserInfo.NickName = nickName
	}
	config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	return nil
}

//上传文件
func (c *userService) UpLoadFile(file multipart.File, url string) error {
	out, _ := os.OpenFile(url, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer out.Close()
	io.Copy(out, file)
	return nil
}

//获取用户Service
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}
