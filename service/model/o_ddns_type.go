package model

type DDNSTypeDBModel struct {
	Id      uint   `gorm:"column:id;primary_key" json:"id"`
	Name    string `json:"name"`
	ApiHost string `json:"api_host"`
}

func (p *DDNSTypeDBModel) TableName() string {
	return "o_ddns_type"
}
