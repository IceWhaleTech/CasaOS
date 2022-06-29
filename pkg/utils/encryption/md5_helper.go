/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-14 14:33:25
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-14 14:33:49
 * @FilePath: /CasaOS/pkg/utils/encryption/md5_helper.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5ByStr(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
