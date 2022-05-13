package version

import (
	json2 "encoding/json"
	"strconv"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/tidwall/gjson"
)

func IsNeedUpdate() (bool, model.Version) {
	var version model.Version
	v := httper.OasisGet(config.ServerInfo.ServerApi + "/v1/sys/version")
	data := gjson.Get(v, "data")
	json2.Unmarshal([]byte(data.String()), &version)

	v1 := strings.Split(version.Version, ".")
	v2 := strings.Split(types.CURRENTVERSION, ".")

	for len(v1) < len(v2) {
		v1 = append(v1, "0")
	}
	for len(v2) < len(v1) {
		v2 = append(v2, "0")
	}
	for i := 0; i < len(v1); i++ {
		a, _ := strconv.Atoi(v1[i])
		b, _ := strconv.Atoi(v2[i])
		if a > b {
			return true, version
		}
	}
	return false, version
}

func IsClientNeedUpdate() (bool, model.Version) {
	var version model.Version
	v := httper.OasisGet(config.ServerInfo.ServerApi + "/v1/sys/client/version")
	data := gjson.Get(v, "data")
	json2.Unmarshal([]byte(data.String()), &version)

	v1 := strings.Split(version.Version, ".")
	v2 := strings.Split(types.CURRENTVERSION, ".")

	for len(v1) < len(v2) {
		v1 = append(v1, "0")
	}
	for len(v2) < len(v1) {
		v2 = append(v2, "0")
	}
	for i := 0; i < len(v1); i++ {
		a, _ := strconv.Atoi(v1[i])
		b, _ := strconv.Atoi(v2[i])
		if a > b {
			return true, version
		}
	}
	return false, version
}

//a版本大于b版本
func VersionCompared(a string, b string) bool {
	v1 := strings.Split(a, ".")
	v2 := strings.Split(b, ".")
	for len(v1) < len(v2) {
		v1 = append(v1, "0")
	}
	for len(v2) < len(v1) {
		v2 = append(v2, "0")
	}
	for i := 0; i < len(v1); i++ {
		a, _ := strconv.Atoi(v1[i])
		b, _ := strconv.Atoi(v2[i])
		if a > b {
			return true
		}
	}
	return false
}
