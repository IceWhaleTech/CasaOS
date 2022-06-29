package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ServerAppListCollection struct {
	List      []ServerAppList `json:"list"`
	Recommend []ServerAppList `json:"recommend"`
	Community []ServerAppList `json:"community"`
	Version   string          `json:"version"`
}

type ServerAppList struct {
	Id             uint      `gorm:"column:id;primary_key" json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Tagline        string    `json:"tagline"`
	Tags           Strings   `gorm:"type:json" json:"tags"`
	Icon           string    `json:"icon"`
	ScreenshotLink Strings   `gorm:"type:json" json:"screenshot_link"`
	Category       string    `json:"category"`
	CategoryId     int       `json:"category_id"`
	CategoryFont   string    `json:"category_font"`
	PortMap        string    `json:"port_map"`
	ImageVersion   string    `json:"image_version"`
	Tip            string    `json:"tip"`
	Envs           EnvArray  `json:"envs"`
	Ports          PortArray `json:"ports"`
	Volumes        PathArray `json:"volumes"`
	Devices        PathArray `json:"devices"`
	NetworkModel   string    `json:"network_model"`
	Image          string    `json:"image"`
	Index          string    `json:"index"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	State          string    `json:"state"`
	Author         string    `json:"author"`
	MinMemory      int       `json:"min_memory"`
	MinDisk        int       `json:"min_disk"`
	MaxMemory      uint64    `json:"max_memory"`
	Thumbnail      string    `json:"thumbnail"`
	Healthy        string    `json:"healthy"`
	Plugins        Strings   `json:"plugins"`
	Origin         string    `json:"origin"`
	Type           int       `json:"type"`
	QueryCount     int       `json:"query_count"`
	Developer      string    `json:"developer"`
	HostName       string    `json:"host_name"`
	Privileged     bool      `json:"privileged"`
	CapAdd         Strings   `json:"cap_add"`
	Cmd            Strings   `json:"cmd"`
}

type Ports struct {
	ContainerPort uint   `json:"container_port"`
	CommendPort   int    `json:"commend_port"`
	Desc          string `json:"desc"`
	Type          int    `json:"type"` //  1:必选 2:可选 3:默认值不必显示 4:系统处理  5:container内容也可编辑
}

type Volume struct {
	ContainerPath string `json:"container_path"`
	Path          string `json:"path"`
	Desc          string `json:"desc"`
	Type          int    `json:"type"` //  1:必选 2:可选 3:默认值不必显示 4:系统处理   5:container内容也可编辑
}

type Envs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"` //  1:必选 2:可选 3:默认值不必显示 4:系统处理 5:container内容也可编辑
}

type Devices struct {
	ContainerPath string `json:"container_path"`
	Path          string `json:"path"`
	Desc          string `json:"desc"`
	Type          int    `json:"type"` //  1:必选 2:可选 3:默认值不必显示 4:系统处理 5:container内容也可编辑
}

type configures struct {
	TcpPorts []Ports   `json:"tcp_ports"`
	UdpPorts []Ports   `json:"udp_ports"`
	Envs     []Envs    `json:"envs"`
	Volumes  []Volume  `json:"volumes"`
	Devices  []Devices `json:"devices"`
}

/****************使gorm支持[]string结构*******************/
type Strings []string

func (c Strings) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Strings) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

/****************使gorm支持[]string结构*******************/

/****************使gorm支持[]string结构*******************/
type MapStrings []map[string]string

func (c MapStrings) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *MapStrings) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

/****************使gorm支持[]string结构*******************/
