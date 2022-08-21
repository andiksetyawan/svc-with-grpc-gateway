# svc-with-grpc-gateway
A simple grpc+http1 service using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) and listening on the same port.
- [Installation](#installation)
- [Usage](#usage)
- [Documentation](#documentation)

## Installation

### Requirements

1. [Go](https://golang.org/doc/install) 1.16+
2. [Open Telemetry Collector](https://opentelemetry.io/docs/collector/getting-started/), Tracer & metrics collector
3. [Docker](https://docs.docker.com/engine/install/), used for testing with [testcontainer](https://www.testcontainers.org/) 
4. [Buf](https://docs.buf.build/introduction) for generating the grpc stubs
5. [protoc-gen-go](#), [protoc-gen-go-grpc](%), [](),

### Setting up environment
Default env:
```
SERVICE_NAME=svc-with-grpc-gateway
ADDRESS=:8080
OTLP_COLLECTOR_URL=localhost:4317
```
you can setup environment using .env file or environment variables (OS)
## Getting Started
## Usage

### Development
```
 go mod tidy
```

After all installed properly, because this project we use opentelemetry collector (jaeger + prometheus) for tracing and metrics, so Run the otelcollector, jaeger and prometheus servers using docker-compose:

```docker-compose -f deploy/docker-compose/docker-compose.yaml up```


and then start:

```
 go run cmd/server/main.go
```

#### Make requests:
HTTP/1.1 POST API with curl
```
$ curl \
--header "Content-Type: application/json" \
--data '{"name": "John"}' \
http://localhost:8080/user.v1.UserService/Create
```
gRPC with grpccurl
```
$ grpcurl \
-protoset <(buf build -o -) -plaintext \
-d '{"name": "John"}' \
localhost:8080 user.v1.UserService/Create
```

reponse:
```
{"message": "OK"}
```

#### Make test:
```
go test ./...
```

## Documentation
### Api specs:
openapi:
```api/v1/user.swagger.json```

proto :
```api/v1/user.proto```

