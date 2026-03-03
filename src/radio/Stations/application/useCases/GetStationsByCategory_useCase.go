package usecases

import (
	"context"
	// Recuerda cambiar "lively" por el nombre de tu módulo real
	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"
)

// Definimos la estructura del caso de uso inyectando la interfaz del repositorio
type GetStationsByCategoryUseCase struct {
	stationRepo repository.IStationRepository
}

// Constructor "New" para instanciar el caso de uso
func NewGetStationsByCategoryUseCase(repo repository.IStationRepository) *GetStationsByCategoryUseCase {
	return &GetStationsByCategoryUseCase{
		stationRepo: repo,
	}
}

// Execute es el método que llamará tu controlador
func (uc *GetStationsByCategoryUseCase) Execute(ctx context.Context, category string, limit int) ([]entity.Station, error) {

	// Aquí podríamos agregar lógica de negocio pura si la tuviéramos.
	// Por ejemplo: validar que la categoría no esté vacía.
	if category == "" {
		category = "pop" // Un valor por defecto si el frontend manda vacío
	}

	if limit <= 0 {
		limit = 20 // Límite por defecto
	}

	// Delegamos la obtención de datos al puerto (repositorio)
	stations, err := uc.stationRepo.GetByCategory(ctx, category, limit)
	if err != nil {
		return nil, err
	}

	return stations, nil
}
