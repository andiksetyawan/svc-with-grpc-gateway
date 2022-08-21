//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	svcUserV1 "svc-with-grpc-gateway/api/v1"
	"svc-with-grpc-gateway/internal/handler"
	"svc-with-grpc-gateway/internal/repository"
	"svc-with-grpc-gateway/internal/service"
	"svc-with-grpc-gateway/internal/store"
)

var userSets = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
)

func InitializedUserServiceHandlerServer() svcUserV1.UserServiceServer {
	wire.Build(
		store.NewSQLLite,
		userSets,
	)
	return nil
}
