/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-15 11:30:47
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-23 18:40:40
 * @FilePath: /CasaOS/model/system_model/verify_information.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package system_model

type VerifyInformation struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
