package zerotier

import (
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/tidwall/gjson"
	"net/http"
)

func PostData(url, token string, data string) interface{} {

	body, code := httper2.ZeroTierPostJson(url, data, GetHead(token))

	if code != http.StatusOK {
		return ""
	}
	result := gjson.Parse(body)
	return result.Value()
}

func GetData(url, token string) interface{} {

	body, code := httper2.ZeroTierGet(url, GetHead(token))

	if code != http.StatusOK {
		return ""
	}
	result := gjson.Parse(body)
	return result.Value()
}

func DeleteMember(url, token string) interface{} {

	body, code := httper2.ZeroTierDelete(url, GetHead(token))

	if code != http.StatusOK {
		return ""
	}
	result := gjson.Parse(body)
	return result.Value()
}

func GetHead(token string) map[string]string {
	var head = make(map[string]string)
	head["Authorization"] = "Bearer " + token
	head["Content-Type"] = "application/json"
	return head
}
