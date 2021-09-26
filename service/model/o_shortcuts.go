package model

import "time"

type ShortcutsDBModel struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"<-:create;<-:update" json:"updated_at"`
}

func (p *ShortcutsDBModel) TableName() string {
	return "o_shortcuts"
}
