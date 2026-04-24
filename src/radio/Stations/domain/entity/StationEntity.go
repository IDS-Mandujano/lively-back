package entity

type Station struct {
	ID          string   `json:"id"`           
	Name        string   `json:"name"`
	StreamURL   string   `json:"stream_url"`
	ImageURL    string   `json:"image_url"`
	CountryCode string   `json:"country_code"`
	Tags        []string `json:"tags"`
	StationUUID string   `json:"stationuuid,omitempty"`
	URLResolved string   `json:"url_resolved,omitempty"`
	Favicon     string   `json:"favicon,omitempty"`
}