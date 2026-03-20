package routes

import (
	"lively-backend/src/core/database"
	usecases "lively-backend/src/radio/Rooms/application/useCases"
	"lively-backend/src/radio/Rooms/infraestructure/controllers"
	"lively-backend/src/radio/Rooms/infraestructure/mysql"
	"lively-backend/src/radio/Rooms/infraestructure/sockets"
	"net/http"
)

func SetupRoomRoutes(mux *http.ServeMux) {
	// --- BASE DE DATOS (REST API) ---
	roomRepo := mysql.NewMySQLRoomRepository(database.DB)

	// Ruta para crear (POST)
	createRoomUC := usecases.NewCreateRoomUseCase(roomRepo)
	createRoomCtrl := controllers.NewCreateRoomController(createRoomUC)
	mux.HandleFunc("/api/rooms", createRoomCtrl.Handle)

	// Ruta para listar (GET) <--- NUEVO
	getAllRoomsUC := usecases.NewGetAllRoomsUseCase(roomRepo)
	getAllRoomsCtrl := controllers.NewGetAllRoomsController(getAllRoomsUC)
	mux.HandleFunc("/api/rooms/list", getAllRoomsCtrl.Handle)

	// --- TIEMPO REAL (WEBSOCKETS) ---
	roomManager := sockets.NewManager()
	wsController := sockets.NewWsController(roomManager)
	mux.HandleFunc("/api/ws/rooms", wsController.HandleConnections)
}
