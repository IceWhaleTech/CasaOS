package upnp

import (
	ip_helper2 "oasis/pkg/utils/ip_helper"
	"testing"
)

func TestGetCtrlUrl(t *testing.T) {
	upnp, err := Gateway()
	if err == nil {
		upnp.CtrlUrl = GetCtrlUrl(upnp.GatewayHost, upnp.DeviceDescUrl)
		upnp.LocalHost = ip_helper2.GetLoclIp()
		upnp.AddPortMapping(8090, 8090, "TCP")
		//upnp.DelPortMapping(9999,  "TCP")
	}
}
