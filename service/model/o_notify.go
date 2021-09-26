package model

type AppNotify struct {
	CustomId    string `gorm:"column:custom_id;primary_key" json:"custom_id"`
	ContainerId string `json:"container_id,omitempty"`
	State       int    `json:"state"` //0:一直在变动的未读消息 1:未读  2:已读
	Message     string `json:"message"`
	CreatedAt   string `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt   string `gorm:"<-:create;<-:update;autoUpdateTime" json:"updated_at"`
	Speed       int    `json:"speed"`
	Id          string `gorm:"-" json:"id"`
	Type        int    `json:"type"` // 1:显示即为已读 2:必须手动点掉 3:error
}

func (p *AppNotify) TableName() string {
	return "o_notify"
}
