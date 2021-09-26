package upnp

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

//
////添加一个端口映射
func (n *Upnp)AddPortMapping(localPort, remotePort int, protocol string) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			//log.Println("upnp模块报错了", errTemp)
			err = errTemp.(error)
		}
	}(err)
	if issuccess := addSend(localPort, remotePort, protocol,n.GatewayHost, n.CtrlUrl,n.LocalHost); issuccess {
		return nil
	} else {
		return errors.New("添加一个端口映射失败")
	}
	return
}

func addSend(localPort, remotePort int, protocol, host, ctrUrl,localHost string) bool {
	request := addRequest(localPort, remotePort, protocol, host, ctrUrl,localHost)
	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()
	//resultBody, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(resultBody))
	if response.StatusCode == 200 {
		return true
	}

	return false
}

type Node struct {
	Name    string
	Content string
	Attr    map[string]string
	Child   []Node
}

func addRequest(localPort, remotePort int, protocol string, gatewayHost, ctlUrl,localHost string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:AddPortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}

	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewInternalPort", Content: strconv.Itoa(localPort)}
	childList3 := Node{Name: "NewProtocol", Content: protocol}
	childList4 := Node{Name: "NewEnabled", Content: "1"}
	childList5 := Node{Name: "NewInternalClient", Content: localHost}
	childList6 := Node{Name: "NewLeaseDuration", Content: "0"}
	childList7 := Node{Name: "NewPortMappingDescription", Content: "Oasis"}
	childList8 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childTwo.AddChild(childList4)
	childTwo.AddChild(childList5)
	childTwo.AddChild(childList6)
	childTwo.AddChild(childList7)
	childTwo.AddChild(childList8)

	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()
	//请求
	request, _ := http.NewRequest("POST", "http://"+gatewayHost+ctlUrl,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

func (n *Node) AddChild(node Node) {
	n.Child = append(n.Child, node)
}

func (n *Node) BuildXML() string {
	buf := bytes.NewBufferString("<")
	buf.WriteString(n.Name)
	for key, value := range n.Attr {
		buf.WriteString(" ")
		buf.WriteString(key + "=" + value)
	}
	buf.WriteString(">" + n.Content)

	for _, node := range n.Child {
		buf.WriteString(node.BuildXML())
	}
	buf.WriteString("</" + n.Name + ">")
	return buf.String()
}

func (n *Upnp)DelPortMapping(remotePort int, protocol string) bool {
	issuccess := delSendSend(remotePort, protocol,n.GatewayHost,n.CtrlUrl)
	if issuccess {
		//this.MappingPort.delMapping(remotePort, protocol)
		//fmt.Println("删除了一个端口映射： remote:", remotePort)
	}
	return issuccess
}

func delSendSend(remotePort int, protocol,host,ctlUrl string) bool {
	delrequest := delbuildRequest(remotePort, protocol,host,ctlUrl)
	response, _ := http.DefaultClient.Do(delrequest)
	//resultBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if response.StatusCode == 200 {
		// log.Println(string(resultBody))
		return true
	}
	return false
}
func delbuildRequest(remotePort int, protocol,host,ctlUrl string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:DeletePortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewProtocol", Content: protocol}
	childList3 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	//请求
	request, _ := http.NewRequest("POST", "http://"+host+ctlUrl,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}