package annotators

import (
	"sort"

	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
	"github.com/jaegertracing/jaeger/model"
)

func AbsoluteCriticalService(sdg []model.DependencyLink) models.Annotation {
	var services = utils.ExtractServicesFromSDG(sdg)

	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return len(services[keys[i]].Dependencies)*len(services[keys[i]].Dependents) < len(services[keys[j]].Dependencies)*len(services[keys[j]].Dependents)
	})

	return models.Annotation{Services: []string{keys[len(keys)-1]}, AnnotationType: models.Criticality}
}

func AbsoluteDependenceService(sdg []model.DependencyLink) models.Annotation {
	var services = utils.ExtractServicesFromSDG(sdg)

	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return len(services[keys[i]].Dependencies) < len(services[keys[j]].Dependencies)
	})

	return models.Annotation{Services: []string{keys[len(keys)-1]}, AnnotationType: models.Criticality}
}
