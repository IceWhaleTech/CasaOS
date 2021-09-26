package ddns

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"oasis/pkg/ddns"
	"time"
)

type GoDaddyService interface {
	Update(ctx context.Context, client *http.Client)
}
type GoDaddy struct {
	Host    string `json:"host"`
	Key     string `json:"key"`
	Secret  string `json:"secret"`
	Domain  string `json:"domain"`
	IPV4    string `json:"ipv_4"`
	IPV6    string `json:"ipv_6"`
	ApiHost string `json:"api_host"`
}

func (g *GoDaddy) Update() string {
	client := &http.Client{Timeout: 30 * time.Second}
	recordType := ddns.A
	buffer1 := bytes.NewBuffer(nil)
	buffer1.WriteString(`[{"data":"`)
	buffer1.WriteString(g.IPV4)
	buffer1.WriteString(`"}]`)
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/domains/%s/records/%s/%s", g.ApiHost, g.Domain, recordType, g.Host), buffer1)
	if err != nil {
		return err.Error()
	}
	g.setHead(request)
	response, err := client.Do(request)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()
	if len(b) > 0 {
		r := gjson.GetBytes(b, "message")
		return r.String()
	}
	if len(g.IPV6) > 0 {
		recordType = ddns.AAAA
		buffer1 := bytes.NewBuffer(nil)
		buffer1.WriteString(`[{"data":"`)
		buffer1.WriteString(g.IPV6)
		buffer1.WriteString(`"}]`)
		request6, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/domains/%s/records/%s/%s", ddns.GODADDYAPIURL, g.Domain, recordType, g.Host), buffer1)
		if err != nil {
			return err.Error()
		}
		g.setHead(request6)
		response6, err := client.Do(request6)
		if err != nil {
			return err.Error()
		}
		defer response6.Body.Close()

		d, err := ioutil.ReadAll(response6.Body)
		if err != nil {
			return err.Error()
		}
		if len(d) > 0 {
			r := gjson.GetBytes(d, "message")
			return r.String()
		}
	}
	return ""
}

func (g *GoDaddy) setHead(request *http.Request) {
	SetUserAgent(request)
	SetAuthSSOKey(request, g.Key, g.Secret)
	SetContentType(request, "application/json")
	SetAccept(request, "application/json")
}
