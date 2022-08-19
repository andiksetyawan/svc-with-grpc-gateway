package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	svcUserV1 "svc-with-grpc-gateway/api/v1"
	"svc-with-grpc-gateway/internal/model/entity"
	"svc-with-grpc-gateway/internal/service"
)

type userHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) svcUserV1.UserServiceServer {
	return &userHandler{userService: userService}
}

func (u *userHandler) Create(ctx context.Context, request *svcUserV1.CreateRequest) (*svcUserV1.CreateResponse, error) {
	_, err := u.userService.Create(ctx, &entity.User{Name: request.Name})
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return &svcUserV1.CreateResponse{
		Error:   false,
		Message: "OK",
	}, nil
}
