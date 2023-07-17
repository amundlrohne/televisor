package annotators

import (
	"fmt"

	"github.com/amundlrohne/televisor/models"
)

func UnderUtilizationCPUAnnotator(services map[string]models.TelevisorService, req float64) []models.Annotation {
	annotations := []models.Annotation{}

	for _, s := range services {
		if s.Cpu.P99 <= req {
			annotations = append(annotations, models.Annotation{
				Services:       []string{s.Name},
				AnnotationType: models.UnderUtilizationCPU,
				YChartLevel:    models.ServiceLevel,
				Message:        fmt.Sprintf("Service %s is using %v of the CPU", s.Name, s.Cpu.P99),
			})
		}
	}

	return annotations
}

func UnderUtilizationMemoryAnnotator(services map[string]models.TelevisorService, req float64) []models.Annotation {
	annotations := []models.Annotation{}

	for _, s := range services {
		if s.Memory.P99 <= req {
			annotations = append(annotations, models.Annotation{
				Services:       []string{s.Name},
				AnnotationType: models.UnderUtilizationMemory,
				YChartLevel:    models.ServiceLevel,
				Message:        fmt.Sprintf("Service %s is using %v of the memory", s.Name, s.Memory.P99),
			})
		}
	}

	return annotations
}
