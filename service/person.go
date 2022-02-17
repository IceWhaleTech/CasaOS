package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
)

type PersonService interface {
	GetPersionInfo(token string) (m model.PersionModel, err error)
	Handshake(m model.ConnectState)
}

type personService struct {
}

var IpInfo model.PersionModel

func PushIpInfo(token string) {

	m := model.PersionModel{}
	m.Ips = GetDeviceAllIP()
	m.Token = token
	b, _ := json.Marshal(m)

	if reflect.DeepEqual(IpInfo, m) {
		return
	}
	head := make(map[string]string)
	infoS := httper2.Post(config.ServerInfo.Handshake+"/v1/update", b, "application/json", head)
	fmt.Println(infoS)
}
func (p *personService) GetPersionInfo(token string) (m model.PersionModel, err error) {
	infoS := httper2.Get(config.ServerInfo.Handshake+"/v1/ips/"+token, nil)
	err = json.Unmarshal([]byte(infoS), &m)
	return
}

//尝试连接
func (p *personService) Handshake(m model.ConnectState) {
	//1先进行udp打通成功

	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: 9901} //注意端口必须固定
	dstAddr := &net.UDPAddr{
		IP: net.ParseIP(config.ServerInfo.Handshake), Port: 9527}
	//DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"udp"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
	//(conn)UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(m)
	if _, err = conn.Write(b); err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 1024)
	//ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
	//ReadFromUDP方***在超过一个固定的时间点之后超时，并返回一个错误。
	n, _, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	toPersion := model.PersionModel{}
	err = json.Unmarshal(data[:n], &toPersion)
	if err != nil {
		fmt.Println(err)
	}
	// bidirectionHole(srcAddr, &anotherPeer)

	//2udp打洞成功向服务器汇报打洞结果
	//3转udp打洞

}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {

	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {

		fmt.Println("send handshake:", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from []")); err != nil {

				log.Println("send msg fail", err)
			}
		}
	}()

	for {

		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {

			log.Printf("error during read:%s\n", err)
		} else {

			log.Printf("收到数据：%s\n", data[:n])
		}
	}
}

func NewPersonService() PersonService {
	return &personService{}
}
