package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"lively-backend/src/radio/Stations/application/useCases"

)

type GetStationsByCategoryController struct {
	useCase *usecases.GetStationsByCategoryUseCase
}

func NewGetStationsByCategoryController(uc *usecases.GetStationsByCategoryUseCase) *GetStationsByCategoryController {
	return &GetStationsByCategoryController{
		useCase: uc,
	}
}

func (c *GetStationsByCategoryController) Handle(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("tag")
	if category == "" {
		category = "pop" 
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	stations, err := c.useCase.Execute(r.Context(), category, limit)
	if err != nil {
	
		http.Error(w, "Error obteniendo las estaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(stations); err != nil {
		http.Error(w, "Error formateando la respuesta", http.StatusInternalServerError)
		return
	}
}