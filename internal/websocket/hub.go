package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	RoomID  uint            `json:"room_id"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	ID     uint
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID uint
}

type Hub struct {
	Rooms       map[uint]map[*Client]bool
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	LastMessage map[uint]Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:       make(map[uint]map[*Client]bool),
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		LastMessage: make(map[uint]Message),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			log.Printf("Usuario %d se unió a la sala %d", client.ID, client.RoomID)

		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID][client]; ok {
				delete(h.Rooms[client.RoomID], client)
				close(client.Send)
			}

		case msg := <-h.Broadcast:
			// Guardamos el último mensaje conocido por sala para nuevas conexiones
			h.LastMessage[msg.RoomID] = msg

			for client := range h.Rooms[msg.RoomID] {
				payload, _ := json.Marshal(msg)
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Rooms[msg.RoomID], client)
				}
			}
		case <-ticker.C:
			// Periodic resync: reenviamos el último mensaje a cada sala para corregir drift
			for roomID, last := range h.LastMessage {
				for client := range h.Rooms[roomID] {
					payload, _ := json.Marshal(last)
					select {
					case client.Send <- payload:
					default:
						close(client.Send)
						delete(h.Rooms[roomID], client)
					}
				}
			}
		}
	}
}
