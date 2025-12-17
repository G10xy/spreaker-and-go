package models


type Chapter struct {
	ChapterID int `json:"chapter_id"`

	StartsAt int `json:"starts_at"`

	Title string `json:"title"`

	ExternalURL string `json:"external_url,omitempty"`

	ImageURL string `json:"image_url,omitempty"`

	ImageOriginalURL string `json:"image_original_url,omitempty"`
}


type ChapterResponse struct {
	Chapter Chapter `json:"chapter"`
}
