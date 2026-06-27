package service_test

import (
	"runtime"
	"testing"

	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/stretchr/testify/assert"
)

func TestPorts(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("port.ListPortsInUse reads /proc/net, which is only available on linux")
	}

	service := service.NewHealthService()

	tcpPorts, udpPorts, err := service.Ports()
	assert.NoError(t, err)

	assert.NotEmpty(t, tcpPorts)
	assert.NotEmpty(t, udpPorts)
}
