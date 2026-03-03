package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	// Ajusta la importación según la ruta exacta de tu proyecto
	"lively-backend/src/radio/Stations/application/useCases"

)

// Estructura del controlador que inyecta el caso de uso
type GetStationsByCategoryController struct {
	useCase *usecases.GetStationsByCategoryUseCase
}

// Constructor
func NewGetStationsByCategoryController(uc *usecases.GetStationsByCategoryUseCase) *GetStationsByCategoryController {
	return &GetStationsByCategoryController{
		useCase: uc,
	}
}

// Handle es la función que responderá a la petición HTTP
func (c *GetStationsByCategoryController) Handle(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer parámetros de la URL (ejemplo: /stations/category?tag=rock&limit=15)
	category := r.URL.Query().Get("tag")
	if category == "" {
		// Podemos definir un comportamiento por defecto si el frontend no manda nada
		category = "pop" 
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // Límite por defecto
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 2. Ejecutar el Caso de Uso (La regla de negocio)
	// Le pasamos el Contexto de la petición HTTP por si el usuario cancela la carga en el frontend
	stations, err := c.useCase.Execute(r.Context(), category, limit)
	if err != nil {
		// Si algo falla en la API externa o en nuestro código, devolvemos un error 500
		http.Error(w, "Error obteniendo las estaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Responder al Frontend en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Convertimos nuestra entidad limpia de Go a JSON y la enviamos
	if err := json.NewEncoder(w).Encode(stations); err != nil {
		http.Error(w, "Error formateando la respuesta", http.StatusInternalServerError)
		return
	}
}