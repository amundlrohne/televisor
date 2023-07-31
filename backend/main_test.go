package main

import (
	"encoding/json"
	"io/ioutil"
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

	expectedService := "service-d"

	if megaservice[0].Services[0] != expectedService {
		t.Fatalf("MegaserviceAnnotator().Services[0] = %s, want %s, error", megaservice[0].Services[0], expectedService)
	}

	annotations = append(annotations, megaservice...)
}

func TestAnalyzeInappropriateIntimacyServiceAnnotator(t *testing.T) {
	greedy := annotators.InappropriateIntimacyServiceAnnotator(operations)

	if len(greedy) != 2 {
		t.Fatalf("InappropriateIntimacyServiceAnnotator() = %v, want 2, error", len(greedy))
	}

	expectedServices := map[string]bool{"api-gateway": false, "service-f": false, "service-a": false, "service-b": false}

	for _, g := range greedy {
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

	expectedOperations := map[string]bool{"op2-subop1": false, "op2-subop3": false, "op4-subop4": false, "op4-subop5": false}

	for _, g := range greedy {
		for _, o := range g.Operations {
			if _, ok := expectedOperations[o]; ok {
				expectedOperations[o] = true
			}
		}
	}

	for k, v := range expectedOperations {
		if !v {
			t.Fatalf("InappropriateIntimacyServiceAnnotator().Operations = %v doesn't exist", k)
		}
	}

	annotations = append(annotations, greedy...)
}

/* func TestCyclicServiceAnnotator(t *testing.T) {
	cyclic := annotators.CyclicDependencyAnnotator(operations, services)

	if len(cyclic) != 2 {
		t.Fatalf("CyclicServiceAnnotator() = %v, want 2, error", len(cyclic))
	}

	expectedCycles := map[string]bool{"service-b": false, "service-f": false}

	for _, a := range cyclic {
		service := a.Services[0]
		if _, ok := expectedCycles[service]; ok {
			expectedCycles[service] = true
		}
	}

	for k, v := range expectedCycles {
		if !v {
			t.Fatalf("CyclicServiceAnnotator().Service = %v doesn't exist", k)
		}
	}

	annotations = append(annotations, cyclic...)
} */

func TestAnalyzeDependenceAnnotator(t *testing.T) {
	dependence := annotators.AbsoluteDependenceService(services)

	expectedService := "service-e"

	if dependence.Services[0] != expectedService {
		t.Log(dependence.Message)
		t.Fatalf(`AbsoluteDependenceAnnotator() = %s, want %s`, dependence.Services[0], expectedService)
	}

	expectedMessage := "Service service-e has 3 dependencies"

	if dependence.Message != expectedMessage {
		t.Fatalf(`AbsoluteDependenceAnnotator().Message = %s, want %s`, dependence.Message, expectedMessage)
	}

	annotations = append(annotations, dependence)
}

func TestAnalyzeCriticalityAnnotator(t *testing.T) {
	criticality := annotators.AbsoluteCriticalService(services)

	expectedService := "service-b"

	if criticality.Services[0] != expectedService {
		t.Fatalf(`AbsoluteCriticalityAnnotator() = %s, want %s`, criticality.Services[0], expectedService)
	}

	expectedMessage := "Service service-b has 4 dependents and 2 dependencies"

	if criticality.Message != expectedMessage {
		t.Fatalf(`AbsoluteCriticalityAnnotator().Message = %s, want %s`, criticality.Message, expectedMessage)
	}

	annotations = append(annotations, criticality)
}

func TestPrintToJSON(t *testing.T) {
	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}

	file, _ := json.MarshalIndent(yCharModel, "", " ")

	_ = ioutil.WriteFile("../y-chart-test.json", file, 0644)
}

func TestRecommendationEngine(t *testing.T) {
	for _, a := range annotations {
		if a.AnnotationType == models.Megaservice {
			ss, o := recommenders.MegaserviceRecommender(services[a.Services[0]], operations[a.InitiatingOperation])
			delete(services, a.Services[0])
			for _, v := range ss {
				services[v.Name] = v
			}
			operations[a.InitiatingOperation] = o
		}
	}

	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}
	file, _ := json.MarshalIndent(yCharModel, "", " ")
	_ = ioutil.WriteFile("../y-chart-recommendation.json", file, 0644)

}
