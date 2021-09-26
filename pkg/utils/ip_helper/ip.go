package ip_helper

import (
	"net"
	httper2 "oasis/pkg/utils/httper"
	"strings"
)

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

//获取外网ip
func GetExternalIPV4() string {
	return httper2.Get("https://api.ipify.org", nil)
}

//获取外网ip
func GetExternalIPV6() string {
	return httper2.Get("https://api6.ipify.org", nil)
}

//获取本地ip
func GetLoclIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return "127.0.0.1"
}
