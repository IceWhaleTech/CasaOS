package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"

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
	usedPorts := map[string]map[int]struct{}{
		"tcp": {},
		"udp": {},
	}

	for _, protocol := range []string{"tcp", "udp"} {
		filename := fmt.Sprintf("/proc/net/%s", protocol)

		file, err := os.Open(filename)
		if err != nil {
			return nil, nil, errors.New("Failed to open " + filename)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			localAddress := fields[1]
			addressParts := strings.Split(localAddress, ":")
			if len(addressParts) < 2 {
				continue
			}

			portHex := addressParts[1]
			port, err := strconv.ParseInt(portHex, 16, 0)
			if err != nil {
				continue
			}

			usedPorts[protocol][int(port)] = struct{}{}
		}

		if err := scanner.Err(); err != nil {
			return nil, nil, errors.New("Error reading from " + filename)
		}
	}

	return lo.Keys(usedPorts["tcp"]), lo.Keys(usedPorts["udp"]), nil
}

func NewHealthService() HealthService {
	return &service{}
}
