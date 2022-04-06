package model

import "time"

type Path struct {
	Name  string    `json:"name"`   //File name or document name
	Path  string    `json:"path"`   //Full path to file or folder
	IsDir bool      `json:"is_dir"` //Is it a folder
	Date  time.Time `json:"date"`
	Size  int64     `json:"size"` //File Size
	Type  string    `json:"type,omitempty"`
	Label string    `json:"label,omitempty"`
}
