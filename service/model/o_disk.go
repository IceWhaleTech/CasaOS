package model

//SerialAdvanced Technology Attachment (STAT)
type SerialDisk struct {
	Id         uint   `gorm:"column:id;primary_key" json:"id"`
	Serial     string `json:"serial"`
	Path       string `json:"path"`
	State      int    `json:"state"`
	MountPoint string `json:"mount_point"`
}

func (p *SerialDisk) TableName() string {
	return "o_disk"
}
