package models

import (
	"fmt"
)

type Episode struct {
	EpisodeID int `json:"episode_id"`

	Title string `json:"title"`

	Description string `json:"description"`

	ShowID int `json:"show_id"`

	Show *Show `json:"show,omitempty"`

	Author *User `json:"author,omitempty"`

	AuthorID int `json:"author_id"`
	SiteURL string `json:"site_url"`

	ImageURL string `json:"image_url"`

	ImageOriginalURL string `json:"image_original_url"`

	Duration int `json:"duration"`

	PlayCount int `json:"plays_count"`

	LikesCount int `json:"likes_count"`

	MessagesCount int `json:"messages_count"`

	DownloadEnabled bool `json:"download_enabled"`

	Explicit bool `json:"explicit"`

	Hidden bool `json:"hidden"`

	Tags []string `json:"tags,omitempty"`

	PublishedAt *CustomTime `json:"published_at,omitempty"`

	EncodingStatus string `json:"encoding_status"`

	MediaURL string `json:"media_url,omitempty"`

	DownloadURL string `json:"download_url,omitempty"`
}

type EpisodeResponse struct {
	Episode Episode `json:"episode"`
}

// DurationFormatted returns the episode duration as a human-readable string.
// The duration is stored in milliseconds.
func (e *Episode) DurationFormatted() string {
	totalSeconds := e.Duration / 1000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
