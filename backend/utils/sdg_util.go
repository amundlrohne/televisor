package utils

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
	"github.com/jaegertracing/jaeger/model"
)

func ExtractServicesFromSDG(sdg []model.DependencyLink) map[string]*models.ExtendedService {
	var services = make(map[string]*models.ExtendedService)

	for _, value := range sdg {
		var parentDeps []string = services[value.Parent].Dependents
		services[value.Parent].Dependents = append(parentDeps, value.Child)

		var childDeps = services[value.Child].Dependencies
		services[value.Child].Dependencies = append(childDeps, value.Parent)
	}

	fmt.Println(services)

	return services
}

func CPUUtilizationToSDG(data models.PrometheusAPIResponse, sdg map[string]*models.ExtendedService) {
	for _, d := range data.Data.Result {
		if service, ok := sdg[d.Metric.Name]; ok {
			service.Cpu.P99 = d.Value[1].(float64) // Needs testing
			sdg[d.Metric.Name] = service
		}
	}
}
