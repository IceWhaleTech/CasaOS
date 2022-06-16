/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2021-09-30 18:18:14
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-16 18:05:00
 * @FilePath: /CasaOS/pkg/utils/jwt/jwt.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package jwt

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Id       int    `json:"id"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

//创建token
func GenerateToken(username, password string, id int, issuer string, t time.Duration) (string, error) {
	clims := Claims{
		username,
		password,
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, clims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err

}

//解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if clims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return clims, nil
		}
	}
	return nil, err
}
