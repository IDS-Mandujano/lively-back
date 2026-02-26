package handler

import (
	"net/http"
	"strconv"
	"lively-backend/internal/websocket"

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
	Hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{Hub: hub}
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
}