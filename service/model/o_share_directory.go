package model

import "time"

type ShareDirDBModel struct {
	Id         uint      `gorm:"column:id;primary_key" json:"id"`
	Name       string    `gorm:"size:50" json:"name"`
	Comment    string    `gorm:"size:200" json:"comment"`
	Path       string    `json:"path"`
	ReadOnly   bool      `json:"read_only"`
	Writeable  bool      `json:"writeable"`
	Browseable bool      `json:"browseable"`
	ValidUsers string    `gorm:"size:200" json:"valid_users"` //可以访问的用户  多用户用 , 分割
	CreatedAt  time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt  time.Time `gorm:"<-:create;<-:update" json:"updated_at"`
}

func (p *ShareDirDBModel) TableName() string {
	return "o_share_directory"
}
