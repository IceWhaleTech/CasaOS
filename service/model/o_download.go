package model

type PersionDownloadDBModel struct {
	UUID      string `gorm:"column:uuid;primary_key" json:"uuid"`
	State     int    `json:"state"` //
	Type      int    `json:"type"`  //defult 1
	Name      string `json:"name"`  //file name
	Size      int64  `json:"size"`  //file size
	BlockSize int    `json:"block_size"`
	Length    int    `json:"length"` //slice length
	Hash      string `json:"hash"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (p *PersionDownloadDBModel) TableName() string {
	return "o_persion_download"
}
