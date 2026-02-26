package handler

import (
	"net/http"
	"lively-backend/internal/service"

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
		RoomID  uint `json:"room_id" binding:"required"`
		TrackID int  `json:"track_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de sala o canción inválidos"})
		return
	}

	// Llamamos al servicio que actualiza DB y emite al WebSocket
	err := h.roomService.SyncTrack(input.RoomID, input.TrackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al sincronizar la sala"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Sincronización enviada"})
}