package handler

import (
	"lively-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(rs service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: rs}
}

func (h *RoomHandler) SyncTrack(c *gin.Context) {
	var input struct {
		RoomID      uint   `json:"room_id" binding:"required"`
		TrackID     int    `json:"track_id"`
		Type        string `json:"type"`          // "radio" | "artist" | "auto"
		DeezerRefID int    `json:"deezer_ref_id"` // id de radio o artista
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de sala o canción inválidos"})
		return
	}

	var err error
	if input.TrackID > 0 {
		err = h.roomService.SyncTrack(input.RoomID, input.TrackID)
	} else {
		err = h.roomService.BootstrapRoom(input.RoomID, input.Type, input.DeezerRefID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al sincronizar la sala"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Sincronización enviada"})
}

func (h *RoomHandler) PrevTrack(c *gin.Context) {
	var input struct {
		RoomID uint `json:"room_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de sala inválidos"})
		return
	}
	if err := h.roomService.PrevTrack(input.RoomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al retroceder pista"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Retroceso enviado"})
}
