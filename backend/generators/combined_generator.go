package generators

import "github.com/amundlrohne/televisor/models"

func CombinedGenerator() models.Operations {
	operations := make(models.Operations)

	operations["operation1"] = make(models.Operation)
	operations["operation2"] = make(models.Operation)
	operations["operation3"] = make(models.Operation)

	operations["operation1"].AddEdge("op1-subop1", "api-gateway", "service-a")
	operations["operation1"].AddEdge("op1-subop2", "api-gateway", "service-b")
	operations["operation1"].AddEdge("op1-subop3", "api-gateway", "service-c")
	operations["operation1"].AddEdge("op1-subop4", "api-gateway", "service-d")
	operations["operation1"].AddEdge("op1-subop5", "api-gateway", "service-d")

	operations["operation2"].AddEdge("op2-subop1", "api-gateway", "service-a")
	operations["operation2"].AddEdge("op2-subop5", "service-a", "service-e")
	operations["operation2"].AddEdge("op2-subop6", "service-e", "service-a")
	operations["operation2"].AddEdge("op2-subop2", "api-gateway", "service-b")
	operations["operation2"].AddEdge("op2-subop3", "api-gateway", "service-c")
	operations["operation2"].AddEdge("op2-subop4", "api-gateway", "service-d")

	operations["operation3"].AddEdge("op3-subop1", "api-gateway", "service-a")
	operations["operation3"].AddEdge("op3-subop2", "api-gateway", "service-b")
	operations["operation3"].AddEdge("op3-subop3", "service-b", "service-f")
	operations["operation3"].AddEdge("op3-subop4", "service-b", "service-g")
	operations["operation3"].AddEdge("op3-subop5", "service-b", "service-h")

	return operations
}
