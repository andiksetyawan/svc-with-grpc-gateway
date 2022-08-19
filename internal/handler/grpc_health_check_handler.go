package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type healthCheckServerHandler struct {
}

func NewHealthCheckServerHandler() health.HealthServer {
	return &healthCheckServerHandler{}
}

func (s *healthCheckServerHandler) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	log.Debug().Msg("healthCheckServerHandler Check")

	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (s *healthCheckServerHandler) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
	// Example of how to register both methods but only implement the Check method.
	return status.Error(codes.Unimplemented, "unimplemented")
}
