package service

import (
	"encoding/json"
	"fmt"
	"io"
	"lively-backend/internal/models"
	"net/http"
	"os"
	"time"
)

type DeezerService interface {
	SearchTracks(query string, limit int, index int) ([]models.DeezerTrack, error)
	GetArtistDetails(artistID string) (map[string]interface{}, error)
	GetGenreRadios() (map[string]interface{}, error)
	GetRadioTracks(radioID int, limit int, index int) ([]models.DeezerTrack, error)
	GetTrackByID(trackID int) (*models.DeezerTrack, error)
	GetArtistTopTracks(artistID string, limit int) ([]models.DeezerTrack, error)
}

type deezerService struct {
	baseURL string
	client  *http.Client
}

func NewDeezerService() DeezerService {
	baseURL := os.Getenv("DEEZER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deezer.com"
	}
	return &deezerService{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 8 * time.Second},
	}
}

func (s *deezerService) SearchTracks(query string, limit int, index int) ([]models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/search?q=%s", s.baseURL, query)
	if limit > 0 {
		url = fmt.Sprintf("%s&limit=%d", url, limit)
	}
	if index > 0 {
		url = fmt.Sprintf("%s&index=%d", url, index)
	}
	return s.fetchTracks(url)
}

func (s *deezerService) GetRadioTracks(radioID int, limit int, index int) ([]models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/radio/%d/tracks", s.baseURL, radioID)
	if limit > 0 {
		url = fmt.Sprintf("%s?limit=%d", url, limit)
		if index > 0 {
			url = fmt.Sprintf("%s&index=%d", url, index)
		}
	} else if index > 0 {
		url = fmt.Sprintf("%s?index=%d", url, index)
	}
	return s.fetchTracks(url)
}

func (s *deezerService) fetchTracks(url string) ([]models.DeezerTrack, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deezer %d: %s", resp.StatusCode, string(b))
	}
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data []models.DeezerTrack `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	// Filtrar pistas sin preview o no legibles
	filtered := make([]models.DeezerTrack, 0, len(result.Data))
	for _, t := range result.Data {
		if t.Preview != "" && t.Readable {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func (s *deezerService) GetTrackByID(trackID int) (*models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/track/%d", s.baseURL, trackID)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deezer %d: %s", resp.StatusCode, string(b))
	}
	body, _ := io.ReadAll(resp.Body)
	var track models.DeezerTrack
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, err
	}
	return &track, nil
}

func (s *deezerService) GetArtistDetails(artistID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/artist/%s", s.baseURL, artistID)
	return s.fetchMap(url)
}

func (s *deezerService) GetGenreRadios() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/radio", s.baseURL)
	return s.fetchMap(url)
}

func (s *deezerService) fetchMap(url string) (map[string]interface{}, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deezer %d: %s", resp.StatusCode, string(b))
	}
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *deezerService) GetArtistTopTracks(artistID string, limit int) ([]models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/artist/%s/top", s.baseURL, artistID)
	if limit > 0 {
		url = fmt.Sprintf("%s?limit=%d", url, limit)
	}
	return s.fetchTracks(url)
}
