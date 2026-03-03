package routes

import (
	"net/http"
	// Ajusta siempre estos imports al nombre de tu módulo
	"lively-backend/src/radio/Stations/application/useCases"
	"lively-backend/src/radio/Stations/infraestructure/api"
	"lively-backend/src/radio/Stations/infraestructure/controllers"

)

// SetupStationRoutes inicializa las rutas y hace la inyección de dependencias
func SetupStationRoutes(mux *http.ServeMux) {
	
	// 1. Instanciamos el Adaptador de Infraestructura (La conexión a la API externa)
	radioAPI := api.NewRadioBrowserAPI()

	// 2. Instanciamos el Caso de Uso inyectándole el adaptador (que cumple con la interfaz IStationRepository)
	getStationsByCategoryUC := usecases.NewGetStationsByCategoryUseCase(radioAPI)

	// 3. Instanciamos el Controlador inyectándole el Caso de Uso
	getStationsByCategoryCtrl := controllers.NewGetStationsByCategoryController(getStationsByCategoryUC)

	// 4. Registramos la ruta en el servidor HTTP
	// Cuando el frontend llame a esta URL, se ejecutará la función Handle de nuestro controlador
	mux.HandleFunc("/api/stations/category", getStationsByCategoryCtrl.Handle)
	
	// Aquí en el futuro agregaremos:
	// mux.HandleFunc("/api/stations/search", searchStationsByNameCtrl.Handle)
	// mux.HandleFunc("/api/stations/top", getTopStationsCtrl.Handle)
}