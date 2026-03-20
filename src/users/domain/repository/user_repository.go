package repository

import (
	"context"
	"lively-backend/src/users/domain/entity"
)

// IUserRepository define las operaciones obligatorias para los usuarios
type IUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
