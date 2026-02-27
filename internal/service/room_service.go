package service

import (
	"encoding/json"
	"lively-backend/internal/repository"
	"lively-backend/internal/websocket"
	"time"
)

type RoomService interface {
	SyncTrack(roomID uint, trackID int) error
}

type roomService struct {
	repo   repository.RoomRepository
	hub    *websocket.Hub
	deezer DeezerService
}

func NewRoomService(repo repository.RoomRepository, hub *websocket.Hub, deezer DeezerService) RoomService {
	return &roomService{repo: repo, hub: hub, deezer: deezer}
}

func (s *roomService) SyncTrack(roomID uint, trackID int) error {
	// 1. Persistencia (SSOT): Guardamos en MySQL primero
	err := s.repo.UpdateRoomTrack(roomID, trackID)
	if err != nil {
		return err
	}

	// 2. Construimos payload enriquecido: track_id, timestamp de inicio y referencia Deezer
	startTs := time.Now().UnixMilli()
	deezerRef := 0
	trackURL := ""
	if room, err := s.repo.GetRoomByID(roomID); err == nil && room != nil {
		deezerRef = room.DeezerReferenceID
	}

	// Intentamos obtener la URL/preview y metadatos desde Deezer
	title := ""
	cover := ""
	duration := 0
	if t, err := s.deezer.GetTrackByID(trackID); err == nil && t != nil {
		trackURL = t.Preview
		title = t.Title
		cover = t.Album.CoverMedium
		duration = t.Duration
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"track_id":            trackID,
		"start_timestamp_ms":  startTs,
		"deezer_reference_id": deezerRef,
		"track_url":           trackURL,
		"title":               title,
		"cover_medium":        cover,
		"duration":            duration,
	})

	s.hub.Broadcast <- websocket.Message{
		Type:    "UPDATE_TRACK",
		RoomID:  roomID,
		Payload: payload,
	}

	return nil
}
