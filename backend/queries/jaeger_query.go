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
	request := &pb.GetDependenciesRequest{StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)}
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
