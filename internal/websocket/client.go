package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (c *Client) ReadPump() {
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()
    
    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    
    for {
        _, message, err := c.Conn.ReadMessage() // Ahora capturamos el 'message'
        if err != nil {
            break
        }
        
        // Creamos la estructura del mensaje para el Hub
        var msg Message
        if err := json.Unmarshal(message, &msg); err == nil {
            // Si el JSON es válido, lo mandamos al Hub para que lo reparta
            c.Hub.Broadcast <- msg
        } else {
            // Si envías texto plano como "Hola", lo envolvemos en un tipo genérico
            payload, _ := json.Marshal(map[string]string{"text": string(message)})
            c.Hub.Broadcast <- Message{
                Type:    "CHAT",
                RoomID:  c.RoomID,
                Payload: payload,
            }
        }
    }
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}