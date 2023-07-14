package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

func CyclicDependencyAnnotator(operations models.Operations, services map[string]models.ExtendedService) []models.Annotation {
	cycles := []models.ExtendedService{}
	ops := operations.ClearReflexiveEdges()

	for k, service := range services {
		for _, op := range ops {
			if op.IsConnected(k, k) {
				cycles = append(cycles, service)
			}
		}
	}

	annotations := []models.Annotation{}

	for _, c := range cycles {
		annotations = append(annotations, models.Annotation{
			Services:       []string{c.Name},
			AnnotationType: models.Cyclic,
			YChartLevel:    models.OperationLevel,
			Message:        fmt.Sprintf("Service %s has a cyclic relationship", c.Name),
		})
	}

	return annotations

}
