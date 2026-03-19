package sockets

import (
	"sync"
)

// Manager es el gerente general de todas las salas de Lively
type Manager struct {
	Rooms map[string]*Room // Un diccionario con todas las salas activas (ID -> Sala)
	mu    sync.Mutex       // Candado para evitar que dos personas creen la misma sala al mismo milisegundo
}

// NewManager crea un gerente vacío al arrancar el servidor
func NewManager() *Manager {
	return &Manager{
		Rooms: make(map[string]*Room),
	}
}

// GetOrCreateRoom busca una sala. Si no existe, la construye y arranca su motor.
func (m *Manager) GetOrCreateRoom(roomID string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Si la sala ya existe, simplemente la devolvemos
	if room, exists := m.Rooms[roomID]; exists {
		return room
	}

	// Si no existe, la creamos usando el constructor que hiciste
	newRoom := NewRoom(roomID)
	m.Rooms[roomID] = newRoom

	// ¡SÚPER IMPORTANTE! Arrancamos el motor (goroutine) de la sala en segundo plano
	// Si no ponemos el "go" aquí, el servidor se quedaría trabado para siempre
	go newRoom.Run()

	return newRoom
}
