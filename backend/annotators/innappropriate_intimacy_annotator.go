package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)


func InappropriateIntimacyServiceAnnotatorV2(request models.Operation) []string {
    var divergingNodes = make(map[string]int)
    var convergingNodes = make(map[string]int)

    for _, r := range request {
        if i, ok := divergingNodes[r.From]; ok {
            divergingNodes[r.From] = i + 1
        } else {
            divergingNodes[r.From] = 1
        }

        if i, ok := convergingNodes[r.To]; ok {
            convergingNodes[r.To] = i + 1
        } else {
            convergingNodes[r.To] = 1
        }
    }

    for d, di := range divergingNodes {
        if di < 2 {
            delete(divergingNodes, d)
        }
    }

    for c, ci := range convergingNodes {
        if ci < 2 {
            delete(convergingNodes, c)
        }
    }

    for d := range divergingNodes {
        for c := range convergingNodes {
            if d == c {
                delete(divergingNodes, d)
                delete(convergingNodes, c)
            }
        }
    }

    var dag []models.OperationEdge

    for _, r := range request {
        dag = append(dag, r)
    }

    if len(divergingNodes) > 0 && len(convergingNodes) > 0 {
        var nodeSets [][]string
        for d := range divergingNodes {
            for c := range convergingNodes {
               nodeSets = append(nodeSets, servicesBetween(d, c, dag))
            }
        }

        //log.Printf("sets: %v", nodeSets)

        var selected int
        largest := 0
        for i, ns := range nodeSets {
            if len(ns) > largest {
                largest = len(ns)
                selected = i
            }
        }

        return nodeSets[selected]
    }

    return []string{}
}

func servicesBetween(from string, to string, dag []models.OperationEdge) []string {
    var edges = make(map[string][]string)
    var counts = make(map[string]int)

    for _, e := range dag {
        edges[e.From] = append(edges[e.From], e.To)

        if _, ok := counts[e.From]; !ok {
            counts[e.From] = 0
        }

        if _, ok := counts[e.To]; !ok {
            counts[e.To] = 1
        } else {
            counts[e.To] += 1
        }
    }

    sorted := topologicalSort(edges, counts)
    nPathsBetween := numberOfPaths(sorted, from, to, edges, counts)

    if nPathsBetween > 1 {
        //log.Printf("sorted: %v", sorted)

        for i := 0; i < len(sorted); i++ {
            if sorted[i] == from {
                sorted = sorted[i:]
                break
            }
        }

        for i := len(sorted) - 1; i >= 0; i-- {
            if sorted[i] == to {
                sorted = sorted[:i+1]
                break
            }
        }

        return sorted
    }

    return []string{}
}

func topologicalSort(edges map[string][]string, counts map[string]int) []string {
    var queue []string

    for e, c := range counts {
        if c == 0 {
            queue = append(queue, e)
        }
    }

    var result []string

    for len(queue) > 0 {
        u := queue[0]
        queue = queue[1:]

        result = append(result, u)

        for _, c := range edges[u] {
            counts[c] -= 1

            if counts[c] == 0 {
                queue = append(queue, c)
            }
        }

    }

    return result
}

func numberOfPaths(sorted []string, source string, destination string, edges map[string][]string, counts map[string]int) int {
    result := make(map[string]int)

    result[destination] = 1

    for i := len(sorted)-1; i >= 0; i-- {
        for _, e := range edges[sorted[i]] {
           result[sorted[i]] += result[e]
        }
    }

    return result[source]

}

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

func InappropriateIntimacyServiceAnnotator(requests models.Operations, cyclicOperations []models.Annotation) []models.Annotation {
	annotations := []models.Annotation{}
	nonReflexiveOperations := requests.ClearReflexiveEdges()

	for _, c := range cyclicOperations {
		delete(nonReflexiveOperations, c.InitiatingOperation)
	}


    for req, edges := range nonReflexiveOperations {
        responsibleServices := InappropriateIntimacyServiceAnnotatorV2(edges)

        if len(responsibleServices) > 0 {
            annotations = append(annotations, models.Annotation{
                Services: responsibleServices,
                InitiatingOperation: req,
                AnnotationType: models.InappropriateIntimacy,
                YChartLevel: models.OperationLevel,
                AnnotationLevel: models.Critical,
                Message: fmt.Sprintf("Services: %v should be merged.", responsibleServices),
            })
        }

    }

    return annotations

//	convergingServices := findConvergingServices(nonReflexiveOperations)
//
//	for req, conv := range convergingServices {
//		for _, c := range conv {
//			services := []string{}
//			for _, e := range findConvergingPaths(c, nonReflexiveOperations[req]) {
//				services = append(services, e.From)
//			}
//
//			annotations = append(annotations, models.Annotation{
//				Services:            services,
//				InitiatingOperation: req,
//				AnnotationType:      models.InappropriateIntimacy,
//				YChartLevel:         models.OperationLevel,
//				AnnotationLevel:     models.Critical,
//				Message:             fmt.Sprintf("Services: %v should be merged.", services),
//			})
//		}
//	}
//
//	return annotations
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
