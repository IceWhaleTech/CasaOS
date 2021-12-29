package env_helper

import "strings"

func ReplaceDefaultENV(key, tz string) string {
	temp := ""
	switch key {
	case "$DefaultPassword":
		temp = "casaos"
	case "$DefaultUserName":
		temp = "admin"

	case "$PUID":
		temp = "1000"
	case "$PGID":
		temp = "1000"
	case "$TZ":
		temp = tz
	}
	return temp
}

//replace env default setting
func ReplaceStringDefaultENV(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "$DefaultPassword", ReplaceDefaultENV("$DefaultPassword", "")), "$DefaultUserName", ReplaceDefaultENV("$DefaultUserName", ""))
}
