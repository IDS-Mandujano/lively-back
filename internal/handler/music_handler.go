package handler

import (
	"lively-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	deezerService service.DeezerService
}

func NewMusicHandler(ds service.DeezerService) *MusicHandler {
	return &MusicHandler{deezerService: ds}
}

func (h *MusicHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Falta el parámetro de búsqueda"})
		return
	}

	tracks, err := h.deezerService.SearchTracks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// En MusicHandler añade:

func (h *MusicHandler) GetArtist(c *gin.Context) {
	id := c.Param("id")
	artist, err := h.deezerService.GetArtistDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, artist)
}

func (h *MusicHandler) GetRadios(c *gin.Context) {
	radios, err := h.deezerService.GetGenreRadios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, radios)
}

func (h *MusicHandler) GetTrack(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "track id inválido"})
		return
	}
	track, err := h.deezerService.GetTrackByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, track)
}
