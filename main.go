package main

import (
	"fmt"
	"log"
	"net/http"
	// Les ponemos un alias a los imports para diferenciarlos
	"lively-backend/src/core/database"
	roomRoutes "lively-backend/src/radio/Rooms/infraestructure/routes"
	stationRoutes "lively-backend/src/radio/Stations/infraestructure/routes"
	userRoutes "lively-backend/src/users/infraestructure/routes"

)

func main() {

	database.Connect()
	// Creamos un nuevo enrutador HTTP estándar de Go
	mux := http.NewServeMux()

	// Llamamos a nuestra función para que registre las rutas de "Stations" (REST API)
	stationRoutes.SetupStationRoutes(mux)
	roomRoutes.SetupRoomRoutes(mux)
	userRoutes.SetupUserRoutes(mux) // <--- NUEVO

	// Definimos el puerto donde correrá el backend
	port := ":8080"
	fmt.Printf("Servidor de Lively Backend corriendo en http://localhost%s\n", port)
	fmt.Printf("Prueba el endpoint en: http://localhost%s/api/stations/category?tag=rock\n", port)
	fmt.Printf("Conexión WebSocket lista en: ws://localhost%s/api/ws/rooms\n", port)

	// Arrancamos el servidor
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
