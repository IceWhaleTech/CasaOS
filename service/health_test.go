package service_test

import (
	"testing"

	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/stretchr/testify/assert"
)

func TestPorts(t *testing.T) {
	service := service.NewHealthService()

	tcpPorts, udpPorts, err := service.Ports()
	assert.NoError(t, err)

	assert.NotEmpty(t, tcpPorts)
	assert.NotEmpty(t, udpPorts)
}
