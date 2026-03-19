package sockets

import "encoding/json"

// WsMessage es el "sobre" estándar. Todo lo que envíe Android debe tener este formato.
type WsMessage struct {
	Type    string          `json:"type"`    // Ej: "propose", "vote", "chat"
	Payload json.RawMessage `json:"payload"` // El contenido real (se procesa dependiendo del Type)
}

// ProposePayload es lo que viene dentro del Payload cuando alguien propone cambiar de estación
type ProposePayload struct {
	StationID   string `json:"station_id"`
	StationName string `json:"station_name"`
	StreamURL   string `json:"stream_url"`
}

// VotePayload es lo que viene dentro del Payload cuando un usuario presiona "Sí" o "No"
type VotePayload struct {
	Vote bool `json:"vote"` // true = Sí quiero cambiar, false = No quiero cambiar
}