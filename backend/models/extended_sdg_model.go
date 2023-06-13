package models

import (
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
