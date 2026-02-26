package repository

import (
	"lively-backend/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	UpdateRoomTrack(roomID uint, trackID int) error
	GetRoomByID(roomID uint) (*models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) UpdateRoomTrack(roomID uint, trackID int) error {
	// Actualiza solo el campo current_track_id en la tabla rooms
	return r.db.Model(&models.Room{}).Where("id = ?", roomID).Update("current_track_id", trackID).Error
}

func (r *roomRepository) GetRoomByID(roomID uint) (*models.Room, error) {
	var room models.Room
	if err := r.db.Preload("Members").First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}