package model

type ZeroTierUpData struct {
	Config ZeroTierConfig `json:"config"`
}

type ZeroTierConfig struct {
	Private bool `json:"private"`
}
