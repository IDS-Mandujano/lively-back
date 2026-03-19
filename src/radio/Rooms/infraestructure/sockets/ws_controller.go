package sockets

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader transforma una petición HTTP normal en un túnel WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permite conexiones desde cualquier origen (necesario para Android y Postman)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsController struct {
	manager *Manager
}

func NewWsController(m *Manager) *WsController {
	return &WsController{
		manager: m,
	}
}

func (c *WsController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Obtenemos el ID de la sala de la URL (ej. ?id=sala-rock)
	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		http.Error(w, "Falta el ID de la sala (?id=lo-que-sea)", http.StatusBadRequest)
		return
	}

	// 2. Transformamos la conexión a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al crear el WebSocket:", err)
		return
	}

	// 3. El Manager busca la sala o crea una nueva
	room := c.manager.GetOrCreateRoom(roomID)

	// 4. Creamos al cliente (el usuario que acaba de entrar)
	client := &Client{
		Conn: ws,
		Room: room,
		Send: make(chan []byte, 256),
	}

	// 5. Metemos al cliente en la sala
	room.mu.Lock()
	room.Clients[client] = true

	// --- NUEVO: Sincronización Inicial ---
	// Si la sala ya tiene una estación reproduciéndose, se la mandamos directo a este usuario
	var currentStationMsg []byte
	if room.CurrentStation != nil {
		currentStationMsg, _ = json.Marshal(map[string]interface{}{
			"type":    "station_changed", // Reutilizamos este evento para que Android le dé Play
			"payload": room.CurrentStation,
		})
	}
	room.mu.Unlock()

	// Si había un mensaje de estación actual, lo metemos al buzón de este cliente
	if currentStationMsg != nil {
		client.Send <- currentStationMsg
	}
	// -------------------------------------

	log.Printf("¡Nuevo usuario conectado a la sala: %s!", roomID)

	// 6. Arrancamos las orejas y la boca en segundo plano
	go client.ReadPump()
	go client.WritePump()
}
