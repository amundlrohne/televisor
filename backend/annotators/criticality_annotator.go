package annotators

import (
	"fmt"
	"sort"

	"github.com/amundlrohne/televisor/models"
)

func AbsoluteCriticalService(services map[string]models.TelevisorService) []models.Annotation {
	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return len(services[keys[i]].Dependencies)*len(services[keys[i]].Dependents) < len(services[keys[j]].Dependencies)*len(services[keys[j]].Dependents)
	})

    //	criticalService := services[keys[len(keys)-1]]

    annotations := []models.Annotation{}

    for _, k := range keys[len(keys)-4:len(keys)-1] {
        service := services[k]
        annotations = append(annotations, models.Annotation{
            Services: []string{service.Name},
            AnnotationType: models.Criticality,
            YChartLevel: models.ServiceLevel,
            Message: fmt.Sprintf("Service %s has %v dependents and %v dependencies", service.Name, len(service.Dependents), len(service.Dependencies)),
        })
    }

    return annotations
    //    return models.Annotation{
    //    	Services:       []string{criticalService.Name},
    //		AnnotationType: models.Criticality,
    //		YChartLevel:    models.ServiceLevel,
    //		Message:        fmt.Sprintf("Service %s has %v dependents and %v dependencies", criticalService.Name, len(criticalService.Dependents), len(criticalService.Dependencies)),
    //	}
}
