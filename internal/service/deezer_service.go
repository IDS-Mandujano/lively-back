package service

import (
	"encoding/json"
	"fmt"
	"lively-backend/internal/models"
	"net/http"
	"os"
)

type DeezerService interface {
	SearchTracks(query string) ([]models.DeezerTrack, error)
	GetArtistDetails(artistID string) (map[string]interface{}, error)
	GetGenreRadios() (map[string]interface{}, error)
	GetTrackByID(trackID int) (*models.DeezerTrack, error)
}

type deezerService struct {
	baseURL string
}

func NewDeezerService() DeezerService {

	baseURL := os.Getenv("DEEZER_BASE_URL")

	if baseURL == "" {
		baseURL = "https://api.deezer.com"
	}

	return &deezerService{
		baseURL: baseURL,
	}
}

func (s *deezerService) SearchTracks(query string) ([]models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/search?q=%s", s.baseURL, query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []models.DeezerTrack `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Data, nil
}

func (s *deezerService) GetArtistDetails(artistID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/artist/%s", s.baseURL, artistID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func (s *deezerService) GetGenreRadios() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/radio", s.baseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func (s *deezerService) GetTrackByID(trackID int) (*models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/track/%d", s.baseURL, trackID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var track models.DeezerTrack
	if err := json.NewDecoder(resp.Body).Decode(&track); err != nil {
		return nil, err
	}
	return &track, nil
}
