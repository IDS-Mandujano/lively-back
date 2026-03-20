package repository

import (
	"context"
	"lively-backend/src/radio/Rooms/domain/entity"
)

// IRoomRepository define cómo interactuamos con la tabla de salas en la base de datos
type IRoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	GetAll(ctx context.Context) ([]entity.Room, error)
}
