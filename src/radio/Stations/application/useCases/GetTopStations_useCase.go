package usecases

import (
	"context"
	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"

)

type GetTopStationsUseCase struct {
	stationRepo repository.IStationRepository
}

func NewGetTopStationsUseCase(repo repository.IStationRepository) *GetTopStationsUseCase {
	return &GetTopStationsUseCase{
		stationRepo: repo,
	}
}

func (uc *GetTopStationsUseCase) Execute(ctx context.Context, limit int) ([]entity.Station, error) {
	if limit <= 0 {
		limit = 20
	}

	return uc.stationRepo.GetTop(ctx, limit)
}
