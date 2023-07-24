package models

import (
	"regexp"
)

type TelevisorService struct {
	Name         string      `json:"name"`
	Dependents   []string    `json:"dependents"`
	Dependencies []string    `json:"dependencies"`
	Cpu          Utilization `json:"cpu"`
	Memory       Utilization `json:"memory"`
	Network      Utilization `json:"network"`
}

type Utilization struct {
	Quantile float64 `json:"quantile"`
	Mean     float64 `json:"mean"`
	Stdev    float64 `json:"stdev"`
}

func (service TelevisorService) IsServiceInContainer(containerName string) bool {
	serviceName := service.Name

	match, _ := regexp.MatchString(serviceName, containerName)

	return match
}
