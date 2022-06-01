/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-27 15:01:58
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-05-31 14:51:21
 * @FilePath: /CasaOS/model/notify/application.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package notify

type Application struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	Message  string `json:"message"`
	Finished bool   `json:"finished"`
	Success  bool   `json:"success"`
}
