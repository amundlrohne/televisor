package queries

import (
	"context"
	"encoding/json"
	"io"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	request := &pb.FindTracesRequest{Query: &pb.TraceQueryParameters{
		ServiceName:   service,
		OperationName: operation,
		StartTimeMin:  time.Now().Add(-time.Hour * 24),
		StartTimeMax:  time.Now(),
		DurationMin:   time.Duration(0),
		DurationMax:   time.Hour,
		SearchDepth:   1,
	}}

	stream, err := qsc.FindTraces(ctx, request)

	if err != nil {
		log.Fatalf("could not get stream: %v", err)
	}

	var result []model.Span

	for {
		spans, err := stream.Recv()

		if err == io.EOF {
			return result
		} else if err == nil {
			result = append(result, spans.Spans...)
		}

		if err != nil {
			log.Fatalf("stream failed: %v", err)
			return []model.Span{}
		}
	}
}

func JaegerTrace(qsc pb.QueryServiceClient, traceId model.TraceID) []model.Span {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	request := &pb.GetTraceRequest{TraceID: traceId}

	jsonRequest, _ := json.Marshal(request)

	var newRequest pb.GetTraceRequest
	json.Unmarshal(jsonRequest, &newRequest)
	client, err := qsc.GetTrace(ctx, &newRequest)

	if err != nil {
		log.Fatalf("Trace doesn't exist: %v", err)
	}

	r, err := client.Recv()

	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	return r.Spans
}
