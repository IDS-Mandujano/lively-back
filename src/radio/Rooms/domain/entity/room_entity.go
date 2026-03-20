package entity

import "time"

// Room representa los datos de una sala guardada permanentemente en MySQL
type Room struct {
	ID          string    `json:"id"`          // Ej: "sala-rock-123"
	Name        string    `json:"name"`        // Ej: "La Cueva del Rock"
	Description string    `json:"description"` // Ej: "Puro rock de los 80s"
	CreatedBy   int       `json:"created_by"`  // El ID del usuario que la creó
	CreatedAt   time.Time `json:"created_at"`
}