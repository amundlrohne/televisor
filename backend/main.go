package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "jaeger-idl/api_v2"

	"github.com/amundlrohne/televisor/connectors"
	"github.com/amundlrohne/televisor/queries"
)

var (
	jaeger_addr = flag.String("jaeger_addr", "localhost:16685", "jaeger address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the Jaeger Server.
	conn := connectors.JaegerConnect(*jaeger_addr)
	defer conn.Close()

	qsc := pb.NewQueryServiceClient(&conn)

	log.Printf("SDG: %v", queries.JaegerSDG(qsc))
	log.Printf("Services: %v", queries.JaegerServices(qsc))

	res, err := http.Get("http://localhost:9090/api/v1/query?query=sum%20by%20(cpu)%20(process_cpu_seconds_total{mode!=%22idle%22})")
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

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
