package service

import (
	"gorm.io/gorm"
	ip_helper2 "oasis/pkg/utils/ip_helper"
	loger2 "oasis/pkg/utils/loger"
	"oasis/service/ddns"
	"oasis/service/model"
	"os/exec"
)

type ddnsStruct struct {
	db  *gorm.DB
	log loger2.OLog
}

type DDNSService interface {
	IsExis(t int, domain string, host string) bool
	GetExternalIP() (string, string)
	GetConfigList() *[]model.DDNSList
	DeleteConfig(id uint) bool
	GetType(name string) (uint, string)
	SaveConfig(model model.DDNSUpdataDBModel) error
}

//判断当前添加的是否存在
func (d *ddnsStruct) IsExis(t int, domain string, host string) bool {
	var count int64
	d.db.Table(model.DDNSLISTTABLENAME).Where("type=? AND domain=? AND host=?", t, domain, host).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

//前台获取已配置的ddns列表
func (d *ddnsStruct) GetConfigList() *[]model.DDNSList {
	var s []model.DDNSList
	d.db.Table(model.DDNSLISTTABLENAME).Select("o_ddns_type.name as name,o_ddns.id,host,domain,created_at,updated_at,message,state").Joins("left join o_ddns_type on o_ddns.type=o_ddns_type.id").Scan(&s)
	return &s
}

func (d *ddnsStruct) DeleteConfig(id uint) bool {
	d.db.Delete(&model.DDNSUpdataDBModel{Id: id})
	return true
}

func (d *ddnsStruct) GetExternalIP() (string, string) {
	ipv4 := make(chan string)
	ipv6 := make(chan string)
	go func() { ipv4 <- ip_helper2.GetExternalIPV4() }()
	go func() { ipv6 <- ip_helper2.GetExternalIPV6() }()

	return <-ipv4, <-ipv6
}

func NewDDNSService(db *gorm.DB, log loger2.OLog) DDNSService {
	return &ddnsStruct{db, log}
}

//根据名称获取类型
func (d *ddnsStruct) GetType(name string) (uint, string) {
	var result model.DDNSTypeDBModel
	d.db.Model(&model.DDNSTypeDBModel{}).Where("name = ?", name).First(&result)
	return result.Id, result.Name
}

//保存配置到数据库
func (d *ddnsStruct) GetDockerRootDir(model model.DDNSUpdataDBModel) error {
	return d.db.Create(&model).Error
}

//保存配置到数据库
func (d *ddnsStruct) SaveConfig(model model.DDNSUpdataDBModel) error {
	return d.db.Create(&model).Error
}

//更新数据库ping状态
func chackPing(b chan bool, url string) {
	cmd := exec.Command("ping", url, "-c", "1", "-W", "5")
	err := cmd.Run()
	if err != nil {
		b <- false
	} else {
		b <- true
	}
}

//更新列表
func UpdataDDNSList(db *gorm.DB) {
	var s []model.DDNSCoreList
	db.Table(model.DDNSLISTTABLENAME).Select("o_ddns_type.name as name,o_ddns_type.api_host as api_host,o_ddns.id,`host`,domain,user_name,`password`,`key`,secret,type").Joins("left join o_ddns_type on o_ddns.type=o_ddns_type.id").Scan(&s)
	for _, item := range s {
		var msg string
		switch item.Type {
		case 1:
			var godaddy = &ddns.GoDaddy{
				Host:    item.Host,
				Key:     item.Key,
				Secret:  item.Secret,
				Domain:  item.Domain,
				IPV4:    ip_helper2.GetExternalIPV4(),
				IPV6:    ip_helper2.GetExternalIPV6(),
				ApiHost: item.ApiHost,
			}
			msg = godaddy.Update()
		}

		b := make(chan bool)

		//获取ping状态
		go chackPing(b, item.Host+"."+item.Domain)

		item.State = <-b
		item.Message = msg
		db.Table(model.DDNSLISTTABLENAME).Model(&item).Select("state", "message").Updates(&item)

	}
}
