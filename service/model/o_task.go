package model

type TaskDBModel struct {
	Id        uint   `gorm:"column:id;primary_key" json:"id"`
	Title     string `gorm:"size:200" json:"title"`
	ImageUrl  string `json:"image_url"`
	Content   string `gorm:"size:255" json:"content"`
	Url       string `json:"url"`
	State     int    `json:"state"` // 0:未阅读,1:已阅读
	Type      int    `json:"type"`
	CreatedAt string `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"<-:create;<-:update;autoUpdateTime" json:"updated_at"`
}

func (p *TaskDBModel) TableName() string {
	return "o_task"
}
