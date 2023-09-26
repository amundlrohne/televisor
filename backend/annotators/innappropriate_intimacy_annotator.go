package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

func findConvergingPaths(c convergance, operation models.Operation, result ...models.OperationEdge) []models.OperationEdge {
	//result := []models.OperationEdge{}

	for _, conv := range c.Converging {

		count := 0
		for _, op := range operation {
			if conv == op.From {
				count++
			}
		}

		if count > 1 {
			continue
		}

		for ko, op := range operation {
			if c.Target == op.To && conv == op.From {

				for k, o := range operation {
					if o.To == conv {
						o.To = c.Target
						operation[k] = o
					}
				}
				result = append(result, op)
				delete(operation, ko)
			}
		}
	}

	convergingService := findConvergingService(operation)

	if len(convergingService) > 0 {
		return findConvergingPaths(convergingService[0], operation, result...)
	}

	return result
}

//func testIIRequests(requests models.Operations, cyclicOperations []models.Annotation) []models.Annotation {
//  	annotations := []models.Annotation{}
//	requestsClone:= requests.ClearReflexiveEdges()
//
//	for _, c := range cyclicOperations {
//		delete(nonReflexiveOperations, c.InitiatingOperation)
//	}
//
//    for _, request := range requestsClone {
//        annotations = append(annotations, testMethod(request))
//    }
//
//    return annotations
//}
//
//func testMethod(request models.Operation, cyclicOperations []models.Annotation) []models.Annotation {
//    annotations := []models.Annotation{}
//    divergingServices := divergingServices(request)
//    converingServices := convergingServices(request)
//
//
//
//    return annotations
//}
//
//func convergingServices(request models.Operations) []string {
//    serviceInboundEdges := make(map[string]int)
//    for key, edge := range request {
//        if _, ok := serviceInboundEdges[edge.To]; !ok {
//           serviceInboundEdges[edge.To] = 1
//        } else {
//            serviceInboundEdges[edge.To] += 1
//        }
//    }
//
//    result := []string{}
//
//    for key, q := range serviceInboundEdges {
//        if q > 1 {
//            result = append(result, key)
//        }
//    }
//
//    return result
//}
//
//
//func divergingServices(request models.Operations) []string {
//    serviceOutboundEdges := make(map[string]int)
//    for key, edge := range request {
//        if _, ok := serviceOutboundEdges[edge.To]; !ok {
//           serviceOutboundEdges[edge.From] = 1
//        } else {
//            serviceOutboundEdges[edge.From] += 1
//        }
//    }
//
//    result := []string{}
//
//    for key, q := range serviceOutboundEdges {
//        if q > 1 {
//            result = append(result, key)
//        }
//    }
//
//    return result
//}
//
func InappropriateIntimacyServiceAnnotator(requests models.Operations, cyclicOperations []models.Annotation) []models.Annotation {
	annotations := []models.Annotation{}
	nonReflexiveOperations := requests.ClearReflexiveEdges()

	for _, c := range cyclicOperations {
		delete(nonReflexiveOperations, c.InitiatingOperation)
	}

	convergingServices := findConvergingServices(nonReflexiveOperations)

	for req, conv := range convergingServices {
		for _, c := range conv {
			services := []string{}
			for _, e := range findConvergingPaths(c, nonReflexiveOperations[req]) {
				services = append(services, e.From)
			}

			annotations = append(annotations, models.Annotation{
				Services:            services,
				InitiatingOperation: req,
				AnnotationType:      models.InappropriateIntimacy,
				YChartLevel:         models.OperationLevel,
				AnnotationLevel:     models.Critical,
				Message:             fmt.Sprintf("Services: %v should be merged.", services),
			})
		}
	}

	return annotations
}

func findConvergingServices(reqs models.Operations) map[string][]convergance {
	targetServices := make(map[string][]convergance)
	for req, operations := range reqs {
		x := findConvergingService(operations)
		if len(x) > 0 {

			targetServices[req] = x
		}
	}

	return targetServices
}

func findConvergingService(operations models.Operation) []convergance {
	targetServices := []convergance{}

	incomingEdges := make(map[string]convergance)
	for _, edge := range operations {
		if ie, ok := incomingEdges[edge.To]; !ok {
			incomingEdges[edge.To] = convergance{Target: edge.To, Converging: []string{edge.From}}
		} else {
			if existsInStringArray(edge.From, ie.Converging) {
				continue
			}
			ie.Converging = append(incomingEdges[edge.To].Converging, edge.From)
			incomingEdges[edge.To] = ie
		}
	}

	for service, ies := range incomingEdges {
		if len(ies.Converging) < 2 {
			delete(incomingEdges, service)
		} else {
			targetServices = append(targetServices, ies)
		}
	}

	return targetServices
}

type convergance struct {
	Target     string
	Converging []string
}

func existsInStringArray(target string, array []string) bool {
	for _, a := range array {
		if target == a {
			return true
		}
	}
	return false
}
