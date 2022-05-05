package model

const CONTAINERTABLENAME = "o_container"

//Soon to be removed
type AppListDBModel struct {
	CustomId string `gorm:"column:custom_id;primary_key" json:"custom_id"`
	Title    string `json:"title"`
	//	ScreenshotLink model.Strings `gorm:"type:json" json:"screenshot_link,omitempty"`
	ScreenshotLink string `json:"screenshot_link"`
	Slogan         string `json:"slogan"`
	Description    string `json:"description"`
	//Tags           model.Strings `gorm:"type:json" json:"tags"`
	Tags        string `json:"tags"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	ContainerId string `json:"container_id,omitempty"`
	Image       string `json:"image,omitempty"`
	Index       string `json:"index"`
	CreatedAt   string `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt   string `gorm:"<-:create;<-:update;autoUpdateTime" json:"updated_at"`
	//Port           string           `json:"port,omitempty"`
	PortMap    string `json:"port_map"`
	Label      string `json:"label"`
	EnableUPNP bool   `json:"enable_upnp"`
	Envs       string `json:"envs"`
	Ports      string `json:"ports"`
	Volumes    string `json:"volumes"`
	Devices    string `json:"devices"`
	//Envs      []model.Env      `json:"envs"`
	//Ports     []model.PortMap  `gorm:"type:json" json:"ports"`
	//Volumes   []model.PathMap  `gorm:"type:json" json:"volumes"`
	//Devices   []model.PathMap  `gorm:"type:json" json:"device"`
	Position  bool   `json:"position"`
	NetModel  string `json:"net_model"`
	CpuShares int64  `json:"cpu_shares"`
	Memory    int64  `json:"memory"`
	Restart   string `json:"restart"`
	//Rely      model.MapStrings `gorm:"type:json" json:"rely"` //[{"mysql":"id"},{"mysql":"id"}]
	Origin     string `json:"origin"`
	HostName   string `json:"host_name"`
	Privileged bool   `json:"privileged"`
	CapAdd     string `json:"cap_add"`
	Cmd        string `gorm:"type:json" json:"cmd"`
}

func (p *AppListDBModel) TableName() string {
	return "o_container"
}

type MyAppList struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	State    string `json:"state"`
	CustomId string `gorm:"column:custom_id;primary_key" json:"custom_id"`
	Index    string `json:"index"`
	Port     string `json:"port"`
	UpTime   string `json:"up_time"`
	Slogan   string `json:"slogan"`
	Type     string `json:"type"`
	//Rely       model.MapStrings `json:"rely"` //[{"mysql":"id"},{"mysql":"id"}]
	Image      string `json:"image"`
	Volumes    string `json:"volumes"`
	NewVersion bool   `json:"new_version"`
}
