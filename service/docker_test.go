package service

import (
	"fmt"
	"testing"
)

//func TestDockerImageInfo(t *testing.T) {
//	//DockerImageInfo()
//
//	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
//	if err != nil {
//		fmt.Println(0, err)
//	}
//
//	listener, err := net.ListenTCP("tcp", address)
//	if err != nil {
//		fmt.Println(0, err)
//	}
//
//	defer listener.Close()
//	fmt.Println(listener.Addr().(*net.TCPAddr).Port, nil)
//
//}

//func TestDockerNetwork(t *testing.T) {
//	DockerNetwork()
//}
//
//func TestDockerPull(t *testing.T) {
//	DockerPull()
//}
//
//func TestDockerLog(t *testing.T) {
//	DockerLog()
//}
//func TestDockerLogs(t *testing.T) {
//	DockerLogs()
//}

func TestDockerContainerStats(t *testing.T) {
	fmt.Println(DockerContainerStats1())
}

//func TestDockerImageRemove(t *testing.T) {
//	host, domain, tld := gotld.GetSubdomain("aaa.liru-05.top", 1)
//	fmt.Println(host)
//	fmt.Println(domain)
//	fmt.Println(tld)
//}
