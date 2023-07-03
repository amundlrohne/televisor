package queries

import (
	"context"
	"log"
	"time"

	pb "jaeger-idl/api_v2"

	"github.com/jaegertracing/jaeger/model"
)

// Currently queries 1 hour back in time
func JaegerSDG(qsc pb.QueryServiceClient) []model.DependencyLink {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// StartTime = EndTs && EndTime = loopback (Diff between StartTime and EndTime) BUG??
	request := &pb.GetDependenciesRequest{StartTime: time.Now(), EndTime: time.Now().Add(time.Hour * 24)}
	//log.Printf("EndTs: %v", request.GetStartTime().UnixMilli())
	//log.Printf("Loopback: %v", request.GetEndTime().Sub(request.GetStartTime()).Milliseconds())
	r, err := qsc.GetDependencies(ctx, request)

	if err != nil {
		log.Fatalf("could not get sdg: %v", err)
	}

	return r.Dependencies
}

func JaegerServices(qsc pb.QueryServiceClient) []string {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := qsc.GetServices(ctx, &pb.GetServicesRequest{})

	if err != nil {
		log.Fatalf("could not get services: %v", err)
	}

	return r.Services
}

func JaegerOperations(qsc pb.QueryServiceClient, service string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := qsc.GetOperations(ctx, &pb.GetOperationsRequest{Service: service})

	if err != nil {
		log.Fatalf("could not get operations: %v", err)
	}

	return r.OperationNames
}

func JaegerTraces(qsc pb.QueryServiceClient, service string, operation string) []model.Span {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	request := &pb.FindTracesRequest{Query: &pb.TraceQueryParameters{ServiceName: "nginx-web-server", OperationName: "/wrk2-api/post/compose", SearchDepth: 0}}

	client, err := qsc.FindTraces(ctx, request)

	if err != nil {
		log.Fatalf("could not get traces: %v", err)
	}

	r, err := client.Recv()

	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	return r.Spans
}

func JaegerTrace(qsc pb.QueryServiceClient, traceId model.TraceID) []model.Span {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	request := &pb.GetTraceRequest{TraceID: traceId}

	client, err := qsc.GetTrace(ctx, request)

	if err != nil {
		log.Fatalf("Trace doesn't exist: %v", err)
	}

	r, err := client.Recv()

	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	return r.Spans
}
