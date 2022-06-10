/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-20 16:27:12
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-09 18:18:46
 * @FilePath: /CasaOS/model/file.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type FileOperate struct {
	Type          string     `json:"type" binding:"required"`
	Item          []FileItem `json:"item" binding:"required"`
	TotalSize     int64      `json:"total_size"`
	ProcessedSize int64      `json:"processed_size"`
	To            string     `json:"to" binding:"required"`
	Style         string     `json:"style"`
	Finished      bool       `json:"finished"`
}

type FileItem struct {
	From          string `json:"from" binding:"required"`
	Finished      bool   `json:"finished"`
	Size          int64  `json:"size"`
	ProcessedSize int64  `json:"processed_size"`
}

type FileUpdate struct {
	FilePath    string `json:"path" binding:"required"`
	FileContent string `json:"content" binding:"required"`
}
