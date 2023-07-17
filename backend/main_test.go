package main

import (
	"testing"

	"github.com/amundlrohne/televisor/annotators"
	"github.com/amundlrohne/televisor/generators"
	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/utils"
)

var operations models.Operations
var services map[string]models.TelevisorService

func init() {
	operations = generators.CombinedGenerator()
	combinedEdges := operations.CombineEdges()
	services = utils.ExtractServicesFromSDG(combinedEdges)
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

}

func TestAnalyzeGreedyServiceAnnotator(t *testing.T) {
	greedy := annotators.GreedyServiceAnnotator(operations)

	if len(greedy) != 1 {
		t.Fatalf("GreedyServiceAnnotator() = %v, want 1, error", len(greedy))
	}

	expectedServices := map[string]bool{"api-gateway": false, "service-e": false}

	for _, s := range greedy[0].Services {
		if _, ok := expectedServices[s]; ok {
			expectedServices[s] = true
		}
	}

	for k, v := range expectedServices {
		if !v {
			t.Fatalf("GreedyServiceAnnotator().Service = %v doesn't exist", k)
		}
	}
	expectedOperations := map[string]bool{"op2-subop1": false, "op2-subop6": false}

	for _, s := range greedy[0].Operations {
		if _, ok := expectedOperations[s]; ok {
			expectedOperations[s] = true
		}
	}

	for k, v := range expectedOperations {
		if !v {
			t.Fatalf("GreedyServiceAnnotator().Operations = %v doesn't exist", k)
		}
	}
}

func TestCyclicServiceAnnotator(t *testing.T) {
	cyclic := annotators.CyclicDependencyAnnotator(operations, services)

	if len(cyclic) != 2 {
		t.Fatalf("CyclicServiceAnnotator() = %v, want 2, error", len(cyclic))
	}

	expectedCycles := map[string]bool{"service-e": false, "service-a": false}

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
}

func TestAnalyzeDependenceAnnotator(t *testing.T) {
	dependence := annotators.AbsoluteDependenceService(services)

	expectedService := "service-a"

	if dependence.Services[0] != expectedService {
		t.Log(dependence.Message)
		t.Fatalf(`AbsoluteDependenceAnnotator() = %s, want %s`, dependence.Services[0], expectedService)
	}

	expectedMessage := "Service service-a has 2 dependencies"

	if dependence.Message != expectedMessage {
		t.Fatalf(`AbsoluteDependenceAnnotator().Message = %s, want %s`, dependence.Message, expectedMessage)
	}
}

func TestAnalyzeCriticalityAnnotator(t *testing.T) {
	criticality := annotators.AbsoluteCriticalService(services)

	expectedService := "service-b"

	if criticality.Services[0] != expectedService {
		t.Fatalf(`AbsoluteCriticalityAnnotator() = %s, want %s`, criticality.Services[0], expectedService)
	}

	expectedMessage := "Service service-b has 3 dependents and 1 dependencies"

	if criticality.Message != expectedMessage {
		t.Fatalf(`AbsoluteCriticalityAnnotator().Message = %s, want %s`, criticality.Message, expectedMessage)
	}
}
