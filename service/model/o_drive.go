package model

type PeerDriveDBModel struct {
	ID          string `gorm:"column:id;primary_key" json:"id"`
	Updated     int64  `gorm:"autoUpdateTime"`
	Created     int64  `gorm:"autoCreateTime"`
	UserAgent   string `json:"user_agent"`
	DisplayName string `json:"display_name"`
	DeviceName  string `json:"device_name"`
	IP          string `json:"ip"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Online      bool   `gorm:"-" json:"online"`
}
