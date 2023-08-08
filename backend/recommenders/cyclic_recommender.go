package recommenders

import (
	"github.com/amundlrohne/televisor/models"
)

func CyclicRecommender(services map[string]models.TelevisorService, operations models.Operations, annotation models.Annotation) (map[string]models.TelevisorService, models.Operations, models.Annotation) {
	operation := operations[annotation.InitiatingOperation]
	// operationServices := servicesInOperation(services, operation)

	for _, edgeKey := range annotation.Operations {
		tempOperation := cloneOperation(operation)
		delete(tempOperation, edgeKey)
		isolatedEdge := false
		for _, s := range annotation.Services {
			if !serviceHasToEdgeInOperation(s, tempOperation) {
				isolatedEdge = true
				break
			}
		}

		if !isolatedEdge {
			delete(operation, edgeKey)
			break
		}
	}

	operations[annotation.InitiatingOperation] = operation

	return services, operations, annotation
}

func cloneOperation(operation models.Operation) models.Operation {
	result := make(models.Operation)

	for k, v := range operation {
		result[k] = v
	}

	return result
}

func servicesInOperation(services map[string]models.TelevisorService, operation models.Operation) map[string]models.TelevisorService {
	result := make(map[string]models.TelevisorService)

	for _, o := range operation {
		result[o.From] = services[o.From]
		result[o.To] = services[o.To]
	}

	return result
}

func serviceHasToEdgeInOperation(serviceName string, operation models.Operation) bool {
	for _, o := range operation {
		if o.To == serviceName {
			return true
		}
	}

	return false
}
