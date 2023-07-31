package annotators

import (
	"github.com/amundlrohne/televisor/models"
)

func MegaserviceAnnotator(operations models.Operations) []models.Annotation {
	annotations := []models.Annotation{}

	operations = operations.ClearReflexiveEdges()

	// Loop over all operations
	for apiK, v := range operations {
		// If operation has two or more edges that share the same From and To, To is a megaservice
		for edgeKey, edge := range v {
			megaservice := []models.OperationEdge{edge}
			operations := []string{edgeKey}
			for compKey, comp := range v {
				if comp.From == edge.From && comp.To == edge.To && edgeKey != compKey {
					megaservice = append(megaservice, comp)
					operations = append(operations, compKey)
				}
			}

			if len(megaservice) > 1 {
				exists := false
				for _, a := range annotations {
					if a.Services[0] == megaservice[0].To {
						exists = true
					}
				}

				if !exists {
					annotation := models.Annotation{
						AnnotationType:      models.Megaservice,
						YChartLevel:         models.OperationLevel,
						Services:            []string{megaservice[0].To},
						Operations:          operations,
						InitiatingOperation: apiK,
					}
					annotations = append(annotations, annotation)
				}
			}
		}
	}

	return annotations
}
