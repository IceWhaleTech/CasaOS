package model

import "time"

type PersionModel struct {
	Token     string    `json:"token"`
	Ips       []string  `json:"ips"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//记录链接状态
type ConnectState struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Type      string    `json:"type"` //current state 1:ready 2:ok
	CreatedAt time.Time `json:"created_at"`
	UUId      string    `json:"uuid"` //对接标识
}
type MessageModel struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	UUId string      `json:"uuid"`
	From string      `json:"from"`
	To   string      `json:"to"`
}

type TranFileModel struct {
	Hash   string `json:"hash"` //Verify current fragment integrity
	Length int    `json:"length"`
	Index  int    `json:"index"`
}

//需要获取文件详情
type FileDetailModel struct {
	Path string `json:"path"`
}

//返回文件详情
type FileSummaryModel struct {
	Hash      string `json:"hash"` //Verify file
	Name      string `json:"name"`
	BlockSize int    `json:"block_size"`
	Length    int    `json:"length"`
	Size      int64  `json:"size"`
}
