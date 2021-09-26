package model

import "time"

type Version struct {
	Id        uint      `gorm:"column:id;primary_key" json:"id"`
	ChangLog  string    `json:"chang_log"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
