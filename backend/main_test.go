package main

import (
	"encoding/json"
	"os"
    "testing"

	"github.com/amundlrohne/televisor/annotators"
	"github.com/amundlrohne/televisor/generators"
	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/recommenders"
	"github.com/amundlrohne/televisor/utils"
)

var annotations []models.Annotation
var operations models.Operations
var services map[string]models.TelevisorService

func init() {
	annotations = []models.Annotation{}
	operations = generators.OperationsGenerator()
	combinedEdges := operations.CombineEdges()
	services = utils.ExtractServicesFromSDG(combinedEdges)
	services = generators.ServiceUtilizationGenerator(services)
}

func TestAnalyzeMegaserviceAnnotator(t *testing.T) {
	megaservice := annotators.MegaserviceAnnotator(operations)

	if len(megaservice) != 1 {
		t.Fatalf("len(MegaserviceAnnotator()) = %v, want 1, error", len(megaservice))
	}

	expectedService := "service-f"

	if megaservice[0].Services[0] != expectedService {
		t.Fatalf("MegaserviceAnnotator().Services[0] = %s, want %s, error", megaservice[0].Services[0], expectedService)
	}

	annotations = append(annotations, megaservice...)
}

func TestCyclicServiceAnnotator(t *testing.T) {
	cyclic := annotators.CyclicDependencyAnnotator(operations, services)

	annotations = append(annotations, cyclic...)
}

func TestAnalyzeInappropriateIntimacyServiceAnnotator(t *testing.T) {
	cycles := []models.Annotation{}
	for _, a := range annotations {
		if a.AnnotationType == models.Cyclic {
			cycles = append(cycles, a)
		}
	}

	inappropriateIntimacy := annotators.InappropriateIntimacyServiceAnnotator(operations, cycles)

	if len(inappropriateIntimacy) != 1 {
		t.Fatalf("InappropriateIntimacyServiceAnnotator() = %v, want 1, error", len(inappropriateIntimacy))
	}

	expectedServices := map[string]bool{"service-d": false, "service-b": false, "service-c": false}

	for _, g := range inappropriateIntimacy {
		for _, s := range g.Services {
			if _, ok := expectedServices[s]; ok {
				expectedServices[s] = true
			}
		}
	}

	for k, v := range expectedServices {
		if !v {
			t.Fatalf("InappropriateIntimacyServiceAnnotator().Service = %v doesn't exist", k)
		}
	}

	annotations = append(annotations, inappropriateIntimacy...)
}

func TestGreedyServiceAnnotator(t *testing.T) {
	greedy := annotators.GreedyServiceAnnotator(operations, services)

	annotations = append(annotations, greedy...)
}

//func TestAnalyzeDependenceAnnotator(t *testing.T) {
//	dependence := annotators.AbsoluteDependenceService(services)
//
//	expectedService := "service-e"
//
//	if dependence.Services[0] != expectedService {
//		t.Log(dependence.Message)
//		t.Fatalf(`AbsoluteDependenceAnnotator() = %s, want %s`, dependence.Services[0], expectedService)
//	}
//
//	expectedMessage := "Service service-e has 2 dependencies"
//
//	if dependence.Message != expectedMessage {
//		t.Fatalf(`AbsoluteDependenceAnnotator().Message = %s, want %s`, dependence.Message, expectedMessage)
//	}
//
//	annotations = append(annotations, dependence)
//}
//
//func TestAnalyzeCriticalityAnnotator(t *testing.T) {
//	criticality := annotators.AbsoluteCriticalService(services)
//
//	expectedService := "service-f"
//
//	if criticality.Services[0] != expectedService {
//		t.Fatalf(`AbsoluteCriticalityAnnotator() = %s, want %s`, criticality.Services[0], expectedService)
//	}
//
//	expectedMessage := "Service service-e has 2 dependents and 3 dependencies"
//
//	if criticality.Message != expectedMessage {
//		t.Fatalf(`AbsoluteCriticalityAnnotator().Message = %s, want %s`, criticality.Message, expectedMessage)
//	}
//
//	annotations = append(annotations, criticality)
//}

func TestPrintToJSON(t *testing.T) {
	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}

	file, _ := json.MarshalIndent(yCharModel, "", " ")

	_ = os.WriteFile("../frontend/src/y-chart.json", file, 0644)
}

func TestRecommendationEngine(t *testing.T) {
	for i, a := range annotations {
		if a.AnnotationType == models.Megaservice {
			services, operations, annotations[i] = recommenders.MegaserviceRecommender(services, operations, a)
		}
	}

	annotations = Analyze(operations, services)

	for i, a := range annotations {
		if a.AnnotationType == models.InappropriateIntimacy {
			services, operations, annotations[i] = recommenders.InappropriateIntimacyRecommender(services, operations, a)
		}
	}

	annotations = Analyze(operations, services)

	for i, a := range annotations {
		if a.AnnotationType == models.Greedy {
			services, operations, annotations[i] = recommenders.GreedyServiceRecommender(services, operations, a)
		}
	}

	annotations = Analyze(operations, services)

	for i, a := range annotations {
		if a.AnnotationType == models.Cyclic {
			services, operations, annotations[i] = recommenders.CyclicRecommender(services, operations, a)
		}
	}

	annotations = Analyze(operations, services)

	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}
	file, _ := json.MarshalIndent(yCharModel, "", "    ")
	_ = os.WriteFile("../frontend/src/y-chart-recommendation.json", file, 0644)

}
