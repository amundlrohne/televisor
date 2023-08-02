package recommenders

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
)

func InappropriateIntimacyRecommender(services map[string]models.TelevisorService, operations models.Operations, operation string, serviceNames []string) (map[string]models.TelevisorService, models.Operations) {
	newService := services[serviceNames[0]]
	delete(services, serviceNames[0])

	for _, n := range serviceNames[1:] {
		currentService := services[n]
		newService = models.TelevisorService{
			Name:    fmt.Sprintf("%s/%s", newService.Name, currentService.Name),
			Cpu:     utils.AddUtilizations(newService.Cpu, currentService.Cpu),
			Memory:  utils.AddUtilizations(newService.Memory, currentService.Memory),
			Network: utils.AddUtilizations(newService.Network, currentService.Network),
		}
		delete(services, n)
	}

	services[newService.Name] = newService

	for kos, os := range operations {
		for ko, o := range os {
			if isMemberOf(o.To, serviceNames) {
				o.To = newService.Name
			}

			if isMemberOf(o.From, serviceNames) {
				o.From = newService.Name
			}

			os[ko] = o
		}
		operations[kos] = os
	}

	mergedEdges := make(map[string]models.OperationEdge)
	mergedOperationNames := make(map[string]string)

	for ko, o := range operations[operation] {

		key := fmt.Sprintf("%s-%s", o.From, o.To)
		if m, ok := mergedEdges[key]; !ok {
			mergedEdges[key] = o
			mergedOperationNames[key] = ko
		} else {
			m.Count += 1
			mergedEdges[key] = m
			mergedOperationNames[key] = fmt.Sprintf("%s/%s", mergedOperationNames[key], ko)
		}
		delete(operations[operation], ko)

	}

	for k, o := range mergedEdges {
		operations[operation][mergedOperationNames[k]] = o
	}

	return services, operations
}

func isMemberOf(member string, list []string) bool {
	for _, l := range list {
		if member == l {
			return true
		}
	}

	return false
}
