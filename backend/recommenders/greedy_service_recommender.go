package recommenders

import (
	"fmt"
	"strings"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
)

func NewGreedyServiceRecommender(services map[string]models.TelevisorService, operations models.Operations, annotation models.Annotation) (map[string]models.TelevisorService, models.Operations, models.Annotation) {
	servicesToBeMerged := []models.TelevisorService{}
	parentService := services[annotation.Services[len(annotation.Services)-1]]

	for _, serviceName := range annotation.Services {
		servicesToBeMerged = append(servicesToBeMerged, services[serviceName])
	}

	cpuUtils := utils.SumCPUUtilizations(servicesToBeMerged...)
	memUtils := utils.SumMemoryUtilizations(servicesToBeMerged...)

	if cpuUtils.Quantile > 0.4 || memUtils.Quantile > 0.4 {
		annotation.Recomendation = models.Recommendation{
			Message: fmt.Sprintf("Service is greedy however one or more quantile utilization metrics would exceed requirements if merged with %s", parentService.Name),
		}
		return services, operations, annotation
	}

	for _, service := range servicesToBeMerged {
		delete(services, service.Name)
	}

	for _, operation := range annotation.Operations {
		delete(operations[annotation.InitiatingOperation], operation)
	}

	mergedServiceName := strings.Join(annotation.Services, "/")
	services[mergedServiceName] = models.TelevisorService{
		Name:         mergedServiceName,
		Dependencies: parentService.Dependencies,
		Dependents:   utils.StringXOR(parentService.Dependents, annotation.Services),
		Cpu:          cpuUtils,
		Memory:       memUtils,
	}

	for requestKey, requestValue := range operations {
		for operationKey, operationValue := range requestValue {

			if operationValue.To == parentService.Name && operationValue.From == parentService.Name {
				operationValue.To = mergedServiceName
				operationValue.From = mergedServiceName
				operations[requestKey][operationKey] = operationValue
			} else if operationValue.To == parentService.Name {
				operationValue.To = mergedServiceName
				operations[requestKey][operationKey] = operationValue
			} else if operationValue.From == parentService.Name {
				operationValue.From = mergedServiceName
				operations[requestKey][operationKey] = operationValue
			}

		}
	}

	annotation.Recomendation = models.Recommendation{Message: fmt.Sprintf("Combinining services %s, into resulting service %s", annotation.Services, mergedServiceName)}

	return services, operations, annotation
}

func GreedyServiceRecommender(services map[string]models.TelevisorService, operations models.Operations, annotation models.Annotation) (map[string]models.TelevisorService, models.Operations, models.Annotation) {
	greedyService := services[annotation.Services[0]]
	operation := annotation.Operations[0]
	parentService := services[operations[annotation.InitiatingOperation][operation].From]

	cpuUtils := utils.SumUtilizations(greedyService.Cpu, parentService.Cpu)
	memUtils := utils.SumUtilizations(greedyService.Memory, parentService.Memory)
	netUtils := utils.SumUtilizations(greedyService.Network, parentService.Network)

	if cpuUtils.Quantile > 0.4 || memUtils.Quantile > 0.4 {
		annotation.Recomendation = models.Recommendation{
			Message: fmt.Sprintf("Service is greedy however one or more quantile utilization metrics would exceed requirements if merged with %s", parentService.Name),
		}
		return services, operations, annotation
	}

	delete(services, parentService.Name)
	delete(services, greedyService.Name)
	delete(operations[annotation.InitiatingOperation], operation)

	mergedService := fmt.Sprintf("%s/%s", parentService.Name, greedyService.Name)
	services[mergedService] = models.TelevisorService{
		Name:         mergedService,
		Dependencies: parentService.Dependencies,
		Dependents:   utils.FilterStringArray(parentService.Dependents, greedyService.Name),
		Cpu:          cpuUtils,
		Memory:       memUtils,
		Network:      netUtils,
	}

	for requestKey, requestValue := range operations {
		for operationKey, operationValue := range requestValue {
			if operationValue.To == parentService.Name || operationValue.To == greedyService.Name {
				operationValue.To = mergedService
				operations[requestKey][operationKey] = operationValue
			}
		}
	}

	annotation.Recomendation = models.Recommendation{Message: fmt.Sprintf("Combinining service %s and service %s, into resulting service %s", greedyService.Name, parentService.Name, mergedService)}

	return services, operations, annotation
}
