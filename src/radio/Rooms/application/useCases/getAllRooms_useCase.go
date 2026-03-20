package usecases

import (
	"context"
	"lively-backend/src/radio/Rooms/domain/entity"
	"lively-backend/src/radio/Rooms/domain/repository"
)

type GetAllRoomsUseCase struct {
	roomRepo repository.IRoomRepository
}

func NewGetAllRoomsUseCase(repo repository.IRoomRepository) *GetAllRoomsUseCase {
	return &GetAllRoomsUseCase{
		roomRepo: repo,
	}
}

func (uc *GetAllRoomsUseCase) Execute(ctx context.Context) ([]entity.Room, error) {
	return uc.roomRepo.GetAll(ctx)
}
