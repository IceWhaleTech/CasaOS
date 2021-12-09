package model

type DockerStatsModel struct {
	Icon  string      `json:"icon"`
	Title string      `json:"title"`
	Data  interface{} `json:"data"`
	Pre   interface{} `json:"pre"`
}
