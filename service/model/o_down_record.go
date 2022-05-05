package model

type PersonDownRecordDBModel struct {
	UUID       string `gorm:"column:uuid;primary_key" json:"uuid"`
	Name       string `json:"name"` //file name
	Type       int    `json:"type"`
	Size       int64  `json:"size"`       //file size
	Downloader string `json:"downloader"` //Error message
	Path       string `json:"path"`
	Created    int64  `gorm:"autoCreateTime" json:"created"`
	Updated    int64  `gorm:"autoCreateTime;autoUpdateTime" json:"updated"`
}

func (p *PersonDownRecordDBModel) TableName() string {
	return "o_person_down_record"
}
