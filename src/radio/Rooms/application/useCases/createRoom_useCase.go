package usecases

import (
	"context"
	"errors"
	"lively-backend/src/radio/Rooms/domain/entity"
	"lively-backend/src/radio/Rooms/domain/repository"
)

type CreateRoomInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   int    `json:"created_by"`
}

type CreateRoomUseCase struct {
	roomRepo repository.IRoomRepository
}

func NewCreateRoomUseCase(repo repository.IRoomRepository) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomRepo: repo,
	}
}

func (uc *CreateRoomUseCase) Execute(ctx context.Context, input CreateRoomInput) (*entity.Room, error) {

	if input.Name == "" {
		return nil, errors.New("el nombre de la sala es obligatorio")
	}

	if input.CreatedBy <= 0 {
		return nil, errors.New("se requiere el ID de un usuario válido para crear la sala")
	}

	room := &entity.Room{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	}

	err := uc.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, errors.New("no se pudo crear la sala. Es posible que el ID ya exista")
	}

	return room, nil
}
