package repository

import (
	"context"
	"lively-backend/src/users/domain/entity"
)

type IUserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
