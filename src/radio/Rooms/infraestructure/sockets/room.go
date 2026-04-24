package sockets

import (
	"encoding/json"
	"lively-backend/src/radio/Stations/domain/entity"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Room *Room
	Send chan []byte
}

type Room struct {
	ID             string
	Clients        map[*Client]bool
	CurrentStation *entity.Station
	Broadcast      chan []byte
	mu             sync.Mutex

	VoteActive      bool
	VotesInFavor    int
	VotesAgainst    int
	VotedClients    map[*Client]bool
	ProposedStation *entity.Station
}

func (r *Room) finalizeElectionIfNeededLocked() (bool, string, []byte) {
	if !r.VoteActive {
		return false, "", nil
	}

	totalVotes := r.VotesInFavor + r.VotesAgainst
	totalClients := len(r.Clients)

	// Si ya no hay nadie conectado, cerramos la votación para evitar sala bloqueada.
	if totalClients == 0 {
		r.VoteActive = false
		r.ProposedStation = nil
		r.VotedClients = make(map[*Client]bool)
		return false, "", nil
	}

	if totalVotes < totalClients {
		return false, "", nil
	}

	approved := r.VotesInFavor > r.VotesAgainst
	r.VoteActive = false

	var resultType string
	var resultPayload interface{}
	if approved && r.ProposedStation != nil {
		r.CurrentStation = r.ProposedStation
		resultType = "station_changed"
		resultPayload = r.CurrentStation
	} else {
		resultType = "vote_failed"
		resultPayload = map[string]string{"reason": "No alcanzó mayoría o fue un empate"}
	}

	notification, _ := json.Marshal(map[string]interface{}{
		"type":    resultType,
		"payload": resultPayload,
	})
	return true, resultType, notification
}

func NewRoom(id string) *Room {
	return &Room{
		ID:           id,
		Clients:      make(map[*Client]bool),
		Broadcast:    make(chan []byte),
		VotedClients: make(map[*Client]bool), // Inicializamos el mapa de votantes
		CurrentStation: &entity.Station{
			ID:        "default-lofi",
			Name:      "Lofi Girl Radio",
			StreamURL: "https://play.streamafrica.net/lofiradio",
		},
	}
}

// Run es el motor de la sala: se queda escuchando infinitamente
func (r *Room) Run() {
	for {
		// Cuando alguien manda un mensaje al canal Broadcast...
		message := <-r.Broadcast

		r.mu.Lock()
		// ...recorremos a TODOS los clientes conectados a esta sala
		for client := range r.Clients {
			// Y les enviamos el mensaje a su buzón personal
			select {
			case client.Send <- message:
			default:
				// Si el buzón del cliente está lleno o se desconectó a la fuerza, lo borramos
				close(client.Send)
				delete(r.Clients, client)
			}
		}
		r.mu.Unlock()
	}
}

// ReadPump se queda escuchando todo lo que el celular Android nos mande
func (c *Client) ReadPump() {
	defer func() {
		c.Room.mu.Lock()
		delete(c.Room.Clients, c)
		electionFinished, _, notification := c.Room.finalizeElectionIfNeededLocked()
		c.Room.mu.Unlock()

		if electionFinished && notification != nil {
			c.Room.Broadcast <- notification
		}

		c.Conn.Close()
		log.Println("Un usuario se desconectó de la sala:", c.Room.ID)
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 1. Abrimos el "sobre" para ver de qué tipo es el mensaje
		var wsMsg WsMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Println("Error leyendo JSON del cliente:", err)
			continue
		}

		// 2. Tomamos decisiones dependiendo del tipo de mensaje
		switch wsMsg.Type {
		case "propose":
			c.Room.mu.Lock()
			// Si ya hay votación, ignoramos la nueva propuesta
			if c.Room.VoteActive {
				c.Room.mu.Unlock()
				continue
			}

			// Preparamos la urna de votación
			c.Room.VoteActive = true
			c.Room.VotesInFavor = 0
			c.Room.VotesAgainst = 0
			c.Room.VotedClients = make(map[*Client]bool)

			// Leemos qué estación propusieron
			var proposePayload ProposePayload
			json.Unmarshal(wsMsg.Payload, &proposePayload)

			c.Room.ProposedStation = &entity.Station{
				ID:        proposePayload.StationID,
				Name:      proposePayload.StationName,
				StreamURL: proposePayload.StreamURL,
			}
			c.Room.mu.Unlock()

			// Avisamos a toda la sala que empezó una votación
			notification, _ := json.Marshal(map[string]interface{}{
				"type":    "vote_started",
				"payload": proposePayload,
			})
			c.Room.Broadcast <- notification
			log.Printf("Votación iniciada en %s para la estación: %s", c.Room.ID, proposePayload.StationName)

		case "vote":
			c.Room.mu.Lock()
			// Si no hay votación activa o si este usuario ya votó, ignoramos
			if !c.Room.VoteActive || c.Room.VotedClients[c] {
				c.Room.mu.Unlock()
				continue
			}

			// Leemos si votó "Sí" (true) o "No" (false)
			var votePayload VotePayload
			json.Unmarshal(wsMsg.Payload, &votePayload)

			c.Room.VotedClients[c] = true // Lo marcamos para que no vote doble
			if votePayload.Vote {
				c.Room.VotesInFavor++
			} else {
				c.Room.VotesAgainst++
			}

			votesInFavor := c.Room.VotesInFavor
			votesAgainst := c.Room.VotesAgainst

			electionFinished, resultType, notification := c.Room.finalizeElectionIfNeededLocked()
			c.Room.mu.Unlock()

			// Enviamos el conteo parcial para que Android actualice la UI en tiempo real.
			voteUpdate, _ := json.Marshal(map[string]interface{}{
				"type": "vote_updated",
				"payload": map[string]int{
					"in_favor": votesInFavor,
					"against":  votesAgainst,
				},
			})
			c.Room.Broadcast <- voteUpdate

			// Si la elección terminó, le gritamos el resultado a todos
			if electionFinished {
				if resultType == "station_changed" {
					log.Printf("¡Votación APROBADA en %s! Cambiando de estación...", c.Room.ID)
				} else {
					log.Printf("Votación RECHAZADA en %s.", c.Room.ID)
				}
				if notification != nil {
					c.Room.Broadcast <- notification
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
