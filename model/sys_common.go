/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-18 14:32:16
 * @FilePath: /CasaOS/model/sys_common.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

import "time"

// 系统配置
type SysInfoModel struct {
	Name string // 系统名称
}

// 用户相关
type UserModel struct {
	UserName    string
	PWD         string
	Token       string
	Head        string
	Email       string
	Description string
	Initialized bool
	Avatar      string
	NickName    string
	Public      string
}

// 服务配置
type ServerModel struct {
	HttpPort    string
	RunMode     string
	ServerApi   string
	LockAccount bool
	Token       string
	SocketPort  string
}

// 服务配置
type APPModel struct {
	LogPath        string
	LogSaveName    string
	LogFileExt     string
	DateStrFormat  string
	DateTimeFormat string
	UserDataPath   string
	TimeFormat     string
	DateFormat     string
	DBPath         string
	ShellPath      string
	TempPath       string
}
type CommonModel struct {
	RuntimePath string
}

// 公共返回模型
type Result struct {
	Success int         `json:"success" example:"200"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data" example:"返回结果"`
}

// redis配置文件
type RedisModel struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type SystemConfig struct {
	ConfigPath string `json:"config_path"`
}

type CasaOSGlobalVariables struct {
	AppChange bool
}

type FileSetting struct {
	ShareDir    []string `json:"share_dir" delim:"|"`
	DownloadDir string   `json:"download_dir"`
}
