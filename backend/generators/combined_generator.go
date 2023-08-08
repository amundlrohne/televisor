package generators

import "github.com/amundlrohne/televisor/models"

func OperationsGenerator() models.Operations {
	operations := make(models.Operations)

	operations["operation-inappropriate-intimacy"] = make(models.Operation)
	operations["operation-megaservice"] = make(models.Operation)
	operations["operation-cyclic"] = make(models.Operation)
	operations["operation-greedy"] = make(models.Operation)

	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop1", "api-gateway", "service-a")
	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop2", "service-a", "service-b")
	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop3", "service-a", "service-c")
	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop4", "service-b", "service-e")
	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop5", "service-c", "service-d")
	operations["operation-inappropriate-intimacy"].AddEdge("op1-subop6", "service-d", "service-e")

	operations["operation-megaservice"].AddEdge("op2-subop1", "api-gateway", "service-f")
	operations["operation-megaservice"].AddEdge("op2-subop2", "api-gateway", "service-f")
	operations["operation-megaservice"].AddEdge("op2-subop3", "service-f", "service-g")
	operations["operation-megaservice"].AddEdge("op2-subop4", "service-f", "service-h")

	operations["operation-cyclic"].AddEdge("op3-subop1", "api-gateway", "service-i")
	operations["operation-cyclic"].AddEdge("op3-subop2", "service-i", "service-j")
	operations["operation-cyclic"].AddEdge("op3-subop3", "service-j", "service-k")
	operations["operation-cyclic"].AddEdge("op3-subop4", "service-k", "service-i")

	operations["operation-greedy"].AddEdge("op4-subop1", "api-gateway", "service-l")
	operations["operation-greedy"].AddEdge("op4-subop2", "service-l", "service-m")

	return operations
}

func ServiceUtilizationGenerator(services map[string]models.TelevisorService) map[string]models.TelevisorService {
	result := services

	if r, ok := result["api-gateway"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.54, Mean: 0.23, Stdev: 0.0022}
		r.Memory = models.Utilization{Quantile: 0.27, Mean: 0.23, Stdev: 0.0000}
		result["api-gateway"] = r
	}

	if r, ok := result["service-a"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.45, Mean: 0.23, Stdev: 0.0024}
		r.Memory = models.Utilization{Quantile: 0.19, Mean: 0.16, Stdev: 0.0000}
		result["service-a"] = r
	}

	if r, ok := result["service-b"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.21, Mean: 0.15, Stdev: 0.0013}
		r.Memory = models.Utilization{Quantile: 0.10, Mean: 0.07, Stdev: 0.0003}
		result["service-b"] = r
	}

	if r, ok := result["service-c"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.18, Mean: 0.14, Stdev: 0.0022}
		r.Memory = models.Utilization{Quantile: 0.32, Mean: 0.23, Stdev: 0.0009}
		result["service-c"] = r
	}

	if r, ok := result["service-d"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.14, Mean: 0.13, Stdev: 0.0022}
		r.Memory = models.Utilization{Quantile: 0.07, Mean: 0.04, Stdev: 0.0020}
		result["service-d"] = r
	}

	if r, ok := result["service-e"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.24, Mean: 0.20, Stdev: 0.0005}
		r.Memory = models.Utilization{Quantile: 0.15, Mean: 0.10, Stdev: 0.0001}
		result["service-e"] = r
	}

	if r, ok := result["service-f"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.55, Mean: 0.37, Stdev: 0.0017}
		r.Memory = models.Utilization{Quantile: 0.38, Mean: 0.32, Stdev: 0.0005}
		result["service-f"] = r
	}

	if r, ok := result["service-g"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.25, Mean: 0.17, Stdev: 0.0003}
		r.Memory = models.Utilization{Quantile: 0.04, Mean: 0.03, Stdev: 0.0000}
		result["service-g"] = r
	}

	if r, ok := result["service-h"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.14, Mean: 0.12, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.33, Mean: 0.28, Stdev: 0.0002}
		result["service-h"] = r
	}

	if r, ok := result["service-i"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.08, Mean: 0.05, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.21, Mean: 0.17, Stdev: 0.0002}
		result["service-i"] = r
	}

	if r, ok := result["service-j"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.09, Mean: 0.06, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.10, Mean: 0.3, Stdev: 0.0002}
		result["service-j"] = r
	}

	if r, ok := result["service-k"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.09, Mean: 0.06, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.10, Mean: 0.3, Stdev: 0.0002}
		result["service-k"] = r
	}

	if r, ok := result["service-l"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.11, Mean: 0.07, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.09, Mean: 0.06, Stdev: 0.0002}
		result["service-l"] = r
	}

	if r, ok := result["service-m"]; ok {
		r.Cpu = models.Utilization{Quantile: 0.02, Mean: 0.01, Stdev: 0.0007}
		r.Memory = models.Utilization{Quantile: 0.03, Mean: 0.01, Stdev: 0.0002}
		result["service-m"] = r
	}

	return result
}
