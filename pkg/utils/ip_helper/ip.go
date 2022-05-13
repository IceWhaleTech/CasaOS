package ip_helper

import (
	"net"
	"strings"

	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
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
func GetDeviceAllIP(port string) []string {
	var address []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return address
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To16() != nil {
				address = append(address, ipNet.IP.String()+":"+port)
			}
		}
	}
	return address
}

func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	ip.String()

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}
