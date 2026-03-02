package service

import (
	"encoding/json"
	"fmt"
	"lively-backend/internal/models"
	"lively-backend/internal/repository"
	"lively-backend/internal/websocket"
	"log"
	"time"
)

type RoomService interface {
	SyncTrack(roomID uint, trackID int) error
	PrevTrack(roomID uint) error
	BootstrapRoom(roomID uint, roomType string, deezerRef int) error
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
	err := s.repo.UpdateRoomTrack(roomID, trackID)
	if err != nil {
		return err
	}

	startTs := time.Now().UnixMilli()
	var t *models.DeezerTrack
	radioIDUsed := 0

	isValid := func(tr *models.DeezerTrack) bool {
		return tr != nil && tr.Preview != "" && tr.Readable
	}

	t, _ = s.deezer.GetTrackByID(trackID)

	if !isValid(t) {
		radioIDUsed = trackID
		log.Printf("RoomService: ID %d not playable. Trying as Radio...", trackID)

		radioTracks, err := s.deezer.GetRadioTracks(trackID, 50, 0)
		if err == nil && len(radioTracks) > 0 {
			for _, rt := range radioTracks {
				if rt.Preview != "" {
					t = &rt
					log.Printf("RoomService: Found playable track in radio: %s", t.Title)
					break
				}
			}
		}
	}

	if !isValid(t) {
		log.Printf("RoomService: CRITICAL - No audio found for room %d", roomID)
		return nil
	}

	_ = s.repo.UpdateRoomTrack(roomID, t.ID)

	payload, _ := json.Marshal(map[string]interface{}{
		"track_id":           t.ID,
		"id":                 fmt.Sprintf("%d", t.ID),
		"start_timestamp_ms": startTs,
		"radio_id_used":      radioIDUsed,
		"track_url":          t.Preview,
		"preview":            t.Preview,
		"title":              t.Title,
		"artist": map[string]string{
			"name":           t.Artist.Name,
			"picture_medium": t.Artist.PictureMedium,
		},
		"album": map[string]string{
			"cover_medium": t.Album.CoverMedium,
		},
		"cover_medium":         t.Album.CoverMedium,
		"duration":             t.Duration,
		"requires_deezer_auth": false,
	})

	s.hub.Broadcast <- websocket.Message{
		Type:    "UPDATE_TRACK",
		RoomID:  roomID,
		Payload: payload,
	}

	return nil
}

func (s *roomService) PrevTrack(roomID uint) error {
	last, err := s.repo.PopLastHistory(roomID)
	if err != nil || last == 0 {
		return nil
	}
	err = s.repo.UpdateRoomTrack(roomID, last)
	if err != nil {
		return err
	}
	t, err := s.deezer.GetTrackByID(last)
	if err != nil || t == nil || t.Preview == "" || !t.Readable {
		return nil
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"track_id":           t.ID,
		"id":                 fmt.Sprintf("%d", t.ID),
		"start_timestamp_ms": time.Now().UnixMilli(),
		"track_url":          t.Preview,
		"preview":            t.Preview,
		"title":              t.Title,
		"artist": map[string]string{
			"name":           t.Artist.Name,
			"picture_medium": t.Artist.PictureMedium,
		},
		"album": map[string]string{
			"cover_medium": t.Album.CoverMedium,
		},
		"cover_medium":         t.Album.CoverMedium,
		"duration":             t.Duration,
		"requires_deezer_auth": false,
	})
	s.hub.Broadcast <- websocket.Message{
		Type:    "UPDATE_TRACK",
		RoomID:  roomID,
		Payload: payload,
	}
	return nil
}

func (s *roomService) BootstrapRoom(roomID uint, roomType string, deezerRef int) error {
	_ = s.repo.UpdateRoomMeta(roomID, fmt.Sprintf("room-%d", roomID), roomType, deezerRef)
	var chosen *models.DeezerTrack
	isValid := func(tr *models.DeezerTrack) bool {
		return tr != nil && tr.Preview != "" && tr.Readable
	}
	if roomType == "radio" && deezerRef > 0 {
		if tracks, err := s.deezer.GetRadioTracks(deezerRef, 50, 0); err == nil {
			for _, rt := range tracks {
				if rt.Preview != "" && rt.Readable {
					chosen = &rt
					break
				}
			}
		}
	}
	if !isValid(chosen) && deezerRef > 0 {
		if tracks, err := s.deezer.GetArtistTopTracks(fmt.Sprintf("%d", deezerRef), 50); err == nil {
			for _, rt := range tracks {
				if rt.Preview != "" && rt.Readable {
					chosen = &rt
					break
				}
			}
		}
	}
	if !isValid(chosen) {
		return nil
	}
	_ = s.repo.UpdateRoomTrack(roomID, chosen.ID)
	payload, _ := json.Marshal(map[string]interface{}{
		"track_id":           chosen.ID,
		"id":                 fmt.Sprintf("%d", chosen.ID),
		"start_timestamp_ms": time.Now().UnixMilli(),
		"track_url":          chosen.Preview,
		"preview":            chosen.Preview,
		"title":              chosen.Title,
		"artist": map[string]string{
			"name":           chosen.Artist.Name,
			"picture_medium": chosen.Artist.PictureMedium,
		},
		"album": map[string]string{
			"cover_medium": chosen.Album.CoverMedium,
		},
		"cover_medium":         chosen.Album.CoverMedium,
		"duration":             chosen.Duration,
		"requires_deezer_auth": false,
	})
	s.hub.Broadcast <- websocket.Message{
		Type:    "UPDATE_TRACK",
		RoomID:  roomID,
		Payload: payload,
	}
	return nil
}
