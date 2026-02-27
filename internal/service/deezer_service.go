package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lively-backend/internal/models"
	"log"
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
		log.Printf("Deezer SearchTracks request error: %v url=%s", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("Deezer SearchTracks non-200 status: %d url=%s body=%s", resp.StatusCode, url, string(body))
		return nil, fmt.Errorf("deezer search status %d", resp.StatusCode)
	}

	var result struct {
		Data []models.DeezerTrack `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Deezer SearchTracks decode error: %v body=%s", err, string(body))
		return nil, err
	}
	return result.Data, nil
}

func (s *deezerService) GetArtistDetails(artistID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/artist/%s", s.baseURL, artistID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Deezer GetArtistDetails request error: %v url=%s", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("Deezer GetArtistDetails non-200 status: %d url=%s body=%s", resp.StatusCode, url, string(body))
		return nil, fmt.Errorf("deezer artist status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Deezer GetArtistDetails decode error: %v body=%s", err, string(body))
		return nil, err
	}
	return result, nil
}

func (s *deezerService) GetGenreRadios() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/radio", s.baseURL)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Deezer GetGenreRadios request error: %v url=%s", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("Deezer GetGenreRadios non-200 status: %d url=%s body=%s", resp.StatusCode, url, string(body))
		return nil, fmt.Errorf("deezer radio status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Deezer GetGenreRadios decode error: %v body=%s", err, string(body))
		return nil, err
	}
	return result, nil
}

func (s *deezerService) GetTrackByID(trackID int) (*models.DeezerTrack, error) {
	url := fmt.Sprintf("%s/track/%d", s.baseURL, trackID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Deezer GetTrackByID request error: %v url=%s", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("Deezer GetTrackByID non-200 status: %d url=%s body=%s", resp.StatusCode, url, string(body))
		return nil, fmt.Errorf("deezer track status %d", resp.StatusCode)
	}

	var track models.DeezerTrack
	if err := json.Unmarshal(body, &track); err != nil {
		log.Printf("Deezer GetTrackByID decode error: %v body=%s", err, string(body))
		return nil, err
	}
	return &track, nil
}
