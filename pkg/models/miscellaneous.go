package models

// -----------------------------------------------------------------------------
// Category Models
// -----------------------------------------------------------------------------

type Category struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	Permalink  string `json:"permalink,omitempty"`
	Level      int    `json:"level"` // 1 = top-level, 2 = subcategory
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type GooglePlayCategory struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
}

type GooglePlayCategoriesResponse struct {
	GooglePlayCategories []GooglePlayCategory `json:"googleplay_categories"`
}

// -----------------------------------------------------------------------------
// Language Models
// -----------------------------------------------------------------------------

type LanguagesResponse struct {
	Languages map[string]string `json:"languages"`
}

type Language struct {
	Code string
	Name string
}
