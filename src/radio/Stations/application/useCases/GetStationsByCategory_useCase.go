package usecases

import (
	"context"
	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"
)

type GetStationsByCategoryUseCase struct {
	stationRepo repository.IStationRepository
}

func NewGetStationsByCategoryUseCase(repo repository.IStationRepository) *GetStationsByCategoryUseCase {
	return &GetStationsByCategoryUseCase{
		stationRepo: repo,
	}
}

func (uc *GetStationsByCategoryUseCase) Execute(ctx context.Context, category string, limit int) ([]entity.Station, error) {

	if category == "" {
		category = "pop"
	}

	if limit <= 0 {
		limit = 20
	}

	stations, err := uc.stationRepo.GetByCategory(ctx, category, limit)
	if err != nil {
		return nil, err
	}

	return stations, nil
}
