package websocket

import (
	"encoding/json"
	"log"

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
	Rooms      map[uint]map[*Client]bool 
	Broadcast  chan Message     
	Register   chan *Client      
	Unregister chan *Client          
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uint]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
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
			
			for client := range h.Rooms[msg.RoomID] {
				payload, _ := json.Marshal(msg)
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Rooms[msg.RoomID], client)
				}
			}
		}
	}
}