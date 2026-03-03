package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/radio/Stations/application/useCases"
	"net/http"
	"strconv"
)

type GetTopStationsController struct {
	useCase *usecases.GetTopStationsUseCase
}

func NewGetTopStationsController(uc *usecases.GetTopStationsUseCase) *GetTopStationsController {
	return &GetTopStationsController{
		useCase: uc,
	}
}

func (c *GetTopStationsController) Handle(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	stations, err := c.useCase.Execute(r.Context(), limit)
	if err != nil {
		http.Error(w, "Error obteniendo top estaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stations); err != nil {
		http.Error(w, "Error formateando la respuesta", http.StatusInternalServerError)
		return
	}
}
