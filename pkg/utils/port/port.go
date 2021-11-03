package port

import (
	"fmt"
	"net"
)

// 获取可用端口
func GetAvailablePort(t string) (int, error) {
	address := fmt.Sprintf("%s:0", "0.0.0.0")
	if t == "udp" {
		add, err := net.ResolveUDPAddr(t, address)
		if err != nil {
			return 0, err
		}

		listener, err := net.ListenUDP(t, add)
		if err != nil {
			return 0, err
		}

		defer listener.Close()
		return listener.LocalAddr().(*net.UDPAddr).Port, nil
	} else {
		add, err := net.ResolveTCPAddr(t, address)
		if err != nil {
			return 0, err
		}

		listener, err := net.ListenTCP(t, add)
		if err != nil {
			return 0, err
		}

		defer listener.Close()
		return listener.Addr().(*net.TCPAddr).Port, nil
	}

}

// 判断端口是否可以（未被占用）
// param t tcp/udp
func IsPortAvailable(port int, t string) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	if t == "udp" {
		sadd, err := net.ResolveUDPAddr("udp", address)
		uc, err := net.ListenUDP("udp", sadd)

		if err != nil {
			fmt.Println(err.Error())
			return false
		} else {
			defer uc.Close()
			return true
		}

	} else {
		listener, err := net.Listen(t, address)

		if err != nil {
			//log.Infof("port %s is taken: %s", address, err)
			return false
		}
		defer listener.Close()
		return true
	}

}
