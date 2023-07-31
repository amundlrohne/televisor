package recommenders

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
)

func MegaserviceRecommender(service models.TelevisorService, operation models.Operation) ([]models.TelevisorService, models.Operation) {
	ins := models.Operation{}
	outs := models.Operation{}

	for op, edge := range operation {
		if edge.To == service.Name {
			ins.AddEdge(op, edge.From, edge.To)
		}

		if edge.From == service.Name {
			outs.AddEdge(op, edge.From, edge.To)
		}
	}

	splitServices := []models.TelevisorService{}

	isFirst := true
	for k, v := range ins {
		serviceName := fmt.Sprintf("%s/%s", service.Name, k)
		if isFirst {
			serviceName = service.Name
			isFirst = false
		}
		splitServices = append(splitServices, models.TelevisorService{
			Name:       serviceName,
			Dependents: []string{v.From},
			Cpu:        utils.SplitUtilization(service.Cpu, len(ins)),
			Memory:     utils.SplitUtilization(service.Memory, len(ins)),
			Network:    utils.SplitUtilization(service.Network, len(ins)),
		})
		v.To = serviceName
		ins[k] = v
	}

	outsKeys := make([]string, 0, len(outs))
	for k := range outs {
		outsKeys = append(outsKeys, k)
	}

	// Round robin assignment of dependencies
	for i := 0; i < len(outs); i++ {
		splitServices[i%len(splitServices)].Dependencies = append(splitServices[i%len(splitServices)].Dependencies, outs[outsKeys[i]].To)
		o := models.OperationEdge{
			From:  splitServices[i%len(splitServices)].Name,
			To:    outs[outsKeys[i]].To,
			Count: outs[outsKeys[i]].Count}
		outs[outsKeys[i]] = o
	}

	for k, v := range ins {
		operation[k] = v
	}

	for k, v := range outs {
		operation[k] = v
	}

	return splitServices, operation
}
