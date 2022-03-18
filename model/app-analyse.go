package model

type AppAnalyse struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	UUId     string `json:"uuid"`
	Language string `json:"language"`
}

type ConnectionStatus struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Error string `json:"error"`
	UUId  string `json:"uuid"`
	Event string `json:"event"`
}
