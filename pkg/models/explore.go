package models

// ExploreShow represents a show item returned by explore endpoints.
// Used by: GET /v2/explore/categories/{category_id}/items
type ExploreShow struct {
	ShowID           int    `json:"show_id"`
	Title            string `json:"title"`
	SiteURL          string `json:"site_url"`
	ImageURL         string `json:"image_url"`
	ImageOriginalURL string `json:"image_original_url"`
	AuthorID         int    `json:"author_id"`
}
