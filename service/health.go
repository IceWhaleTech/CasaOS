package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/utils/port"
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
)

type HealthService interface {
	Services() (map[bool]*[]string, error)
	Ports() ([]int, []int, error)
}

type service struct{}

func (s *service) Services() (map[bool]*[]string, error) {
	services, err := systemctl.ListServices("casaos*")
	if err != nil {
		return nil, err
	}

	var running, notRunning []string

	for _, service := range services {
		if service.Running {
			running = append(running, service.Name)
		} else {
			notRunning = append(notRunning, service.Name)
		}
	}

	result := map[bool]*[]string{
		true:  &running,
		false: &notRunning,
	}

	return result, nil
}

func (s *service) Ports() ([]int, []int, error) {
	return port.ListPortsInUse()
}

func NewHealthService() HealthService {
	return &service{}
}
