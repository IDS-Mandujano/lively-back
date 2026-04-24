package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/radio/Rooms/application/useCases"
	"log"
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
	log.Printf("[rooms:create] request method=%s remote=%s", r.Method, r.RemoteAddr)
	if r.Method != http.MethodPost {
		log.Printf("[rooms:create] rejected method=%s", r.Method)
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input usecases.CreateRoomInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("[rooms:create] invalid json: %v", err)
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	log.Printf("[rooms:create] payload name=%q created_by=%d", input.Name, input.CreatedBy)

	room, err := c.useCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("[rooms:create] usecase error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[rooms:create] success id=%d name=%q", room.ID, room.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sala creada exitosamente",
		"room":    room,
	})
}
