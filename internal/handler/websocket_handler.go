package handler

import (
	"encoding/json"
	"fmt"
	"lively-backend/internal/repository"
	"lively-backend/internal/service"
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
	Svc  service.RoomService
}

func NewWebSocketHandler(hub *websocket.Hub, repo repository.RoomRepository, svc service.RoomService) *WebSocketHandler {
	return &WebSocketHandler{Hub: hub, Repo: repo, Svc: svc}
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
		// Bootstrap inicial: si no hay estado, podemos recibir query ?type=radio|artist&deezer_ref=123
		roomType := c.Query("type")
		deezerRef, _ := strconv.Atoi(c.DefaultQuery("deezer_ref", "0"))
		_ = h.Repo.UpdateRoomMeta(client.RoomID, fmt.Sprintf("room-%d", client.RoomID), roomType, deezerRef)
		_ = h.Svc.BootstrapRoom(client.RoomID, roomType, deezerRef)
	}
}
