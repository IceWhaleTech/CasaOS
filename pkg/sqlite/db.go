/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-31 13:39:24
 * @FilePath: /CasaOS/pkg/sqlite/db.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package sqlite

import (
	"errors"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var gdb *gorm.DB

var ErrDatabaseConnection = errors.New("failed to connect to sqlite database")

func GetDb(dbPath string) (*gorm.DB, error) {
	if gdb != nil {
		return gdb, nil
	}

	file.IsNotExistMkDir(dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath+"/casaOS.db"), &gorm.Config{})
	if err != nil {
		logger.Error("sqlite connect error", zap.Error(err), zap.String("path", dbPath))
		return nil, ErrDatabaseConnection
	}

	c, err := db.DB()
	if err != nil {
		logger.Error("failed to get underlying sql.DB", zap.Error(err))
		return nil, err
	}
	c.SetMaxIdleConns(10)
	c.SetMaxOpenConns(1)
	c.SetConnMaxIdleTime(time.Second * 1000)
	gdb = db

	err = db.AutoMigrate(&model2.AppNotify{}, model2.SharesDBModel{}, model2.ConnectionsDBModel{}, model2.PeerDriveDBModel{})
	if err != nil {
		logger.Error("failed to auto-migrate database models", zap.Error(err))
	}

	db.Exec("DROP TABLE IF EXISTS o_application")
	db.Exec("DROP TABLE IF EXISTS o_friend")
	db.Exec("DROP TABLE IF EXISTS o_person_download")
	db.Exec("DROP TABLE IF EXISTS o_person_down_record")
	return db, nil
}
