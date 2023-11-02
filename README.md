# televisor
Kubernetes Operator identifying and suggesting mitigation techniques for architectural anti-patterns in microservice applications based on telemetry from Jaeger and Prometheus.

## Setup
This project requires Go version **1.20**, and Node version **v18** or newer.


In order to connect to the gRPC client, you first must update both Git sub-modules (backend/jaeger-idl and backend/jaeger-idl/opentelemetry-proto). Then follow the steps in ```backend/jaeger-idl/README.md``` to generate the necessary protobuf code.

This is done with the ```git submodule update``` command.

To let Televisor interface with the generated code create a ```go.mod``` file with the ```module jaeger-idl``` tag  in the ```proto-gen-go``` directory and add all files with a ```*.pb.go.``` into a new ```api_v2``` directory in the same location. With the exception of the ```model.pb.go``` file, that has to placed in a separate ```/proto-gen-go/api_v2/model``` directory.

As mentioned, Televisor is dependent on telemetry instrumentation in the target application. The microservice application needs to be instrumented with distributed tracing (OpenTelemetry), it needs to run Jaeger and also needs a cAdvisor instance on the platform. Televisor targets a Prometheus instance running on either the host machine or locally. However, Prometheus needs to target the cAdvisor instance on the host machine regardless of where it is run. This can be configured in ```prometheus.yml```.


## Target Application
Note that Jaeger does not expose the :16685 port by default which is required by Televisor to pull telemetry over gRPC.

## Usage
Televisor is a two module application. First we execute the backend module in order to gather telemetry, generate mitigation suggestions, and produce two JSON files with the results.

To run the backend module after the necessary setup steps have been made execute the following command in the ```backend``` directory. Jaeger and Prometheus flags can be omitted if localhost with default ports is wanted.

```
go run main.go --jaeger_addr <jaeger address> --prom_addr <prometheus address> --api_gateway <service name of api gateway>
```

This will produce two JSON files in the ```frontend/src``` directory. While the files can be examined manually, it is recommended to use the frontend module for visualization.

To run the frontend module, navigate to the ```frontend``` directory and run:

```
npm start
```

The frontend visualization will then be available on ```localhost:3000```.
