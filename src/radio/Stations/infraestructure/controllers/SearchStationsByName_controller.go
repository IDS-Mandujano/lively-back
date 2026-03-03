package controllers

import (
	"encoding/json"
	usecases "lively-backend/src/radio/Stations/application/useCases"
	"net/http"
	"strconv"
)

type SearchStationsByNameController struct {
	useCase *usecases.SearchStationsByNameUseCase
}

func NewSearchStationsByNameController(uc *usecases.SearchStationsByNameUseCase) *SearchStationsByNameController {
	return &SearchStationsByNameController{
		useCase: uc,
	}
}

func (c *SearchStationsByNameController) Handle(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Falta el parámetro de búsqueda 'name'", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	stations, err := c.useCase.Execute(r.Context(), name, limit)
	if err != nil {
		http.Error(w, "Error buscando estaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stations); err != nil {
		http.Error(w, "Error formateando la respuesta", http.StatusInternalServerError)
		return
	}
}
