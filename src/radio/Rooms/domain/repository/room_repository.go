package repository

import (
	"context"
	"lively-backend/src/radio/Rooms/domain/entity"
)

type IRoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	GetAll(ctx context.Context) ([]entity.Room, error)
}
