package annotators

import (
	"github.com/amundlrohne/televisor/models"
)

func MegaserviceAnnotator(operations models.Operations) []models.Annotation {
	annotations := []models.Annotation{}

	// Loop over all operations
	for _, v := range operations {
		// If operation has two or more edges that share the same From and To, To is a megaservice
		for edgeKey, edge := range v {
			// If operation is reflexive, it is not a megaservice
			if edge.To == edge.From {
				continue
			}
			megaservice := []models.OperationEdge{edge}
			operations := []string{edgeKey}
			for compKey, comp := range v {
				// If operation is reflexive, it is not a megaservice
				if comp.To == comp.From {
					continue
				}
				if comp.From == edge.From && comp.To == edge.To && edgeKey != compKey {
					megaservice = append(megaservice, comp)
					operations = append(operations, compKey)
				}
			}

			if len(megaservice) > 1 {
				annotation := models.Annotation{
					AnnotationType: models.Megaservice,
					YChartLevel:    models.OperationLevel,
					Services:       []string{megaservice[0].To},
					Operations:     operations,
				}
				annotations = append(annotations, annotation)
			}
		}
	}

	return annotations
}
