package env_helper

import "strings"

func ReplaceDefaultENV(key string) string {
	temp := ""
	switch key {
	case "$DefaultPassword":
		temp = "casaos"
	case "$DefaultUserName":
		temp = "admin"
	}
	return temp
}

//replace env default setting
func ReplaceStringDefaultENV(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "$DefaultPassword", ReplaceDefaultENV("$DefaultPassword")), "$DefaultUserName", ReplaceDefaultENV("$DefaultUserName"))
}
