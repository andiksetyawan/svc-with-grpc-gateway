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

func InitializedServer() svcUserV1.UserServiceServer {
	wire.Build(
		store.NewSQLLite,
		userSets,
		//grpc.NewServer,
		////wire.Bind(new(http.Handler), new(*runtime.ServeMux)),
		//runtime.NewServeMux,
		////wire.Bind(new(*runtime.ServeMux), new(http.Handler)),
		//middleware.GrpcHttpMiddleware,
	)
	return nil
}
