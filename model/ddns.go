package model

type GoDaddyModel struct {
	Type    uint   `json:"type"`
	ApiHost string `json:"api_host"`
	Key     string `json:"key"`
	Secret  string `json:"secret"`
	Host    string `json:"host"`
}
