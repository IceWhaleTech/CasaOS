/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-01 18:32:57
 * @FilePath: /CasaOS/model/zima.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

import "time"

type Path struct {
	Name       string                 `json:"name"`   //File name or document name
	Path       string                 `json:"path"`   //Full path to file or folder
	IsDir      bool                   `json:"is_dir"` //Is it a folder
	Date       time.Time              `json:"date"`
	Size       int64                  `json:"size"` //File Size
	Type       string                 `json:"type,omitempty"`
	Label      string                 `json:"label,omitempty"`
	Write      bool                   `json:"write"`
	Extensions map[string]interface{} `json:"extensions"`
}
