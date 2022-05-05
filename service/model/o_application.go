package model

import (
	"time"
)

type ApplicationModel struct {
	Id        int       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	State     int       `json:"state"`
	Type      string    `json:"type"`
	Index     int       `json:"index"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"<-:create;<-:update" json:"updated_at"`
}

func (p *ApplicationModel) TableName() string {
	return "o_application"
}
