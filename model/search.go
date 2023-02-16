package model

type SearchEngine struct {
	Name      string   `json:"name"`
	Icon      string   `json:"icon"`
	SearchUrl string   `json:"search_url"`
	RecoUrl   string   `json:"reco_url"`
	Data      []string `json:"data"`
}
