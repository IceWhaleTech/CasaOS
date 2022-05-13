/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-03-18 11:40:55
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-13 14:48:01
 * @FilePath: /CasaOS/model/app-analyse.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type AppAnalyse struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	UUId     string `json:"uuid"`
	Language string `json:"language"`
	Version  string `json:"version"`
}

type ConnectionStatus struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Error string `json:"error"`
	UUId  string `json:"uuid"`
	Event string `json:"event"`
}
