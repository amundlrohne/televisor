package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

// Check services that have a singular function
// Microservices with very limited functionalities (e.g., a microservice serving only one static HTML page).

func GreedyServiceAnnotator(operations models.Operations, services map[string]models.TelevisorService) []models.Annotation {
	annotations := []models.Annotation{}

	serviceDegrees := make(map[string][]greedyMicroservice)
    operations = operations.ClearReflexiveEdges()

	for rootKey, rootOperations := range operations {
		for opK, op := range rootOperations {
			if service, ok := serviceDegrees[op.To]; !ok {
				serviceDegrees[op.To] = []greedyMicroservice{{Request: rootKey, Operation: opK, ServiceName: op.To, ParentServiceName: op.From}}
			} else {
				service = append(service, greedyMicroservice{Request: rootKey, Operation: opK, ServiceName: op.To, ParentServiceName: op.From})
				serviceDegrees[op.To] = service
			}

			if service, ok := serviceDegrees[op.From]; !ok {
				serviceDegrees[op.From] = []greedyMicroservice{{Request: rootKey, Operation: opK, ServiceName: op.From}}
			} else {
				service = append(service, greedyMicroservice{Request: rootKey, Operation: opK, ServiceName: op.From})
				serviceDegrees[op.From] = service
			}
		}
	}

    parentServices := make(map[string][]greedyMicroservice)

    for _, data := range serviceDegrees {
        if len(data) == 1 {
            if children, ok := parentServices[data[0].ParentServiceName]; !ok {
                parentServices[data[0].ParentServiceName] = data
            } else {
                children = append(children, data[0])
                parentServices[data[0].ParentServiceName] = children
            }
        }
    }



    // Move utils check to recommendation engine
	for service, data := range parentServices {
        servicesToBeMerged := getServiceNames(data)
        servicesToBeMerged = append(servicesToBeMerged, service)
        operationsToBeRemoved := getServiceOperations(data)
        annotations = append(annotations, models.Annotation{
            Services:            servicesToBeMerged,
            Operations:          operationsToBeRemoved,
            InitiatingOperation: data[0].Request,
            AnnotationType:      models.Greedy,
            Message:             fmt.Sprintf("Greedy services: %s", servicesToBeMerged),
            AnnotationLevel:     models.Info,
        })
	}

	return annotations
}

type greedyMicroservice struct {
    ServiceName string
    ParentServiceName string
    Request     string
    Operation   string
}

func getServiceNames(greedy []greedyMicroservice) []string {
    result := []string{}
    for _, data := range greedy {
       result = append(result, data.ServiceName)
    }

    return result
}

func getServiceOperations(greedy []greedyMicroservice) []string {
    result := []string{}
    for _, data := range greedy {
        result = append(result, data.Operation)
    }

    return result
}
