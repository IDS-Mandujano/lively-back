package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/users/application/useCases"
	"net/http"
)

type RegisterUserController struct {
	useCase *usecases.RegisterUserUseCase
}

func NewRegisterUserController(uc *usecases.RegisterUserUseCase) *RegisterUserController {
	return &RegisterUserController{
		useCase: uc,
	}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input usecases.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	user, err := c.useCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuario registrado exitosamente",
		"user":    user,
	})
}
