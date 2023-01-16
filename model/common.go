package model

type PageResp struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}
