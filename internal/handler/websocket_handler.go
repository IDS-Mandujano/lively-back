package handler

import (
	"encoding/json"
	"lively-backend/internal/repository"
	"lively-backend/internal/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

// Upgrader necesario para convertir la petición HTTP a WebSocket
var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// En desarrollo, permitimos cualquier origen para evitar errores de CORS con Android
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	Hub  *websocket.Hub
	Repo repository.RoomRepository
}

func NewWebSocketHandler(hub *websocket.Hub, repo repository.RoomRepository) *WebSocketHandler {
	return &WebSocketHandler{Hub: hub, Repo: repo}
}

func (h *WebSocketHandler) HandleWS(c *gin.Context) {
	// Obtenemos los parámetros de la URL: /ws/:roomID/:userID
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	userID, _ := strconv.Atoi(c.Param("userID"))

	// Hacemos el "Upgrade" de la conexión
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Creamos el nuevo cliente para el Hub
	client := &websocket.Client{
		ID:     uint(userID),
		Hub:    h.Hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		RoomID: uint(roomID),
	}

	// Lo registramos en el Hub
	h.Hub.Register <- client

	// Iniciamos los hilos de lectura y escritura para este cliente
	go client.WritePump()
	go client.ReadPump()

	// Si el Hub tiene un último mensaje para la sala, enviarlo al cliente recién conectado
	if msg, ok := h.Hub.LastMessage[client.RoomID]; ok {
		if data, err := json.Marshal(msg); err == nil {
			select {
			case client.Send <- data:
			default:
			}
		}
	} else {
		// Fallback: si no hay mensaje en memoria, intentar leer DB y enviar solo track_id
		if room, err := h.Repo.GetRoomByID(client.RoomID); err == nil && room != nil {
			payload, _ := json.Marshal(map[string]interface{}{"track_id": room.CurrentTrackID})
			msg := websocket.Message{
				Type:    "UPDATE_TRACK",
				RoomID:  client.RoomID,
				Payload: payload,
			}
			if client.Send != nil {
				if data, err := json.Marshal(msg); err == nil {
					select {
					case client.Send <- data:
					default:
					}
				}
			}
		}
	}
}
