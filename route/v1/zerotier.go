package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/IceWhaleTech/CasaOS/common"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func ZerotierProxy(c *gin.Context) {
	// Read the port number from the file
	w := c.Writer
	r := c.Request
	port, err := ioutil.ReadFile("/var/lib/zerotier-one/zerotier-one.port")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the request path and remove "/zt"
	path := strings.TrimPrefix(r.URL.Path, "/v1/zt")
	fmt.Println(path)

	// Build the target URL
	targetURL := fmt.Sprintf("http://localhost:%s%s", strings.TrimSpace(string(port)), path)

	// Create a new request
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the X-ZT1-AUTH header
	authToken, err := ioutil.ReadFile("/var/lib/zerotier-one/authtoken.secret")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-ZT1-AUTH", strings.TrimSpace(string(authToken)))

	copyHeaders(req.Header, r.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func copyHeaders(destination, source http.Header) {
	for key, values := range source {
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}

func CheckNetwork() {
	//先获取所有已创建的网络
	respBody, err := httper.ZTGet("/controller/network")
	if err != nil {
		fmt.Println(err)
		return
	}
	networkId := ""
	address := ""
	networkNames := gjson.ParseBytes(respBody).Array()
	for _, v := range networkNames {
		res, err := httper.ZTGet("/controller/network/" + v.Str)
		if err != nil {
			fmt.Println(err)
			return
		}
		name := gjson.GetBytes(res, "name").Str
		if name == common.RANW_NAME {
			fmt.Println(string(res))
			networkId = gjson.GetBytes(res, "id").Str
			break
		}
	}
	if len(networkId) == 0 {
		if len(address) == 0 {
			address = GetAddress()
		}
		networkId = CreateNet(address)
	}
	res, err := httper.ZTGet("/network")
	if err != nil {
		fmt.Println(err)
		return
	}
	joined := false
	networks := gjson.GetBytes(res, "#.id").Array()
	for _, v := range networks {
		if v.Str == networkId {
			fmt.Println("已加入网络")
			joined = true
			break
		}
	}
	if !joined {
		JoinAndUpdateNet(address, networkId)
	}
}
func GetAddress() string {
	//获取nodeId
	nodeRes, err := httper.ZTGet("/status")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return gjson.GetBytes(nodeRes, "address").String()
}
func JoinAndUpdateNet(address, networkId string) {
	res, err := httper.ZTPost("/network/"+networkId, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))
	if len(address) == 0 {
		address = GetAddress()
	}
	b := `{
		"authorized": true,
		"activeBridge": true,
		"ipAssignments": [
		  "10.147.20.1"
		]
	  }`
	r, err := httper.ZTPost("/controller/network/"+networkId+"/member/"+address, b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(r))
}
func CreateNet(address string) string {
	body := `{
		"name": "` + common.RANW_NAME + `",
		"private": false,
		"v4AssignMode": {
		"zt": true
		},
		"ipAssignmentPools": [
		{
		"ipRangeStart": "10.147.20.1",
		"ipRangeEnd": "10.147.20.254"
		}
		],
		"routes": [
		{
		"target": "10.147.20.0/24"
		}
		],
		"rules": [
		{
		"etherType": 2048,
		"not": true,
		"or": false,
		"type": "MATCH_ETHERTYPE"
		},
		{
		"etherType": 2054,
		"not": true,
		"or": false,
		"type": "MATCH_ETHERTYPE"
		},
		{
		"etherType": 34525,
		"not": true,
		"or": false,
		"type": "MATCH_ETHERTYPE"
		},
		{
		"type": "ACTION_DROP"
		},
		{
		"type": "ACTION_ACCEPT"
		}
		]
		}`
	createRes, err := httper.ZTPost("/controller/network/"+address+"______", body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return gjson.GetBytes(createRes, "id").Str
}
