/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-15 11:30:47
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 15:10:10
 * @FilePath: /CasaOS/model/system_model/verify_information.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package system_model

type VerifyInformation struct {
	RefreshToken string `json:"refresh"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    string `json:"expires_at"`
}
