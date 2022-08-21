package service

import (
	"context"

	"go.opentelemetry.io/otel"
	"svc-with-grpc-gateway/config"
	"svc-with-grpc-gateway/internal/model/entity"
	"svc-with-grpc-gateway/internal/repository"
)

type IUserService interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
}

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &userService{userRepository: userRepository}
}

func (u *userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	//Tracer
	_, span := otel.Tracer(config.App.ServiceName).Start(ctx, "service.user.Create")
	defer span.End()

	//TODO business logic
	return u.userRepository.Create(ctx, user)
}
