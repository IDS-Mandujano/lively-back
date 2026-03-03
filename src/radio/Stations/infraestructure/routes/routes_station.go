package routes

import (
	"net/http"
	"lively-backend/src/radio/Stations/application/useCases"
	"lively-backend/src/radio/Stations/infraestructure/api"
	"lively-backend/src/radio/Stations/infraestructure/controllers"

)

func SetupStationRoutes(mux *http.ServeMux) {
	
	radioAPI := api.NewRadioBrowserAPI()

	getStationsByCategoryUC := usecases.NewGetStationsByCategoryUseCase(radioAPI)
	getTopStationsUC := usecases.NewGetTopStationsUseCase(radioAPI)
	searchStationsByNameUC := usecases.NewSearchStationsByNameUseCase(radioAPI)

	getStationsByCategoryCtrl := controllers.NewGetStationsByCategoryController(getStationsByCategoryUC)
	getTopStationsCtrl := controllers.NewGetTopStationsController(getTopStationsUC)
	searchStationsByNameCtrl := controllers.NewSearchStationsByNameController(searchStationsByNameUC)

	mux.HandleFunc("/api/stations/category", getStationsByCategoryCtrl.Handle)
	mux.HandleFunc("/api/stations/top", getTopStationsCtrl.Handle)
	mux.HandleFunc("/api/stations/search", searchStationsByNameCtrl.Handle)
	
}