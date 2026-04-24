package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ repository.IStationRepository = (*RadioBrowserAPI)(nil)

type RadioBrowserAPI struct {
	client    *http.Client
	userAgent string
}

func NewRadioBrowserAPI() *RadioBrowserAPI {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &RadioBrowserAPI{
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: customTransport,
		},
		userAgent: "LivelyApp/1.0 (lively-backend)",
	}
}

func (api *RadioBrowserAPI) getBaseURL() (string, error) {
	_, srvs, err := net.LookupSRV("api", "tcp", "radio-browser.info")
	if err != nil || len(srvs) == 0 {
		return "", fmt.Errorf("error resolviendo DNS de radio-browser: %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	randomServer := srvs[rand.Intn(len(srvs))]

	host := strings.TrimSuffix(randomServer.Target, ".")
	return fmt.Sprintf("https://%s", host), nil
}

type radioBrowserResponse struct {
	StationUUID string `json:"stationuuid"`
	Name        string `json:"name"`
	URLResolved string `json:"url_resolved"`
	Favicon     string `json:"favicon"`
	CountryCode string `json:"countrycode"`
	Tags        string `json:"tags"`
}

func (api *RadioBrowserAPI) GetByCategory(ctx context.Context, category string, limit int) ([]entity.Station, error) {
	baseURL, err := api.getBaseURL()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/json/stations/bytag/%s?limit=%d&hidebroken=true", baseURL, url.PathEscape(category), limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", api.userAgent)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API externa devolvió status: %d", res.StatusCode)
	}

	var rawStations []radioBrowserResponse
	if err := json.NewDecoder(res.Body).Decode(&rawStations); err != nil {
		return nil, err
	}

	var cleanStations []entity.Station
	for _, raw := range rawStations {
		tagsArray := []string{}
		if raw.Tags != "" {
			tagsArray = strings.Split(raw.Tags, ",")
		}

		cleanStations = append(cleanStations, entity.Station{
			ID:          raw.StationUUID,
			Name:        strings.TrimSpace(raw.Name),
			StreamURL:   raw.URLResolved,
			ImageURL:    raw.Favicon,
			CountryCode: raw.CountryCode,
			Tags:        tagsArray,
			StationUUID: raw.StationUUID,
			URLResolved: raw.URLResolved,
			Favicon:     raw.Favicon,
		})
	}

	return cleanStations, nil
}

func (api *RadioBrowserAPI) SearchByName(ctx context.Context, name string, limit int) ([]entity.Station, error) {
	baseURL, err := api.getBaseURL()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/json/stations/search?name=%s&limit=%d&hidebroken=true&order=clickcount&reverse=true",
		baseURL, url.QueryEscape(name), limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.userAgent)
	req.Header.Set("Accept", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API externa devolvió status: %d", res.StatusCode)
	}

	var rawStations []radioBrowserResponse
	if err := json.NewDecoder(res.Body).Decode(&rawStations); err != nil {
		return nil, err
	}

	var cleanStations []entity.Station
	for _, raw := range rawStations {
		tagsArray := []string{}
		if raw.Tags != "" {
			tagsArray = strings.Split(raw.Tags, ",")
		}

		cleanStations = append(cleanStations, entity.Station{
			ID:          raw.StationUUID,
			Name:        strings.TrimSpace(raw.Name),
			StreamURL:   raw.URLResolved,
			ImageURL:    raw.Favicon,
			CountryCode: raw.CountryCode,
			Tags:        tagsArray,
			StationUUID: raw.StationUUID,
			URLResolved: raw.URLResolved,
			Favicon:     raw.Favicon,
		})
	}

	return cleanStations, nil
}

func (api *RadioBrowserAPI) GetTop(ctx context.Context, limit int) ([]entity.Station, error) {
	baseURL, err := api.getBaseURL()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/json/stations/topclick/%d?hidebroken=true", baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "LivelyApp/1.0 (MandujanoDev)")

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API externa devolvió status: %d", res.StatusCode)
	}

	var rawStations []radioBrowserResponse
	if err := json.NewDecoder(res.Body).Decode(&rawStations); err != nil {
		return nil, err
	}

	var cleanStations []entity.Station
	for _, raw := range rawStations {
		tagsArray := []string{}
		if raw.Tags != "" {
			tagsArray = strings.Split(raw.Tags, ",")
		}

		cleanStations = append(cleanStations, entity.Station{
			ID:          raw.StationUUID,
			Name:        strings.TrimSpace(raw.Name),
			StreamURL:   raw.URLResolved,
			ImageURL:    raw.Favicon,
			CountryCode: raw.CountryCode,
			Tags:        tagsArray,
			StationUUID: raw.StationUUID,
			URLResolved: raw.URLResolved,
			Favicon:     raw.Favicon,
		})
	}

	return cleanStations, nil
}
