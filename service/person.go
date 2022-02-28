package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
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

	//websocket 连接
	// bidirectionHole(srcAddr, &anotherPeer)

	//2udp打洞成功向服务器汇报打洞结果
	//3转udp打洞

}

var ipAddress chan string

func UDPConnect(ips []string) {
	ipAddress = make(chan string)
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: 9901}

	conn, err := net.ListenUDP("udp", srcAddr)

	if err != nil {
		fmt.Println("监听错误", err.Error())
	}
	for _, v := range ips {
		dstAddr := &net.UDPAddr{
			IP: net.ParseIP(v), Port: 9901}

		fmt.Println(v, "开始监听")
		go AsyncUDPConnect(conn, dstAddr)
	}

	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read:%s\n", err)
		} else {
			fmt.Println("收到数据：", string(data[:n]))
		}
	}
}

func AsyncUDPConnect(conn *net.UDPConn, dst *net.UDPAddr) {
	for {
		time.Sleep(2 * time.Second)
		if _, err := conn.WriteToUDP([]byte(dst.IP.String()+" is ok"), dst); err != nil {
			log.Println("send msg fail", err)
			return
		} else {
			fmt.Println(dst.IP)
			fmt.Println(dst.IP.To4())
		}
	}
}
func TestTCPOne() {

	for i := 0; i < 100; i++ {
		fmt.Println(httper.Get("http://18.136.202.206:8088/v1/ping", nil))
		time.Sleep(time.Second * 2)
	}

}

func TCPServer() {
	localAddress := net.TCPAddr{IP: net.IPv4zero, Port: 8087}  //定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", &localAddress) //在刚定义好的地址上进监听请求。
	if err != nil {
		fmt.Println("监听出错：", err)
		return
	}
	defer func() { //担心return之前忘记关闭连接，因此在defer中先约定好关它。
		tcpListener.Close()
	}()
	fmt.Println("正在等待连接...")
	var conn, err2 = tcpListener.AcceptTCP() //接受连接。

	if err2 != nil {
		fmt.Println("接受连接失败：", err2)
		return
	}
	var remoteAddr = conn.RemoteAddr() //获取连接到的对像的IP地址。
	fmt.Println("接受到一个连接：", remoteAddr)
	fmt.Println("正在读取消息...")
	var buf = make([]byte, 1000)
	var n, _ = conn.Read(buf) //读取对方发来的内容。
	fmt.Println("接收到客户端的消息：", string(buf[:n]))
	conn.Write([]byte("hello, Nice to meet you, my name is SongXingzhu")) //尝试发送消息。
	conn.Close()
}

func parseAddrTCP(addr string) net.TCPAddr {

	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.TCPAddr{

		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func TestTCPTwo() {
	//localAddress := net.TCPAddr{IP: net.IPv4zero, Port: 8087} //定义一个本机IP和端口。
	//	t, _ := net.ResolveTCPAddr("tcp", "18.136.202.206:8088")

	//	dstAddr := &net.TCPAddr{IP: net.ParseIP("18.136.202.206"), Port: 8088}
	connTCP, err := net.ResolveTCPAddr("tcp", "18.136.202.206:8088")
	//	ddd,err := net.Dial("tcp", "")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(connTCP)
	// connTCP.Write([]byte("test"))
	// var buf = make([]byte, 1000)
	// var n, _ = connTCP.Read(buf) //读取对方发来的内容。
	// anotherPeer := parseAddrTCP(string(buf[:n]))

	// fmt.Println("接收到消息：", anotherPeer)
	// connTCP.Close()
	// time.Sleep(time.Second * 20)
	//go TCPServer()
	//	bidirectionHoleTCP(&localAddress, &anotherPeer)
}

func bidirectionHoleTCP(srcAddr *net.TCPAddr, anotherAddr *net.TCPAddr) {
	//	t, _ := net.ResolveTCPAddr("tcp", srcAddr.String())
	conn, err := net.Dial("tcp", anotherAddr.String())
	if err != nil {

		fmt.Println("send handshake:", err)
	}
	go func() {

		for {

			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from " + config.ServerInfo.Token)); err != nil {

				log.Println("send msg fail", err)
			}
		}
	}()

	for {

		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {

			log.Printf("error during read:%s\n", err)
		} else {

			log.Printf("收到数据：%s\n", data[:n])
		}
	}
}

func TestTCP() {

	conn, err := net.Dial("tcp", "192.168.2.224:8088")

	// srcAddr := &net.TCPAddr{
	// 	IP: net.IPv4zero, Port: 9901}
	// conn, err := net.ListenTCP("tcp", srcAddr)
	// con, err := conn.AcceptTCP()
	// 连接出错则打印错误消息并退出程序
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	time.Sleep(time.Second * 2)
	// 调用返回的连接对象提供的 Write 方法发送请求
	for i := 0; i < 10; i++ {
		n, err := conn.Write([]byte("aaaa"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}

		fmt.Println(n)
		time.Sleep(time.Second)
	}

	// 通过连接对象提供的 Read 方法读取所有响应数据
	result, err := readFully(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	// 打印响应数据
	fmt.Println(string(result))
	os.Exit(0)
}

func readFully(conn net.Conn) ([]byte, error) {
	// 读取所有响应数据后主动关闭连接
	defer conn.Close()
	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}

func GetUdpConnet() {
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9901}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("18.136.202.206"), Port: 9527}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {

		fmt.Println(err)
	}

	if _, err = conn.Write([]byte("hello,I'm new peer:" + config.ServerInfo.Token)); err != nil {

		fmt.Println("写入错误", err)
	}
	time.Sleep(time.Second)

	data := make([]byte, 1024)
	//ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
	//ReadFromUDP方***在超过一个固定的时间点之后超时，并返回一个错误。
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	fmt.Println(remoteAddr)
	fmt.Println("服务器返回的信息", string(data[:n]))
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer)
	bidirectionHole(&anotherPeer)
}

func bidirectionHole(anotherAddr *net.UDPAddr) {
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: 9901}
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {

		fmt.Println("send handshake:", err)
	}
	go func() {

		for {

			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + config.ServerInfo.Token + "]")); err != nil {

				log.Println("send msg fail", err)
			}
		}
		fmt.Println("退出")
	}()

	for {

		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read:%s\n", err)
		} else {
			log.Printf("本机token：%s\n", config.ServerInfo.Token)
			log.Printf("收到数据：%s\n", data[:n])
		}
	}

}
func parseAddr(addr string) net.UDPAddr {

	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{

		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
func NewPersonService() PersonService {
	return &personService{}
}
