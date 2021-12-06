package ddns

import (
	"net/http"
)

func SetUserAgent(request *http.Request) {
	request.Header.Set("User-Agent", "Oasis lauren.pan@icewhale.org")
}

func SetContentType(request *http.Request, contentType string) {
	request.Header.Set("Content-Type", contentType)
}

func SetAccept(request *http.Request, acceptContent string) {
	request.Header.Set("Accept", acceptContent)
}

func SetAuthBearer(request *http.Request, token string) {
	request.Header.Set("Authorization", "Bearer "+token)
}

func SetAuthSSOKey(request *http.Request, key, secret string) {
	request.Header.Set("Authorization", "sso-key "+key+":"+secret)
}

func SetOauth(request *http.Request, value string) {
	request.Header.Set("oauth", value)
}

func SetXFilter(request *http.Request, value string) {
	request.Header.Set("X-Filter", value)
}
