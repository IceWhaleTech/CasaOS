package model

import "time"

type PersionModel struct {
	Token     string    `json:"token"`
	Ips       []string  `json:"ips"`
	CreatedAt time.Time `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

//记录链接状态
type ConnectState struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Type      string    `json:"type"` //current state 1:ready 2:ok
	CreatedAt time.Time `json:"created_at"`
	UUId      string    `json:"uuid"` //对接标识
}
