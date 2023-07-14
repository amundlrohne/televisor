package main

import (
	"flag"
	"fmt"

	pb "jaeger-idl/api_v2"

	"github.com/amundlrohne/televisor/annotators"
	"github.com/amundlrohne/televisor/connectors"
	"github.com/amundlrohne/televisor/utils"
)

var (
	jaeger_addr = flag.String("jaeger_addr", "localhost:16685", "jaeger address to connect to")
	api_gateway = flag.String("api_gateway", "nginx-web-server", "api gateway in microservice application")
)

func main() {
	analyze()
}

func analyze() {
	//queries.PrometheusContainerCPU()
	//queries.PrometheusContainerMemory()
	//queries.PrometheusContainerNetworkInput()
	//queries.PrometheusContainerNetworkOutput()

	flag.Parse()
	// Set up a connection to the Jaeger Server.
	conn := connectors.JaegerConnect(*jaeger_addr)
	defer conn.Close()

	qsc := pb.NewQueryServiceClient(&conn)

	/* log.Printf("SDG: %v", queries.JaegerSDG(qsc))
	log.Printf("Services: %v", queries.JaegerServices(qsc))
	log.Printf("Operations: %v", queries.JaegerOperations(qsc, "nginx-web-server"))
	log.Printf("Traces: %+v", queries.JaegerTraces(qsc)) */

	operations := utils.GetSubSDGs(qsc, *api_gateway)
	combinedEdges := operations.CombineEdges()
	services := utils.ExtractServicesFromSDG(combinedEdges)

	megaservices := annotators.MegaserviceAnnotator(operations)
	fmt.Printf("Megaservices: %+v \n", megaservices)

	greedy := annotators.GreedyServiceAnnotator(operations)
	fmt.Printf("Greedy: %+v \n", greedy)

	criticality := annotators.AbsoluteCriticalService(services)
	fmt.Printf("Criticality %+v \n", criticality)

	dependence := annotators.AbsoluteDependenceService(services)
	fmt.Printf("Dependence %+v \n", dependence)

	cycles := annotators.CyclicDependencyAnnotator(operations, services)
	fmt.Printf("Cycles %+v \n", cycles)

	//fmt.Printf("Connected: %+v", operations["/wrk2-api/home-timeline/read"].IsConnected("nginx-web-server", "post-storage-service"))

}

/* func main() {
	// Hello world, the web server
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
*/
