package model

import "time"

type Path struct {
	Name  string    `json:"name"`
	Path  string    `json:"path"`
	IsDir bool      `json:"is_dir"`
	Date  time.Time `json:"date"`
	Size  int64     `json:"size"`
}
