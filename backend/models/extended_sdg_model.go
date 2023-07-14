package models

import (
	"regexp"
	"time"

	"github.com/jaegertracing/jaeger/model"
)

type ExtendedSDG struct {
	model.DependencyLink

	AvgLatency time.Duration
	P99Latency time.Duration
}

type ExtendedService struct {
	Name         string
	Dependents   []string
	Dependencies []string
	Cpu          Utiliziation
	Memory       Utiliziation
	Network      Utiliziation
}

type Utiliziation struct {
	P99 float64
}

func (service ExtendedService) IsServiceInContainer(containerName string) bool {
	serviceName := service.Name

	match, _ := regexp.MatchString(serviceName, containerName)

	return match
}

type TelevisorService struct {
	Name       string
	Operations []string
	//Traces        map[OperationName][]model.Span
	Relationships map[string]TelevisorRelationship
}

type TelevisorRelationship struct {
	Count       int
	ServiceName string
}
