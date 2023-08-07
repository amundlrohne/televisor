package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

// Check services that have a singular function
// Microservices with very limited functionalities (e.g., a microservice serving only one static HTML page).

func GreedyServiceAnnotator(operations models.Operations, services map[string]models.TelevisorService) []models.Annotation {
	annotations := []models.Annotation{}

	serviceDegrees := make(map[string][]requestOperationTuple)

	for rootKey, rootOperations := range operations {
		for opK, op := range rootOperations {
			if service, ok := serviceDegrees[op.To]; !ok {
				serviceDegrees[op.To] = []requestOperationTuple{{Request: rootKey, Operation: opK}}
			} else {
				service = append(service, requestOperationTuple{Request: rootKey, Operation: opK})
				serviceDegrees[op.To] = service
			}

			if service, ok := serviceDegrees[op.From]; !ok {
				serviceDegrees[op.From] = []requestOperationTuple{{Request: rootKey, Operation: opK}}
			} else {
				service = append(service, requestOperationTuple{Request: rootKey, Operation: opK})
				serviceDegrees[op.From] = service
			}
		}
	}

	for service, operations := range serviceDegrees {
		if len(operations) == 1 && serviceUtilsIsAcceptable(services[service]) {
			annotations = append(annotations, models.Annotation{
				Services:            []string{service},
				Operations:          []string{operations[0].Operation},
				InitiatingOperation: operations[0].Request,
				AnnotationType:      models.Greedy,
				Message:             fmt.Sprintf("Service %s, only has a single operation (%s). Potential Greedy service as utils are under requirements.", service, operations),
				AnnotationLevel:     models.Info,
			})
		}
	}

	return annotations
}

func serviceUtilsIsAcceptable(service models.TelevisorService) bool {
	if service.Cpu.Quantile >= 0.15 {
		return false
	}

	if service.Memory.Quantile >= 0.2 {
		return false
	}

	return true
}

type requestOperationTuple struct {
	Request   string
	Operation string
}
