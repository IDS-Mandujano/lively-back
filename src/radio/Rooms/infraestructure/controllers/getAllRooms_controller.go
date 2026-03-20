package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/radio/Rooms/application/useCases"
	"lively-backend/src/radio/Rooms/domain/entity"
	"net/http"
)

type GetAllRoomsController struct {
	useCase *usecases.GetAllRoomsUseCase
}

func NewGetAllRoomsController(uc *usecases.GetAllRoomsUseCase) *GetAllRoomsController {
	return &GetAllRoomsController{
		useCase: uc,
	}
}

func (c *GetAllRoomsController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rooms, err := c.useCase.Execute(r.Context())
	if err != nil {
		http.Error(w, "Error obteniendo las salas", http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []entity.Room{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}
