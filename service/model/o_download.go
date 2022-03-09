package model

type PersionDownloadDBModel struct {
	UUID      string `gorm:"column:uuid;primary_key" json:"uuid"`
	State     int    `json:"state"`     //
	Type      int    `json:"type"`      //defult 1
	Name      string `json:"name"`      //file name
	TempPath  string `json:"temp_path"` //temp path
	Size      int64  `json:"size"`      //file size
	Section   string `json:"section"`
	Length    int    `json:"length"` //slice length
	Hash      string `json:"hash"`
	CreatedAt string `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"<-:create;<-:update;autoUpdateTime" json:"updated_at"`
}

func (p *PersionDownloadDBModel) TableName() string {
	return "o_persion_download"
}

type PersionFileSectionModel struct {
	Index int    `json:"index"`
	Hash  string `json:"hash"`
}
