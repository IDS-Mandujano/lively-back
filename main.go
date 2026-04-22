package main

import (
	"fmt"
	"lively-backend/src/core/database"
	roomRoutes "lively-backend/src/radio/Rooms/infraestructure/routes"
	stationRoutes "lively-backend/src/radio/Stations/infraestructure/routes"
	userRoutes "lively-backend/src/users/infraestructure/routes"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[boot] iniciando lively-backend")

	database.Connect()
	mux := http.NewServeMux()

	stationRoutes.SetupStationRoutes(mux)
	roomRoutes.SetupRoomRoutes(mux)
	userRoutes.SetupUserRoutes(mux)

	port := ":8080"
	fmt.Printf("Servidor de Lively Backend corriendo en http://localhost%s\n", port)
	fmt.Printf("Prueba el endpoint en: http://localhost%s/api/stations/category?tag=rock\n", port)
	fmt.Printf("Conexión WebSocket lista en: ws://localhost%s/api/ws/rooms\n", port)
	log.Printf("[boot] escuchando en %s", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
