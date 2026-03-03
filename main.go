package main

import (
	"fmt"
	"log"
	"net/http"
	"lively-backend/src/radio/Stations/infraestructure/routes"

)

func main() {
	// Creamos un nuevo enrutador HTTP estándar de Go
	mux := http.NewServeMux()

	// Llamamos a nuestra función para que registre las rutas de "Stations"
	routes.SetupStationRoutes(mux)

	// Definimos el puerto donde correrá el backend
	port := ":8080"
	fmt.Printf("Servidor de Lively Backend corriendo en http://localhost%s\n", port)
	fmt.Printf("Prueba el endpoint en: http://localhost%s/api/stations/category?tag=rock\n", port)

	// Arrancamos el servidor
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}