/*
 * @Author: a624669980@163.com a624669980@163.com
 * @Date: 2022-12-13 11:05:05
 * @LastEditors: a624669980@163.com a624669980@163.com
 * @LastEditTime: 2022-12-13 11:05:13
 * @FilePath: /drive/internal/driver/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package driver

type Config struct {
	Name        string `json:"name"`
	LocalSort   bool   `json:"local_sort"`
	OnlyLocal   bool   `json:"only_local"`
	OnlyProxy   bool   `json:"only_proxy"`
	NoCache     bool   `json:"no_cache"`
	NoUpload    bool   `json:"no_upload"`
	NeedMs      bool   `json:"need_ms"` // if need get message from user, such as validate code
	DefaultRoot string `json:"default_root"`
	CheckStatus bool
}

func (c Config) MustProxy() bool {
	return c.OnlyProxy || c.OnlyLocal
}
