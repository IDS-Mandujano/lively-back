package models

type DeezerTrack struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Preview   string `json:"preview"`
	Artist    struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album     struct {
		CoverMedium string `json:"cover_medium"`
	} `json:"album"`
}

type DeezerSearchResponse struct {
	Data []DeezerTrack `json:"data"`
}