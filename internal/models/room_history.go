package models

import "time"

type RoomHistory struct {
	ID        uint      `gorm:"primaryKey"`
	RoomID    uint      `gorm:"index"`
	TrackID   int
	CreatedAt time.Time
}
