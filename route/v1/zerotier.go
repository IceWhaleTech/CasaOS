package v1

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/common"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func ZerotierProxy(ctx echo.Context) error {
	// Read the port number from the file
	w := ctx.Response().Writer
	r := ctx.Request()
	port, err := ioutil.ReadFile("/var/lib/zerotier-one/zerotier-one.port")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}

	// Add the X-ZT1-AUTH header
	authToken, err := ioutil.ReadFile("/var/lib/zerotier-one/authtoken.secret")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	req.Header.Set("X-ZT1-AUTH", strings.TrimSpace(string(authToken)))

	copyHeaders(req.Header, r.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return the response to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
	// TODO
	return nil
}

func copyHeaders(destination, source http.Header) {
	for key, values := range source {
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}

func CheckNetwork() {
	logger.Info("start check network")
	respBody, err := httper.ZTGet("/controller/network")
	if err != nil {
		logger.Error("get network error", zap.Error(err))
		return
	}
	networkId := ""
	address := ""
	networkNames := gjson.ParseBytes(respBody).Array()
	routers := ""
	for _, v := range networkNames {
		res, err := httper.ZTGet("/controller/network/" + v.Str)
		if err != nil {
			logger.Error("get network error", zap.Error(err))
			return
		}
		name := gjson.GetBytes(res, "name").Str
		if name == common.RANW_NAME {
			networkId = gjson.GetBytes(res, "id").Str
			routers = gjson.GetBytes(res, "routes.0.target").Str
			break
		}
	}
	ip, s, e, c := getZTIP(routers)
	logger.Info("ip", zap.Any("ip", ip))
	if len(networkId) == 0 {
		if len(address) == 0 {
			address = GetAddress()
		}
		networkId = CreateNet(address, s, e, c)
	}
	res, err := httper.ZTGet("/network")
	if err != nil {
		logger.Error("get network error", zap.Error(err))
		return
	}
	joined := false
	networks := gjson.GetBytes(res, "#.id").Array()
	for _, v := range networks {
		if v.Str == networkId {
			joined = true
			break
		}
	}
	logger.Info("joined", zap.Any("joined", joined))
	if !joined {
		JoinAndUpdateNet(address, networkId, ip)
	}
}

func GetAddress() string {
	nodeRes, err := httper.ZTGet("/status")
	if err != nil {
		logger.Error("get status error", zap.Error(err))
		return ""
	}
	return gjson.GetBytes(nodeRes, "address").String()
}

func JoinAndUpdateNet(address, networkId, ip string) {
	logger.Info("start join network", zap.Any("ip", ip))
	_, err := httper.ZTPost("/network/"+networkId, "")
	if err != nil {
		logger.Error(" get network error", zap.Error(err))
		return
	}

	if len(address) == 0 {
		address = GetAddress()
	}
	b := `{
		"authorized": true,
		"activeBridge": true,
		"ipAssignments": [
		  "` + ip + `"
		]
	  }`
	_, err = httper.ZTPost("/controller/network/"+networkId+"/member/"+address, b)
	if err != nil {
		logger.Error("join network error", zap.Error(err))
		return
	}
}

func CreateNet(address, s, e, c string) string {
	body := `{
		"name": "` + common.RANW_NAME + `",
		"private": false,
		"v4AssignMode": {
		"zt": true
		},
		"ipAssignmentPools": [
		{
		"ipRangeStart": "` + s + `",
		"ipRangeEnd": "` + e + `"
		}
		],
		"routes": [
		{
		"target": "` + c + `"
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
		],
		"v6AssignMode": {
			"rfc4193": true
		   }
		}`
	createRes, err := httper.ZTPost("/controller/network/"+address+"______", body)
	if err != nil {
		logger.Error("post network error", zap.Error(err))
		return ""
	}
	return gjson.GetBytes(createRes, "id").Str
}

func GetZTIPs() []gjson.Result {
	res, err := httper.ZTGet("/network")
	if err != nil {
		logger.Error("get network error", zap.Error(err))
		return []gjson.Result{}
	}
	a := gjson.GetBytes(res, "#.routes.0.target")
	return a.Array()
}

func getZTIP(routes string) (ip, start, end, cidr string) {
	excluded := GetZTIPs()
	cidrs := []string{
		"10.147.11.0/24",
		"10.147.12.0/24",
		"10.147.13.0/24",
		"10.147.14.0/24",
		"10.147.15.0/24",
		"10.147.16.0/24",
		"10.147.17.0/24",
		"10.147.18.0/24",
		"10.147.19.0/24",
		"10.147.20.0/24",
		"10.240.0.0/16",
		"10.241.0.0/16",
		"10.242.0.0/16",
		"10.243.0.0/16",
		"10.244.0.0/16",
		"10.245.0.0/16",
		"10.246.0.0/16",
		"10.247.0.0/16",
		"10.248.0.0/16",
		"10.249.0.0/16",
		"172.21.0.0/16",
		"172.22.0.0/16",
		"172.23.0.0/16",
		"172.24.0.0/16",
		"172.25.0.0/16",
		"172.26.0.0/16",
		"172.27.0.0/16",
		"172.28.0.0/16",
		"172.29.0.0/16",
		"172.30.0.0/16",
	}
	filteredCidrs := make([]string, 0)
	if len(routes) > 0 {
		filteredCidrs = append(filteredCidrs, routes)
	} else {
		for _, cidr := range cidrs {
			isExcluded := false
			for _, excludedIP := range excluded {
				if cidr == excludedIP.Str {
					isExcluded = true
					break
				}
			}
			if !isExcluded {
				filteredCidrs = append(filteredCidrs, cidr)
			}
		}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ip = ""
	if len(filteredCidrs) > 0 {
		randomIndex := rnd.Intn(len(filteredCidrs))
		selectedCIDR := filteredCidrs[randomIndex]
		_, ipNet, err := net.ParseCIDR(selectedCIDR)
		if err != nil {
			logger.Error("ParseCIDR error", zap.Error(err))
			return
		}
		cidr = selectedCIDR
		startIP := ipNet.IP
		endIP := make(net.IP, len(startIP))
		copy(endIP, startIP)

		for i := range startIP {
			endIP[i] |= ^ipNet.Mask[i]
		}
		startIP[3] = 1
		start = startIP.String()
		endIP[3] = 254
		end = endIP.String()
		ipt := ipNet
		ipt.IP[3] = 1
		ip = ipt.IP.String()
		return
	} else {
		logger.Error("No available CIDR found")
	}
	return
}
