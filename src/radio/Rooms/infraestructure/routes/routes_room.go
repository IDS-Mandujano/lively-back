package routes

import (
	"lively-backend/src/radio/Rooms/infraestructure/sockets"
	"net/http"
)

// SetupRoomRoutes configura las rutas de WebSockets para las salas
func SetupRoomRoutes(mux *http.ServeMux) {
	// 1. Creamos al gerente de salas general (solo nace una vez)
	roomManager := sockets.NewManager()

	// 2. Creamos el controlador inyectándole al gerente
	wsController := sockets.NewWsController(roomManager)

	// 3. Abrimos la puerta al público
	mux.HandleFunc("/api/ws/rooms", wsController.HandleConnections)
}
