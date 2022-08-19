package repository

import (
	"context"

	"gorm.io/gorm"
	"svc-with-grpc-gateway/internal/model/entity"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
