package service

import (
	"lively-backend/internal/repository"
	"lively-backend/internal/websocket"
	"encoding/json"
)

type RoomService interface {
	SyncTrack(roomID uint, trackID int) error
}

type roomService struct {
	repo repository.RoomRepository
	hub  *websocket.Hub
}

func NewRoomService(repo repository.RoomRepository, hub *websocket.Hub) RoomService {
	return &roomService{repo: repo, hub: hub}
}

func (s *roomService) SyncTrack(roomID uint, trackID int) error {
	// 1. Persistencia (SSOT): Guardamos en MySQL primero
	err := s.repo.UpdateRoomTrack(roomID, trackID)
	if err != nil {
		return err
	}

	// 2. Tiempo Real: Avisamos a todos en la sala a través del Hub
	payload, _ := json.Marshal(map[string]interface{}{
		"track_id": trackID,
	})

	s.hub.Broadcast <- websocket.Message{
		Type:    "UPDATE_TRACK",
		RoomID:  roomID,
		Payload: payload,
	}

	return nil
}