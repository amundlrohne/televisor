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
