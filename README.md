# televisor
Kubernetes Operator identifying performance bottle-necks in microservice cluster based on telemetry from Jaeger and Prometheus

## Setup 
In order to connect to the gRPC client, you first must update both sub-modules (backend/jaeger-idl and backend/jaeger-idl/opentelemetry-proto) and create a go.mod file with the ```module jaeger-idl``` tag. 