package annotators

import (
	"github.com/amundlrohne/televisor/models"
)

func GreedyServiceAnnotator(operations models.Operations) []models.Annotation {
	annotations := []models.Annotation{}

	nonRefOps := operations.ClearReflexiveEdges()

	// Check all operations
	for _, v := range nonRefOps {
		compEdges := v
		// For every edge check if there is another edge that has the same To but different From
		// Since all edges are from the same request, if they converge they are greedy.
		for edgeKey, edgeValue := range v {
			convergingEdges := []models.OperationEdge{edgeValue} //append(convergingEdges, edgeValue)
			convergingOperations := []string{edgeKey}
			delete(compEdges, edgeKey)
			for compKey, compValue := range compEdges {
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

	return removeSubsets(annotations)
}

func removeSubsets(annotations []models.Annotation) []models.Annotation {
	result := []models.Annotation{}

	for _, a := range annotations {

		isUnique := true
		for j, r := range result {
			if isSubset(a.Services, r.Services) && isSubset(a.Operations, r.Operations) {
				result[j] = a
				isUnique = false
			} else if isSubset(r.Services, a.Services) && isSubset(r.Operations, a.Operations) {
				isUnique = false
				break
			}
		}
		if isUnique {
			result = append(result, a)
		}
	}

	return result
}

func isSubset(set []string, subset []string) bool {
	if len(subset) > len(set) {
		return false
	}

	setmap := make(map[string]int)

	for _, s := range set {
		setmap[s] = 1
	}

	for _, s := range subset {
		if count, found := setmap[s]; !found {
			return false
		} else if count > 1 {
			return false
		} else {
			setmap[s] = count - 1
		}
	}

	return true
}

func minOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
