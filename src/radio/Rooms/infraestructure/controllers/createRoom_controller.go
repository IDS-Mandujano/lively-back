package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/radio/Rooms/application/useCases"
	"net/http"
)

type CreateRoomController struct {
	useCase *usecases.CreateRoomUseCase
}

func NewCreateRoomController(uc *usecases.CreateRoomUseCase) *CreateRoomController {
	return &CreateRoomController{
		useCase: uc,
	}
}

func (c *CreateRoomController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input usecases.CreateRoomInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	room, err := c.useCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sala creada exitosamente",
		"room":    room,
	})
}
