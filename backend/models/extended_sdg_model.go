package models

import (
	"regexp"
)

type TelevisorService struct {
	Name         string
	Dependents   []string
	Dependencies []string
	Cpu          Utilization
	Memory       Utilization
	Network      Utilization
}

type Utilization struct {
	P99 float64
}

func (service TelevisorService) IsServiceInContainer(containerName string) bool {
	serviceName := service.Name

	match, _ := regexp.MatchString(serviceName, containerName)

	return match
}
