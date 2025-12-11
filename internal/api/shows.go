package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Show API Methods
// -----------------------------------------------------------------------------

// GetShow retrieves a single show by ID.
// API: GET /v2/shows/{show_id}
func (c *Client) GetShow(showID int) (*models.Show, error) {
	path := fmt.Sprintf("/shows/%d", showID)

	var resp models.ShowResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Show, nil
}

// CreateShowParams contains parameters for creating a new show.
type CreateShowParams struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
	Language    string `json:"language,omitempty"`
	Explicit    bool   `json:"explicit,omitempty"`
}

// CreateShow creates a new podcast show.
// Requires authentication.
// API: POST /v2/shows
func (c *Client) CreateShow(params CreateShowParams) (*models.Show, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("authentication required")
	}

	// Convert to form fields (Spreaker uses form data, not JSON)
	fields := map[string]string{
		"title": params.Title,
	}
	if params.Description != "" {
		fields["description"] = params.Description
	}
	if params.CategoryID > 0 {
		fields["category_id"] = fmt.Sprintf("%d", params.CategoryID)
	}
	if params.Language != "" {
		fields["language"] = params.Language
	}
	if params.Explicit {
		fields["explicit"] = "true"
	}

	var resp models.ShowResponse
	if err := c.PostForm("/shows", fields, &resp); err != nil {
		return nil, err
	}

	return &resp.Show, nil
}

// UpdateShowParams contains parameters for updating a show.
type UpdateShowParams struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	CategoryID  *int    `json:"category_id,omitempty"`
	Language    *string `json:"language,omitempty"`
	Explicit    *bool   `json:"explicit,omitempty"`
}

// UpdateShow updates an existing show.
// Requires authentication and ownership.
// API: POST /v2/shows/{show_id}
func (c *Client) UpdateShow(showID int, params UpdateShowParams) (*models.Show, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/shows/%d", showID)

	// Build form fields only for non-nil parameters
	fields := make(map[string]string)
	if params.Title != nil {
		fields["title"] = *params.Title
	}
	if params.Description != nil {
		fields["description"] = *params.Description
	}
	if params.CategoryID != nil {
		fields["category_id"] = fmt.Sprintf("%d", *params.CategoryID)
	}
	if params.Language != nil {
		fields["language"] = *params.Language
	}
	if params.Explicit != nil {
		if *params.Explicit {
			fields["explicit"] = "true"
		} else {
			fields["explicit"] = "false"
		}
	}

	var resp models.ShowResponse
	if err := c.PostForm(path, fields, &resp); err != nil {
		return nil, err
	}

	return &resp.Show, nil
}

// DeleteShow deletes a show.
// Requires authentication and ownership.
// API: DELETE /v2/shows/{show_id}
func (c *Client) DeleteShow(showID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/shows/%d", showID)
	return c.Delete(path, nil)
}

// GetShowEpisodes retrieves all episodes of a show.
// API: GET /v2/shows/{show_id}/episodes
func (c *Client) GetShowEpisodes(showID int, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/shows/%d/episodes", showID)
	return GetPaginated[models.Episode](c, path, pagination.ToMap())
}

// AddShowToFavorites adds a show to the user's favorites.
// Requires authentication.
// API: PUT /v2/users/{user_id}/favorites/{show_id}
func (c *Client) AddShowToFavorites(userID, showID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/favorites/%d", userID, showID)
	return c.Put(path, nil)
}

// RemoveShowFromFavorites removes a show from the user's favorites.
// Requires authentication.
// API: DELETE /v2/users/{user_id}/favorites/{show_id}
func (c *Client) RemoveShowFromFavorites(userID, showID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/favorites/%d", userID, showID)
	return c.Delete(path, nil)
}

// GetFavoriteShows retrieves the user's favorite shows.
// API: GET /v2/users/{user_id}/favorites
func (c *Client) GetFavoriteShows(userID int, pagination PaginationParams) (*PaginatedResult[models.Show], error) {
	path := fmt.Sprintf("/users/%d/favorites", userID)
	return GetPaginated[models.Show](c, path, pagination.ToMap())
}