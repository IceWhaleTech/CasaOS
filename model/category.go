package model

type ServerCategoryList struct {
	Id uint `gorm:"column:id;primary_key" json:"id"`
	//CreatedAt time.Time `json:"created_at"`
	//
	//UpdatedAt time.Time `json:"updated_at"`
	Name  string `json:"name"`
	Count uint   `json:"count"`
}
