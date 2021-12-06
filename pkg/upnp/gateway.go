package upnp

import (
	ip_helper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/ip_helper"
	"github.com/pkg/errors"
	"net"
	"strings"
)

func Gateway() (*Upnp, error) {
	result, error := send()
	if result == "" || error != nil {
		return nil, error
	}
	upnp := resolvesss(result)
	return upnp, nil
}

func send() (string, error) {
	var str = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
	var conn *net.UDPConn
	remoteAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		return "", errors.New("组播地址格式不正确")
	}
	localAddr, err := net.ResolveUDPAddr("udp", ip_helper2.GetLoclIp()+":")

	if err != nil {
		return "", errors.New("本地ip地址格式不正确")
	}
	conn, err = net.ListenUDP("udp", localAddr)
	defer conn.Close()
	if err != nil {
		return "", errors.New("监听udp出错")
	}
	_, err = conn.WriteToUDP([]byte(str), remoteAddr)
	if err != nil {
		return "", errors.New("发送msg到组播地址出错")
	}
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return "", errors.New("从组播地址接搜消息出错")
	}
	result := string(buf[:n])
	return result, nil
}

func resolvesss(result string) *Upnp {
	var upnp = &Upnp{}
	lines := strings.Split(result, "\r\n")
	for _, line := range lines {
		//按照第一个冒号分为两个字符串
		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			//fmt.Println(nameValues[1])
		case "CACHE-CONTROL":
			//fmt.Println(nameValues[1])
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			upnp.GatewayHost = (urls[0])
			upnp.DeviceDescUrl = ("/" + urls[1])
		case "SERVER":
			upnp.GatewayName = (nameValues[1])
		default:
		}
	}
	return upnp
}
