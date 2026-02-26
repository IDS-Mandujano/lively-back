package models

import "time"

type Room struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Type              string    `gorm:"not null" json:"type"` 
	DeezerReferenceID int       `json:"deezer_reference_id"`  
	CurrentTrackID    int       `json:"current_track_id"`
	CreatedAt         time.Time `json:"created_at"`

	Members []User `gorm:"many2many:room_members;" json:"members,omitempty"`
}