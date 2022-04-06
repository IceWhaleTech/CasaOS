package model

type PersonDownloadDBModel struct {
	UUID      string `gorm:"column:uuid;primary_key" json:"uuid"`
	State     int    `json:"state"`             //
	Type      int    `json:"type"`              //defult 1
	Name      string `json:"name"`              //file name
	Size      int64  `json:"size"`              //file size
	BlockSize int    `json:"block_size"`        //Size of each file block
	Length    int    `json:"length"`            //slice length
	Hash      string `json:"hash"`              //File hash value
	Error     string `json:"error"`             //
	From      string `json:"from"`              //Error message
	Path      string `json:"path"`              //Full path to the file
	Already   int    `json:"already" gorm:"-"`  //Folder blocks that have been downloaded
	LocalPath string `json:"local_path"`        //The address where the file is saved after download
	Duration  int64  `json:"duration" gorm:"-"` //Length of time
	Created   int64  `gorm:"autoCreateTime" json:"created"`
	Updated   int64  `gorm:"autoCreateTime;autoUpdateTime" json:"updated"`
}

func (p *PersonDownloadDBModel) TableName() string {
	return "o_person_download"
}
