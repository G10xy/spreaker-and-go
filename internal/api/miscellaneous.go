package api

import (
	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Miscellaneous API - Categories
// -----------------------------------------------------------------------------

// GetShowCategories retrieves all available show categories.
// API: GET /v2/show-categories
// Parameters:
//   - locale: Optional locale for category names (e.g., "it_IT" for Italian)
func (c *Client) GetShowCategories(locale string) ([]models.Category, error) {
	params := make(map[string]string)
	if locale != "" {
		params["c"] = locale
	}

	var resp models.CategoriesResponse
	if err := c.Get("/show-categories", params, &resp); err != nil {
		return nil, err
	}

	return resp.Categories, nil
}

// GetGooglePlayCategories retrieves all available Google Play podcast categories.
// API: GET /v2/googleplay-categories
func (c *Client) GetGooglePlayCategories() ([]models.GooglePlayCategory, error) {
	var resp models.GooglePlayCategoriesResponse
	if err := c.Get("/googleplay-categories", nil, &resp); err != nil {
		return nil, err
	}

	return resp.GooglePlayCategories, nil
}

// -----------------------------------------------------------------------------
// Miscellaneous API - Languages
// -----------------------------------------------------------------------------

// GetShowLanguages retrieves all available show languages.
// API: GET /v2/show-languages
// Parameters:
//   - locale: Optional locale for language names (e.g., "it_IT" for Italian)
func (c *Client) GetShowLanguages(locale string) (map[string]string, error) {
	params := make(map[string]string)
	if locale != "" {
		params["c"] = locale
	}

	var resp models.LanguagesResponse
	if err := c.Get("/show-languages", params, &resp); err != nil {
		return nil, err
	}

	return resp.Languages, nil
}

// GetShowLanguagesList is a convenience method that returns languages as a slice
func (c *Client) GetShowLanguagesList(locale string) ([]models.Language, error) {
	langMap, err := c.GetShowLanguages(locale)
	if err != nil {
		return nil, err
	}

	languages := make([]models.Language, 0, len(langMap))
	for code, name := range langMap {
		languages = append(languages, models.Language{
			Code: code,
			Name: name,
		})
	}

	return languages, nil
}
