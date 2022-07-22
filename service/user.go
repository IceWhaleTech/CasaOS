/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-03-18 11:40:55
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-07-12 10:05:37
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

	"github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/gorm"
)

type UserService interface {
	UpLoadFile(file multipart.File, name string) error
	CreateUser(m model.UserDBModel) model.UserDBModel
	GetUserCount() (userCount int64)
	UpdateUser(m model.UserDBModel)
	UpdateUserPassword(m model.UserDBModel)
	GetUserInfoById(id string) (m model.UserDBModel)
	GetUserAllInfoById(id string) (m model.UserDBModel)
	GetUserAllInfoByName(userName string) (m model.UserDBModel)
	DeleteUserById(id string)
	DeleteAllUser()
	GetUserInfoByUserName(userName string) (m model.UserDBModel)
	GetAllUserName() (list []model.UserDBModel)
}

var UserRegisterHash = make(map[string]string)

type userService struct {
	db *gorm.DB
}

func (u *userService) DeleteAllUser() {
	u.db.Where("1=1").Delete(&model.UserDBModel{})
}
func (u *userService) DeleteUserById(id string) {
	u.db.Where("id= ?", id).Delete(&model.UserDBModel{})
}

func (u *userService) GetAllUserName() (list []model.UserDBModel) {
	u.db.Select("username").Find(&list)
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
	u.db.Model(&m).Omit("password").Updates(&m)
}
func (u *userService) UpdateUserPassword(m model.UserDBModel) {
	u.db.Model(&m).Update("password", m.Password)
}
func (u *userService) GetUserAllInfoById(id string) (m model.UserDBModel) {
	u.db.Where("id= ?", id).First(&m)
	return
}
func (u *userService) GetUserAllInfoByName(userName string) (m model.UserDBModel) {
	u.db.Where("username= ?", userName).First(&m)
	return
}
func (u *userService) GetUserInfoById(id string) (m model.UserDBModel) {
	u.db.Select("username", "id", "role", "nickname", "description", "avatar", "email").Where("id= ?", id).First(&m)
	return
}

func (u *userService) GetUserInfoByUserName(userName string) (m model.UserDBModel) {
	u.db.Select("username", "id", "role", "nickname", "description", "avatar", "email").Where("username= ?", userName).First(&m)
	return
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
