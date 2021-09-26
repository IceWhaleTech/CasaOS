package model

import "time"

func (p *DDNSUpdataDBModel) TableName() string {
	return "o_ddns"
}

type DDNSUpdataDBModel struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Ipv4      string    `gorm:"-"`
	Ipv6      string    `gorm:"-"`
	Type      uint      `json:"type" form:"type"`
	Domain    string    `json:"domain" form:"domain"`
	Host      string    `json:"host" form:"host"`
	Key       string    `json:"key" form:"key"`
	Secret    string    `json:"secret" form:"secret"`
	UserName  string    `json:"user_name" form:"user_name"`
	Password  string    `json:"password" form:"password"`
	CreatedAt      time.Time       `gorm:"<-:create" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"<-:create;<-:update" json:"updated_at"`
}

const DDNSLISTTABLENAME = "o_ddns"

//返回给前台使用
type DDNSList struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain" form:"domain"`
	Host      string    `json:"host" form:"host"`
	IPV4      string    `json:"ipv_4" gorm:"-"`
	IPV6      string    `json:"ipv_6" gorm:"-"`
	Message   string    `json:"message"`
	State     bool      `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//定时任务使用
type DDNSCoreList struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Domain   string `json:"domain" form:"domain"`
	Name     string `json:"domain" form:"name"`
	Type     uint   `json:"type"`
	Key      string `json:"key"`
	Message  string `json:"message"`
	State    bool   `json:"state"`
	Secret   string `json:"secret" form:"secret"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
	ApiHost  string `json:"api_host"`
	Host     string `json:"host"`
	IPV4     string `json:"ipv_4" gorm:"-"`
	IPV6     string `json:"ipv_6" gorm:"-"`
}
