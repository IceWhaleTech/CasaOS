package port

import (
	"fmt"
	"testing"
)

func TestPortAvailable(t *testing.T) {
	//	fmt.Println(PortAvailable())
	//fmt.Println(IsPortAvailable(6881,"tcp"))
	p, _ := GetAvailablePort("udp")
	fmt.Println("udp", p)
	fmt.Println(IsPortAvailable(p, "udp"))

	t1, _ := GetAvailablePort("tcp")
	fmt.Println("tcp", t1)
	fmt.Println(IsPortAvailable(t1, "tcp"))
}
