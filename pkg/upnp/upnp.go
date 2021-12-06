package upnp

import (
	"fmt"
	"github.com/prestonTao/upnp"
)

type Upnp struct {
	LocalHost     string `json:"local_host"`
	GatewayName   string `json:"gateway_name"`    //网关名称
	GatewayHost   string `json:"gateway_host"`    //网关ip和端口
	DeviceDescUrl string `json:"device_desc_url"` //设备描述url
	CtrlUrl       string `json:"ctrl_url"`        //控制请求url
}

func Testaaa() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.SearchGateway()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("local ip address: ", upnpMan.LocalHost)
		fmt.Println("gateway ip address: ", upnpMan.Gateway.Host)
	}
}
