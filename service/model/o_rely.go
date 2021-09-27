package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/IceWhaleTech/CasaOS/service/docker_base"
	"time"
)

type RelyDBModel struct {
	Id                uint         `gorm:"column:id;primary_key" json:"id"`
	CustomId          string       ` json:"custom_id"`
	ContainerCustomId string       `json:"container_custom_id"`
	Config            MysqlConfigs `json:"config"`
	ContainerId       string       `json:"container_id,omitempty"`
	Type              int          `json:"type"` //目前暂未使用
	CreatedAt         time.Time    `gorm:"<-:create" json:"created_at"`
	UpdatedAt         time.Time    `gorm:"<-:create;<-:update" json:"updated_at"`
}

/****************使gorm支持[]string结构*******************/
type MysqlConfigs docker_base.MysqlConfig

func (c MysqlConfigs) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *MysqlConfigs) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

/****************使gorm支持[]string结构*******************/

func (p RelyDBModel) TableName() string {
	return "o_rely"
}
