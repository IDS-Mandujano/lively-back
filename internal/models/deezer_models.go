package models

type DeezerTrack struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Readable bool   `json:"readable"`
	Preview  string `json:"preview"`
	Duration int    `json:"duration"`
	Artist   struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		PictureMedium string `json:"picture_medium"`
	} `json:"artist"`
	Album struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		CoverMedium string `json:"cover_medium"`
	} `json:"album"`
}

type DeezerSearchResponse struct {
	Data []DeezerTrack `json:"data"`
}
