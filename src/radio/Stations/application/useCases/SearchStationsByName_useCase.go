package usecases

import (
	"context"
	"errors"
	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"
)

type SearchStationsByNameUseCase struct {
	stationRepo repository.IStationRepository
}

func NewSearchStationsByNameUseCase(repo repository.IStationRepository) *SearchStationsByNameUseCase {
	return &SearchStationsByNameUseCase{
		stationRepo: repo,
	}
}

func (uc *SearchStationsByNameUseCase) Execute(ctx context.Context, name string, limit int) ([]entity.Station, error) {
	if name == "" {
		return nil, errors.New("el término de búsqueda no puede estar vacío")
	}

	if limit <= 0 {
		limit = 20
	}

	return uc.stationRepo.SearchByName(ctx, name, limit)
}
