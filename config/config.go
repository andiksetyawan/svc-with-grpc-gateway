package config

import (
	"svc-with-grpc-gateway/pkg/helper"
)

var (
	ServiceName           = helper.GetEnv("SERVICE_NAME", "svc-with-grpc-gateway")
	Address               = helper.GetEnv("ADDRESS", ":8080")
	OtlpCollectorUrl      = helper.GetEnv("OTLP_COLLECTOR_URL", "localhost:4317")
	InsecureOtlpCollector = helper.GetEnv("INSECURE_OTLP_COLLECTOR", "true")
)
