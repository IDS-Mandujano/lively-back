package repository

import (
	"fmt"
	"lively-backend/internal/models"

	"gorm.io/gorm"
)

type RoomRepository interface {
	UpdateRoomTrack(roomID uint, trackID int) error
	GetRoomByID(roomID uint) (*models.Room, error)
	UpdateRoomMeta(roomID uint, name string, roomType string, deezerRef int) error
	PushHistory(roomID uint, trackID int) error
	PopLastHistory(roomID uint) (int, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) UpdateRoomTrack(roomID uint, trackID int) error {
	// Intentamos obtener la sala; si no existe la creamos y seteamos current_track_id
	var room models.Room
	if err := r.db.First(&room, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Crear sala con valores por defecto para cumplir constraints (Name, Type)
			room = models.Room{
				ID:             roomID,
				Name:           fmt.Sprintf("room-%d", roomID),
				Type:           "auto",
				CurrentTrackID: trackID,
			}
			if err := r.db.Create(&room).Error; err != nil {
				return err
			}
			return nil
		}
		return err
	}
	room.PreviousTrackID = room.CurrentTrackID
	if room.PreviousTrackID != 0 {
		_ = r.PushHistory(roomID, room.PreviousTrackID)
	}
	room.CurrentTrackID = trackID
	return r.db.Model(&room).Updates(map[string]interface{}{
		"previous_track_id": room.PreviousTrackID,
		"current_track_id":  room.CurrentTrackID,
	}).Error
}

func (r *roomRepository) GetRoomByID(roomID uint) (*models.Room, error) {
	var room models.Room
	if err := r.db.Preload("Members").First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) UpdateRoomMeta(roomID uint, name string, roomType string, deezerRef int) error {
	var room models.Room
	if err := r.db.First(&room, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			room = models.Room{
				ID:                roomID,
				Name:              name,
				Type:              roomType,
				DeezerReferenceID: deezerRef,
			}
			return r.db.Create(&room).Error
		}
		return err
	}
	room.Name = name
	room.Type = roomType
	room.DeezerReferenceID = deezerRef
	return r.db.Save(&room).Error
}

func (r *roomRepository) PushHistory(roomID uint, trackID int) error {
	h := models.RoomHistory{RoomID: roomID, TrackID: trackID}
	if err := r.db.Create(&h).Error; err != nil {
		return err
	}
	// Mantener últimas 5
	var ids []uint
	if err := r.db.
		Model(&models.RoomHistory{}).
		Where("room_id = ?", roomID).
		Order("created_at desc").
		Offset(5).
		Pluck("id", &ids).Error; err == nil && len(ids) > 0 {
		_ = r.db.Delete(&models.RoomHistory{}, ids).Error
	}
	return nil
}

func (r *roomRepository) PopLastHistory(roomID uint) (int, error) {
	var h models.RoomHistory
	if err := r.db.
		Where("room_id = ?", roomID).
		Order("created_at desc").
		First(&h).Error; err != nil {
		return 0, err
	}
	_ = r.db.Delete(&h).Error
	return h.TrackID, nil
}
