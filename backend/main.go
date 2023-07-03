package main

import (
	"flag"
	"fmt"

	pb "jaeger-idl/api_v2"

	"github.com/amundlrohne/televisor/connectors"
	"github.com/amundlrohne/televisor/models"
	"github.com/amundlrohne/televisor/queries"
)

var (
	jaeger_addr = flag.String("jaeger_addr", "localhost:16685", "jaeger address to connect to")
)

func main() {
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

	var sdg = make(map[string]models.TelevisorService)

	services := queries.JaegerServices(qsc)

	for _, s := range services {
		operations := queries.JaegerOperations(qsc, s)
		relationships := make(map[string]models.TelevisorRelationship)
		for _, o := range operations {
			traces := queries.JaegerTraces(qsc, s, o)
			for _, t := range traces {
				if relationship, ok := relationships[t.OperationName]; ok {
					relationship.Count++
					relationships[t.OperationName] = relationship
				} else {
					relationship := models.TelevisorRelationship{
						Count:       1,
						ServiceName: t.Process.ServiceName,
					}
					relationships[t.OperationName] = relationship
				}
			}
		}

		service := models.TelevisorService{Name: s, Operations: operations, Relationships: relationships}
		sdg[s] = service
	}

	fmt.Printf("Services: %+v", sdg["nginx-web-server"])
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
