package sockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
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
	// 1. Obtenemos el ID de la sala
	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		http.Error(w, "Falta el ID de la sala", http.StatusBadRequest)
		return
	}

	// --- LA BARRERA DE SEGURIDAD (EL CADENERO) ---
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "No autorizado. Falta el token de acceso", http.StatusUnauthorized)
		return
	}

	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validamos que el método de encriptación sea el correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado")
		}
		return []byte(secretKey), nil
	})

	// Si el token no sirve, caducó, o es inventado, le cerramos la puerta en la cara
	if err != nil || !token.Valid {
		http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
		return
	}

	// Opcional: Extraemos los datos del usuario por si queremos saber quién entró
	var username string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = claims["username"].(string)
	}
	// ---------------------------------------------

	// 2. Si pasó el filtro, transformamos la conexión a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al crear el WebSocket:", err)
		return
	}

	// 3. El Manager busca la sala o crea una nueva
	room := c.manager.GetOrCreateRoom(roomID)

	// 4. Creamos al cliente
	client := &Client{
		Conn: ws,
		Room: room,
		Send: make(chan []byte, 256),
	}

	// 5. Metemos al cliente en la sala
	room.mu.Lock()
	room.Clients[client] = true

	// --- Sincronización Inicial ---
	var currentStationMsg []byte
	if room.CurrentStation != nil {
		currentStationMsg, _ = json.Marshal(map[string]interface{}{
			"type":    "station_changed",
			"payload": room.CurrentStation,
		})
	}
	room.mu.Unlock() // IMPORTANTE: Solo se abre el candado una vez

	// Si había un mensaje de estación actual, lo metemos al buzón de este cliente
	if currentStationMsg != nil {
		client.Send <- currentStationMsg
	}
	// -------------------------------------

	log.Printf("¡Usuario [%s] conectado a la sala: %s!", username, roomID)

	// 6. Arrancamos las orejas y la boca en segundo plano
	go client.ReadPump()
	go client.WritePump()
}
