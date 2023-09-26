package utils

import (
	"errors"
	"fmt"
	"strconv"

	pb "jaeger-idl/api_v2"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/queries"
	"github.com/jaegertracing/jaeger/model"
)

func ExtractServicesFromSDG(sdg []models.OperationEdge) map[string]models.TelevisorService {
	var services = make(map[string]models.TelevisorService)

	for _, value := range sdg {
		if serviceFrom, ok := services[value.From]; ok {
			serviceFrom.Dependents = append(serviceFrom.Dependents, value.To)
			services[value.From] = serviceFrom
		} else {
			services[value.From] = models.TelevisorService{
				Name:       value.From,
				Dependents: []string{value.To},
			}
		}

		if serviceTo, ok := services[value.To]; ok {
			serviceTo.Dependencies = append(serviceTo.Dependencies, value.From)
			services[value.To] = serviceTo
		} else {
			services[value.To] = models.TelevisorService{
				Name:         value.To,
				Dependencies: []string{value.From},
			}
		}
	}

	return services
}

func GetSubSDGs(qsc pb.QueryServiceClient, service string) models.Operations {
	operations := queries.JaegerOperations(qsc, service)

    //fmt.Println(operations)
	result := make(models.Operations)
	spanIDToService := make(map[string]string)

	for _, o := range operations {
		traces := queries.JaegerTraces(qsc, service, o)
        //fmt.Printf("Number of traces for operation %v is %v \n", o, len(traces))
        //fmt.Println(traces)
		root, err := getRootSpan(traces)
		if err != nil {
			fmt.Printf("could not find root span for operation: %v \n", o)
			continue
		}

		// Assign reflexive property to root span
		if _, ok := result[root.OperationName]; !ok {
			result[root.OperationName] = make(map[string]models.OperationEdge)

			result[root.OperationName][root.SpanID.String()] = models.OperationEdge{
                Name:  root.OperationName,
				From:  root.Process.ServiceName,
				To:    root.Process.ServiceName,
				Count: 1,
			}
			spanIDToService[root.SpanID.String()] = root.Process.ServiceName
		}

		// Retrieve all correlations between spanID and service name
		for _, t := range traces {
			spanIDToService[t.SpanID.String()] = t.Process.ServiceName
		}

		// Add result where key is "sub-operation" name and value is relationship
		for _, t := range traces {
			reference, err := getSpanChildRef(t)
			if err != nil {
				//fmt.Printf("could not find child_of service for operation: %v \n", t.OperationName)
				//fmt.Printf("error: %v \n", err)
				continue
			}

			if path, ok := result[root.OperationName][t.SpanID.String()]; ok {
				path.Count++
				result[root.OperationName][t.SpanID.String()] = path
			} else {
				fromService := spanIDToService[reference]

				result[root.OperationName][t.SpanID.String()] = models.OperationEdge{
                    Name:  t.OperationName,
					From:  fromService,
					To:    t.Process.ServiceName,
					Count: 1,
				}
			}
		}
		//fmt.Printf("Operation %s = %v \n\n", o, spanIDToService)
		spanIDToService = make(map[string]string)
	}

	return result
}

func getSpanChildRef(span model.Span) (string, error) {
	references := span.References

	if len(references) == 0 {
		return "", errors.New("no references")
	}

	for _, r := range references {
		if r.RefType == model.SpanRefType_CHILD_OF {
			return r.SpanID.String(), nil
		}
	}

	return "", errors.New("no CHILD_OF reference")
}

func getRootSpan(traces []model.Span) (model.Span, error) {
	for _, t := range traces {
		if t.SpanID.String() == t.TraceID.String() {
			return t, nil
		}
	}

	return model.Span{}, errors.New("can't find root span")
}

func AddCPUQuantileToServices(services map[string]models.TelevisorService, cpuData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range cpuData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Cpu.Quantile = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}

func AddCPUMeanToServices(services map[string]models.TelevisorService, cpuData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range cpuData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Cpu.Mean = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}

func AddCPUStdevToServices(services map[string]models.TelevisorService, cpuData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range cpuData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Cpu.Stdev = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}

func AddMemoryQuantileToServices(services map[string]models.TelevisorService, memoryData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range memoryData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Memory.Quantile = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}

func AddMemoryMeanToServices(services map[string]models.TelevisorService, memoryData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range memoryData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Memory.Mean = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}

func AddMemoryStdevToServices(services map[string]models.TelevisorService, memoryData models.PrometheusContainerMetric) map[string]models.TelevisorService {
	for serviceKey, s := range services {
		for metricKey, m := range memoryData {
			if s.IsServiceInContainer(metricKey) {
				utilization, _ := strconv.ParseFloat(m, 64)
				s.Memory.Stdev = utilization
				services[serviceKey] = s
			}
		}
	}
	return services
}
