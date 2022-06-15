/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-15 16:53:00
 * @FilePath: /CasaOS/service/model/o_user.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

//Soon to be removed
type UserDBModel struct {
	Id          int    `gorm:"column:id;primary_key" json:"id"`
	UserName    string `json:"user_name"`
	Password    string `json:"password omitempty"`
	Role        string `json:"role"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	NickName    string `json:"nick_name"`
	Description string `json:"description"`
}

func (p *UserDBModel) UserDBModel() string {
	return "o_user"
}
