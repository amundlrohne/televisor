package annotators

import (
	"fmt"
	"sort"

	"github.com/amundlrohne/televisor/models"
)

func AbsoluteDependenceService(services map[string]models.TelevisorService) []models.Annotation {
	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return len(services[keys[i]].Dependencies) < len(services[keys[j]].Dependencies)
	})

    //	dependenceService := services[keys[len(keys)-1]]

    annotations :=  []models.Annotation{}

    for _, k := range keys[len(keys)-4:len(keys)-1] {
        service := services[k]
        annotations = append(annotations, models.Annotation{
            Services:       []string{service.Name},
            AnnotationType: models.Dependence,
            YChartLevel:    models.ServiceLevel,
            Message:        fmt.Sprintf("Service %s has %v dependencies", service.Name, len(service.Dependencies)),
        })
    }

    return annotations

	//return models.Annotation{
	//	Services:       []string{dependenceService.Name},
	//	AnnotationType: models.Dependence,
	//	YChartLevel:    models.ServiceLevel,
	//	Message:        fmt.Sprintf("Service %s has %v dependencies", dependenceService.Name, len(dependenceService.Dependencies)),
	//}
}
