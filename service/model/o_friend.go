package model

type FriendModel struct {
	State     int    `json:"state"` //备用
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"`
	NickName  string `json:"nick_name"` //custom name
	Avatar    string `json:"avatar"`    //头像
	Name      string `json:"name"`
	Token     string `gorm:"column:token;primary_key" json:"token"`
	Profile   string `json:"profile"`
}

func (p *FriendModel) TableName() string {
	return "o_friend"
}
