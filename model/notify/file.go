/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-26 14:21:57
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-30 18:51:36
 * @FilePath: /CasaOS/model/notify/file.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package notify

type File struct {
	Finished       bool   `json:"finished"`
	ProcessedSize  int64  `json:"processed_size"`
	ProcessingPath string `json:"processing_path"`
	Status         string `json:"status"`
	TotalSize      int64  `json:"total_size"`
	Id             string `json:"id"`
	To             string `json:"to"`
}
