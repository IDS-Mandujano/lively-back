package entity

import "time"

// User representa a un usuario registrado en Lively
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // El guión oculta este campo al convertirlo a JSON
	CreatedAt    time.Time `json:"created_at"`
}