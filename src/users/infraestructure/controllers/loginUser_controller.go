package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/users/application/useCases"
	"net/http"
)

type LoginUserController struct {
	useCase *usecases.LoginUserUseCase
}

func NewLoginUserController(uc *usecases.LoginUserUseCase) *LoginUserController {
	return &LoginUserController{
		useCase: uc,
	}
}

func (c *LoginUserController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input usecases.LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	token, user, err := c.useCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Le enviamos a Android su Token y sus datos
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Inicio de sesión exitoso",
		"token":   token,
		"user":    user,
	})
}
