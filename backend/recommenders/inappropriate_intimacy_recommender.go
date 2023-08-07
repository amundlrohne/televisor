package recommenders

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
)

func InappropriateIntimacyRecommender(services map[string]models.TelevisorService, operations models.Operations, annotation models.Annotation) (map[string]models.TelevisorService, models.Operations, models.Annotation) {
	operation := operations[annotation.InitiatingOperation]
	servicesToBeMerged := []models.TelevisorService{}

	for _, serviceName := range annotation.Services {
		if _, ok := services[serviceName]; ok {
			servicesToBeMerged = append(servicesToBeMerged, services[serviceName])
			delete(services, serviceName)
		}
	}

	mergedService := models.TelevisorService{
		Name:    servicesToBeMerged[0].Name,
		Cpu:     servicesToBeMerged[0].Cpu,
		Memory:  servicesToBeMerged[0].Memory,
		Network: servicesToBeMerged[0].Network,
	}

	for _, service := range servicesToBeMerged[1:] {
		mergedService = models.TelevisorService{
			Name:    fmt.Sprintf("%s/%s", mergedService.Name, service.Name),
			Cpu:     utils.SumUtilizations(mergedService.Cpu, service.Cpu),
			Memory:  utils.SumUtilizations(mergedService.Memory, service.Memory),
			Network: utils.SumUtilizations(mergedService.Network, service.Network),
		}
	}

	services[mergedService.Name] = mergedService

	operationsTo := make(map[string]string)
	operationsFrom := make(map[string]string)

	for operationKey, edge := range operation {
		if serviceInServices(edge.To, servicesToBeMerged) && serviceInServices(edge.From, servicesToBeMerged) {
			delete(operation, operationKey)
			continue
		}

		if serviceInServices(edge.To, servicesToBeMerged) {
			if _, ok := operationsTo[edge.From]; !ok {
				operationsTo[edge.From] = operationKey
			} else {
				operationsTo[edge.From] = fmt.Sprintf("%s/%s", operationsTo[edge.From], operationKey)
			}
			delete(operation, operationKey)
		}

		if serviceInServices(edge.From, servicesToBeMerged) {
			if _, ok := operationsFrom[edge.To]; !ok {
				operationsFrom[edge.To] = operationKey
			} else {
				operationsFrom[edge.To] = fmt.Sprintf("%s/%s", operationsFrom[edge.To], operationKey)
			}
			delete(operation, operationKey)
		}
	}

	for to, operationName := range operationsTo {
		operation[operationName] = models.OperationEdge{
			To:    mergedService.Name,
			From:  to,
			Count: 1,
		}
	}

	for from, operationName := range operationsFrom {
		operation[operationName] = models.OperationEdge{
			To:    from,
			From:  mergedService.Name,
			Count: 1,
		}
	}

	annotation.Recomendation.Message = fmt.Sprintf("Merged services into %s", mergedService.Name)
	operations[annotation.InitiatingOperation] = operation

	return services, operations, annotation

}

func serviceInServices(target string, array []models.TelevisorService) bool {
	for _, a := range array {
		if target == a.Name {
			return true
		}
	}

	return false
}
