package docker

import "strings"

func GetDir(id, envName string) string {

	if strings.Contains(envName, "$AppID") && len(id) > 0 {
		return strings.ReplaceAll(envName, "$AppID", id)
	}
	return envName
}
