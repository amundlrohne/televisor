package annotators

import (
	"github.com/amundlrohne/televisor/models"
)

func GreedyServiceAnnotator(operations models.Operations) []models.Annotation {
	annotations := []models.Annotation{}

	nonRefOps := operations.ClearReflexiveEdges()

	// Check all operations
	for _, v := range nonRefOps {
		// For every edge check if there is another edge that has the same To but different From
		// Since all edges are from the same request, if they converge they are greedy.
		for edgeKey, edgeValue := range v {
			//var convergingEdges []models.OperationEdge
			convergingEdges := []models.OperationEdge{edgeValue} //append(convergingEdges, edgeValue)
			convergingOperations := []string{edgeKey}
			for compKey, compValue := range v {
				if compValue.From != edgeValue.From && compValue.To == edgeValue.To {
					convergingEdges = append(convergingEdges, compValue)
					convergingOperations = append(convergingOperations, compKey)
				}
			}
			// Create annotation if service has two or more converging edges
			if len(convergingEdges) > 1 {
				annotation := models.Annotation{}
				annotation.AnnotationType = models.Greedy
				annotation.YChartLevel = models.OperationLevel
				annotation.Operations = convergingOperations
				for _, edge := range convergingEdges {
					annotation.Services = append(annotation.Services, edge.From)
				}
				annotations = append(annotations, annotation)
			}

		}
	}

	return annotations
}
