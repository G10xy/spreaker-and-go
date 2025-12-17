package models

// -----------------------------------------------------------------------------
// Episode Message Models
// -----------------------------------------------------------------------------

type Message struct {
	MessageID int `json:"message_id"`

	EpisodeID int `json:"episode_id"`

	Text string `json:"text"`

	CreatedAt string `json:"created_at"`

	AuthorID int `json:"author_id"`

	AuthorUsername string `json:"author_username"`

	AuthorFullname string `json:"author_fullname"`

	AuthorSiteURL string `json:"author_site_url"`

	AuthorImageURL string `json:"author_image_url,omitempty"`

	AuthorImageOriginalURL string `json:"author_image_original_url,omitempty"`

	AuthorIsOwner bool `json:"author_is_owner"`

	AppName string `json:"app_name,omitempty"`

	AppURL string `json:"app_url,omitempty"`
}
