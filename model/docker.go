/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-12-08 18:10:25
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-14 17:20:36
 * @FilePath: /CasaOS/model/docker.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

type DockerStatsModel struct {
	Icon  string      `json:"icon"`
	Title string      `json:"title"`
	Data  interface{} `json:"data"`
	Pre   interface{} `json:"pre"`
}
