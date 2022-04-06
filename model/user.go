package model

type UserInfo struct {
	NickName string `json:"nick_name"`
	Desc     string `json:"desc"`
	ShareId  string `json:"share_id"`
	Avatar   string `json:"avatar"`
	Version  int    `json:"version,omitempty"`
}
