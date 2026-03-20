package repository

import (
	"context"
	"lively-backend/src/radio/Stations/domain/entity"
)

type IStationRepository interface {
	GetByCategory(ctx context.Context, category string, limit int) ([]entity.Station, error)
	SearchByName(ctx context.Context, name string, limit int) ([]entity.Station, error)
	GetTop(ctx context.Context, limit int) ([]entity.Station, error)
}
