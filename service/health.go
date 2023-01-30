package service

import (
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
)

type HealthService interface {
	Services() (map[bool]*[]string, error)
}

type service struct{}

func (s *service) Services() (map[bool]*[]string, error) {
	services, err := systemctl.ListServices("casaos-*")
	if err != nil {
		return nil, err
	}

	var running, notRunning []string

	for _, service := range services {
		if service.Status == "running" {
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

func NewHealthService() HealthService {
	return &service{}
}
