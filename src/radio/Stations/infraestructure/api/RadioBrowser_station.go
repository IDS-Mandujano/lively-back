package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"lively-backend/src/radio/Stations/domain/entity"
	"lively-backend/src/radio/Stations/domain/repository"
)

// Aseguramos en tiempo de compilación que RadioBrowserAPI implementa IStationRepository
var _ repository.IStationRepository = (*RadioBrowserAPI)(nil)

type RadioBrowserAPI struct {
	client    *http.Client
	userAgent string
}

func NewRadioBrowserAPI() *RadioBrowserAPI {
	return &RadioBrowserAPI{
		client: &http.Client{
			Timeout: 10 * time.Second, // Buena práctica: siempre pon un timeout
		},
		userAgent: "LivelyApp/1.0 (lively-backend)", // Requisito de la documentación
	}
}

// getBaseURL resuelve los DNS para encontrar un servidor disponible y elige uno al azar
func (api *RadioBrowserAPI) getBaseURL() (string, error) {
	// La forma más eficiente en Go de obtener los servidores según la doc: buscar el registro SRV
	_, srvs, err := net.LookupSRV("api", "tcp", "radio-browser.info")
	if err != nil || len(srvs) == 0 {
		return "", fmt.Errorf("error resolviendo DNS de radio-browser: %w", err)
	}

	// Aleatorizamos para no saturar un solo servidor
	rand.Seed(time.Now().UnixNano())
	randomServer := srvs[rand.Intn(len(srvs))]

	// Quitamos el punto final del target (ej: "de1.api.radio-browser.info.")
	host := strings.TrimSuffix(randomServer.Target, ".")
	return fmt.Sprintf("https://%s", host), nil
}

// Estructura interna SOLO para parsear el JSON ruidoso de la API externa
type radioBrowserResponse struct {
	StationUUID string `json:"stationuuid"`
	Name        string `json:"name"`
	URLResolved string `json:"url_resolved"`
	Favicon     string `json:"favicon"`
	CountryCode string `json:"countrycode"`
	Tags        string `json:"tags"`
}

// GetByCategory implementa el método de nuestra interfaz
func (api *RadioBrowserAPI) GetByCategory(ctx context.Context, category string, limit int) ([]entity.Station, error) {
	baseURL, err := api.getBaseURL()
	if err != nil {
		return nil, err
	}

	// Construimos la URL del endpoint
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

	// Mapeo: Convertimos la respuesta cruda a nuestra Entidad Limpia
	var cleanStations []entity.Station
	for _, raw := range rawStations {
		// Separamos el string de tags ("pop,rock,dance") en un arreglo real de Go
		tagsArray := []string{}
		if raw.Tags != "" {
			tagsArray = strings.Split(raw.Tags, ",")
		}

		cleanStations = append(cleanStations, entity.Station{
			ID:          raw.StationUUID,
			Name:        strings.TrimSpace(raw.Name),
			StreamURL:   raw.URLResolved, // Usamos url_resolved porque es el stream real
			ImageURL:    raw.Favicon,
			CountryCode: raw.CountryCode,
			Tags:        tagsArray,
		})
	}

	return cleanStations, nil
}

// Placeholder para los otros métodos de la interfaz
func (api *RadioBrowserAPI) SearchByName(ctx context.Context, name string, limit int) ([]entity.Station, error) {
	return nil, nil // Lo implementaremos después si lo necesitas
}

func (api *RadioBrowserAPI) GetTop(ctx context.Context, limit int) ([]entity.Station, error) {
	return nil, nil // Lo implementaremos después si lo necesitas
}