package repository

import (
	"fmt"
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
			return r.db.Create(&room).Error
		}
		return err
	}
	return r.db.Model(&room).Update("current_track_id", trackID).Error
}

func (r *roomRepository) GetRoomByID(roomID uint) (*models.Room, error) {
	var room models.Room
	if err := r.db.Preload("Members").First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
