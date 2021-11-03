package model

import (
	"database/sql/driver"
	"encoding/json"
)

type TcpPorts struct {
	Desc          string `json:"desc"`
	ContainerPort int    `json:"container_port"`
}
type UdpPorts struct {
	Desc          string `json:"desc"`
	ContainerPort int    `json:"container_port"`
}

/*******************使用gorm支持json************************************/

type PortMap struct {
	ContainerPort string `json:"container"`
	CommendPort   string `json:"host"`
	Protocol      string `json:"protocol"`
	Desc          string `json:"desc"`
	Type          int    `json:"type"`
}

type PortArray []PortMap

// Value 实现方法
func (p PortArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *PortArray) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}

/************************************************************************/

/*******************使用gorm支持json************************************/

type Env struct {
	Name  string `json:"container"`
	Value string `json:"host"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"`
}

type JSON json.RawMessage

type EnvArray []Env

// Value 实现方法
func (p EnvArray) Value() (driver.Value, error) {
	return json.Marshal(p)
	//return .MarshalJSON()
}

// Scan 实现方法
func (p *EnvArray) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}

/************************************************************************/

/*******************使用gorm支持json************************************/

type PathMap struct {
	ContainerPath string `json:"container"`
	Path          string `json:"host"`
	Type          int    `json:"type"`
	Desc          string `json:"desc"`
}

type PathArray []PathMap

// Value 实现方法
func (p PathArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *PathArray) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}

/************************************************************************/

//type PostData struct {
//	Envs       EnvArrey  `json:"envs,omitempty"`
//	Udp        PortArrey `json:"udp_ports"`
//	Tcp        PortArrey `json:"tcp_ports"`
//	Volumes    PathArrey `json:"volumes"`
//	Devices    PathArrey `json:"devices"`
//	Port       string    `json:"port,omitempty"`
//	PortMap    string    `json:"port_map"`
//	CpuShares  int64     `json:"cpu_shares,omitempty"`
//	Memory     int64     `json:"memory,omitempty"`
//	Restart    string    `json:"restart,omitempty"`
//	EnableUPNP bool      `json:"enable_upnp"`
//	Label      string    `json:"label"`
//	Position   bool      `json:"position"`
//}

type CustomizationPostData struct {
	Origin       string    `json:"origin"`
	NetworkModel string    `json:"network_model"`
	Index        string    `json:"index"`
	Icon         string    `json:"icon"`
	Image        string    `json:"image"`
	Envs         EnvArray  `json:"envs"`
	Ports        PortArray `json:"ports"`
	Volumes      PathArray `json:"volumes"`
	Devices      PathArray `json:"devices"`
	//Port         string    `json:"port,omitempty"`
	PortMap     string `json:"port_map"`
	CpuShares   int64  `json:"cpu_shares"`
	Memory      int64  `json:"memory"`
	Restart     string `json:"restart"`
	EnableUPNP  bool   `json:"enable_upnp"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Position    bool   `json:"position"`
}
