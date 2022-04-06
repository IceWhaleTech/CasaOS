package model

type FriendModel struct {
	State     int    `json:"state"` //Reserved
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"`
	NickName  string `json:"nick_name"`
	Mark      string `json:"mark"`   //Remarks
	Block     bool   `json:"block"`  //Disable or not
	Avatar    string `json:"avatar"` //User avatar
	Token     string `gorm:"column:token;primary_key" json:"token"`
	Profile   string `json:"profile"` //Description
	OnLine    bool   `json:"on_line" gorm:"-"`
	Version   int    `json:"version"`
}

func (p *FriendModel) TableName() string {
	return "o_friend"
}
