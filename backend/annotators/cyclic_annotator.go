package annotators

import (
	"errors"
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

func CyclicDependencyAnnotator(operations models.Operations, services map[string]models.TelevisorService) []models.Annotation {
	//cycles := []models.TelevisorService{}
	ops := operations.ClearReflexiveEdges()

	cycles := make(map[string][]cycle)

	for k, op := range ops {
		cycle := findCycles(op)
		cycles[k] = cycle
	}

	annotations := []models.Annotation{}

	for k, c := range cycles {
		for _, cc := range c {
			annotations = append(annotations, models.Annotation{
				Services:            cc.Services,
				Operations:          cc.Operations,
				AnnotationType:      models.Cyclic,
				YChartLevel:         models.OperationLevel,
				InitiatingOperation: k,
				Message:             fmt.Sprintf("Services %v has a cyclic relationship", cc.Services),
			})
		}
	}

	return annotations

}

func findCycles(operations models.Operation) []cycle {
	unvisitedServices := make(map[string]bool)
	for _, edge := range operations {
		unvisitedServices[edge.From] = true
		unvisitedServices[edge.To] = true
	}

	cycles := []cycle{}

	for k, edge := range operations {
		operationsClone := models.Operation{}
		for k, v := range operations {
			operationsClone[k] = v
		}

		cycle, err := recursiveFindCycle(k, edge.From, operationsClone, cycle{})
		if err != nil {
			continue
		}
		if len(cycle.Services) != 0 {
			cycles = append(cycles, cycle)
			break
		}
	}
	return cycles
}

func recursiveFindCycle(currentOperation string, targetService string, operations models.Operation, history cycle) (cycle, error) {
	if len(operations) == 0 {
		return cycle{}, errors.New("No cycle")
	}

	if operations[currentOperation].To == targetService {
		history.Operations = append(history.Operations, currentOperation)
		history.Services = append(history.Services, operations[currentOperation].To)
		return history, nil
	}

	for k, o := range operations {
		if o.From == operations[currentOperation].To {

			history.Operations = append(history.Operations, currentOperation)
			history.Services = append(history.Services, operations[currentOperation].To)
			delete(operations, currentOperation)
			return recursiveFindCycle(k, targetService, operations, history)
		}
	}

	return cycle{}, errors.New("Dead end")
}

type cycle struct {
	Services   []string
	Operations []string
}
