package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	pb "jaeger-idl/api_v2"

	"github.com/amundlrohne/televisor/annotators"
	"github.com/amundlrohne/televisor/connectors"
	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/queries"
	"github.com/amundlrohne/televisor/recommenders"
	"github.com/amundlrohne/televisor/utils"
)

var (
	jaeger_addr = flag.String("jaeger_addr", "localhost:16685", "jaeger address to connect to")
	api_gateway = flag.String("api_gateway", "nginx-web-server", "api gateway in microservice application")
	cpu_req     = flag.Float64("cpu_req", 0.6, "cpu utilization requirement")
	mem_req     = flag.Float64("mem_req", 0.2, "memory utilization requirement")
)

func main() {
	operations, services := retrieveTelemetry()

	annotations := Analyze(operations, services)
	recommend(operations, services, annotations)
}

func recommend(operations models.Operations, services map[string]models.TelevisorService, annotations []models.Annotation) {
	for i, a := range annotations {
		if a.AnnotationType == models.Megaservice {
			services, operations, annotations[i] = recommenders.MegaserviceRecommender(services, operations, a)
		}

		if a.AnnotationType == models.InappropriateIntimacy {
			services, operations, annotations[i] = recommenders.InappropriateIntimacyRecommender(services, operations, a)
		}
	}

	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}
	file, _ := json.MarshalIndent(yCharModel, "", " ")
	_ = ioutil.WriteFile("../y-chart-recommendation.json", file, 0644)

}

func retrieveTelemetry() (models.Operations, map[string]models.TelevisorService) {
	flag.Parse()
	// Set up a connection to the Jaeger Server.
	conn := connectors.JaegerConnect(*jaeger_addr)
	defer conn.Close()

	qsc := pb.NewQueryServiceClient(&conn)

	cpuUtilsQuantile := queries.PrometheusContainerCPUQuantile()
	cpuUtilsMean := queries.PrometheusContainerCPUMean()
	cpuUtilsStdev := queries.PrometheusContainerCPUStdev()

	memoryUtilsQuantile := queries.PrometheusContainerMemoryQuantile()
	memoryUtilsMean := queries.PrometheusContainerMemoryMean()
	memoryUtilsStdev := queries.PrometheusContainerMemoryStdev()

	// networkInUtils := queries.PrometheusContainerNetworkInput()
	// networkOutUtils := queries.PrometheusContainerNetworkOutput()

	operations := utils.GetSubSDGs(qsc, *api_gateway)
	combinedEdges := operations.CombineEdges()
	services := utils.ExtractServicesFromSDG(combinedEdges)

	services = utils.AddCPUQuantileToServices(services, cpuUtilsQuantile)
	services = utils.AddCPUMeanToServices(services, cpuUtilsMean)
	services = utils.AddCPUStdevToServices(services, cpuUtilsStdev)

	services = utils.AddMemoryQuantileToServices(services, memoryUtilsQuantile)
	services = utils.AddMemoryMeanToServices(services, memoryUtilsMean)
	services = utils.AddMemoryStdevToServices(services, memoryUtilsStdev)

	return operations, services
}

func Analyze(operations models.Operations, services map[string]models.TelevisorService) []models.Annotation {

	annotations := []models.Annotation{}

	// fmt.Printf("Services %+v \n", services)

	megaservices := annotators.MegaserviceAnnotator(operations)
	annotations = append(annotations, megaservices...)
	// fmt.Printf("Megaservices: %+v \n", megaservices)

	cycles := annotators.CyclicDependencyAnnotator(operations, services)
	annotations = append(annotations, cycles...)
	// fmt.Printf("Cycles %+v \n", cycles)

	innapropriateIntimacy := annotators.InappropriateIntimacyServiceAnnotator(operations, cycles)
	annotations = append(annotations, innapropriateIntimacy...)
	// fmt.Printf("Innapropriate Intimacy: %+v \n", innapropriateIntimacy)

	greedy := annotators.GreedyServiceAnnotator(operations, services)
	annotations = append(annotations, greedy...)
	// fmt.Printf("Greedy %+v \n", greedy)

	criticality := annotators.AbsoluteCriticalService(services)
	annotations = append(annotations, criticality)
	// fmt.Printf("Criticality %+v \n", criticality)

	dependence := annotators.AbsoluteDependenceService(services)
	annotations = append(annotations, dependence)
	// fmt.Printf("Dependence %+v \n", dependence)

	yCharModel := models.YChartModel{Annotations: annotations, Operations: operations, Services: services}
	file, _ := json.MarshalIndent(yCharModel, "", " ")
	_ = ioutil.WriteFile("../y-chart.json", file, 0644)

	return annotations

}
