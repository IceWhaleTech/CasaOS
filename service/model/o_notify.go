package model

type AppNotify struct {
	State     int    `json:"state"` //0:一直在变动的未读消息 1:未读  2:已读
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Id        string `json:"id"`
	Type      int    `json:"type"` // 1:显示即为已读 2:info 3:warning 4:error 5:success
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	CustomId  string `gorm:"column:custom_id;primary_key" json:"custom_id"`
}

func (p *AppNotify) TableName() string {
	return "o_notify"
}
