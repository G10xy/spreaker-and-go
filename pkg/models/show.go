package models


type Show struct {
	ShowID int `json:"show_id"`

	Title string `json:"title"`

	Description string `json:"description"`

	SiteURL string `json:"site_url"`

	ImageURL string `json:"image_url"`

	ImageOriginalURL string `json:"image_original_url"`

	Author *User `json:"author,omitempty"`

	AuthorID int `json:"author_id"`

	Category *Category `json:"category,omitempty"`

	CategoryID int `json:"category_id"`

	Language string `json:"language"`

	EpisodesCount int `json:"episodes_count"`

	FollowersCount int `json:"followers_count"`

	PlayCount int `json:"plays_count"`

	LikesCount int `json:"likes_count"`

	Explicit bool `json:"explicit"`

	LastEpisodeAt *CustomTime `json:"last_episode_at,omitempty"`

	CreatedAt *CustomTime `json:"created_at,omitempty"`
}

type ShowResponse struct {
	Show Show `json:"show"`
}

// Category represents a show category.
type Category struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}
