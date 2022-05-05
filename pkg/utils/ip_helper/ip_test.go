package ip_helper

import (
	"fmt"
	"net"
	"testing"
)

func TestGetExternalIPV4(t *testing.T) {
	ipv4 := make(chan string)
	go func() { ipv4 <- GetExternalIPV4() }()
	fmt.Println(<-ipv4)
}
func TestGetExternalIPV6(t *testing.T) {
	ipv6 := make(chan string)
	go func() { ipv6 <- GetExternalIPV6() }()
	fmt.Println(<-ipv6)

}

func TestGetLoclIp(t *testing.T) {
	fmt.Println(GetLoclIp())
}
func TestHasLocalIP(t *testing.T) {
	fmt.Println("dddd")
	fmt.Println(HasLocalIP(net.ParseIP("192.168.2.10")))
}
