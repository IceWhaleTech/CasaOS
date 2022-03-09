package sqlite

import (
	"fmt"
	"time"

	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gdb *gorm.DB

func GetDb(projectPath string) *gorm.DB {
	if gdb != nil {
		return gdb
	}
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.PWD, m.IP, m.Port, m.DBName)
	//db, err := gorm.Open(mysql2.Open(dsn), &gorm.Config{})
	file.IsNotExistMkDir(projectPath + "/db/")
	db, err := gorm.Open(sqlite.Open(projectPath+"/db/casaOS.db"), &gorm.Config{})
	c, _ := db.DB()
	c.SetMaxIdleConns(10)
	c.SetMaxOpenConns(100)
	c.SetConnMaxIdleTime(time.Second * 1000)
	if err != nil {
		fmt.Println("连接数据失败!")
		panic("数据库连接失败")
		return nil
	}
	gdb = db
	err = db.AutoMigrate(&model2.TaskDBModel{}, &model2.AppNotify{}, &model2.AppListDBModel{}, &model2.SerialDisk{}, model2.PersionDownloadDBModel{})
	if err != nil {
		fmt.Println("检查和创建数据库出错", err)
	}
	return db
}
