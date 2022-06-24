/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-22 18:50:44
 * @FilePath: /CasaOS/pkg/sqlite/db.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package sqlite

import (
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/loger"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gdb *gorm.DB

func GetDb(dbPath string) *gorm.DB {
	if gdb != nil {
		return gdb
	}
	// Refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
	//dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.PWD, m.IP, m.Port, m.DBName)
	//db, err := gorm.Open(mysql2.Open(dsn), &gorm.Config{})
	file.IsNotExistMkDir(dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath+"/casaOS.db"), &gorm.Config{})
	c, _ := db.DB()
	c.SetMaxIdleConns(10)
	c.SetMaxOpenConns(100)
	c.SetConnMaxIdleTime(time.Second * 1000)
	if err != nil {
		loger.Error("sqlite connect error", zap.Any("db connect error", err))
		panic("sqlite connect error")
		return nil
	}
	gdb = db
	err = db.AutoMigrate(&model2.AppNotify{}, &model2.AppListDBModel{}, &model2.SerialDisk{}, model2.PersonDownloadDBModel{}, model2.FriendModel{}, model2.PersonDownRecordDBModel{}, model2.UserDBModel{})
	db.Exec("DROP TABLE IF EXISTS o_application")
	if err != nil {
		loger.Error("check or create db error", zap.Any("error", err))
	}
	return db
}
