package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// SearchParams contains parameters for search queries.
type SearchParams struct {
	Query  string
	Filter string // "listenable" (default) or "editable"
}

func (p SearchParams) ToMap() map[string]string {
	params := make(map[string]string)
	if p.Query != "" {
		params["q"] = p.Query
	}
	if p.Filter != "" {
		params["filter"] = p.Filter
	}
	return params
}

// -----------------------------------------------------------------------------
// Search Shows
// -----------------------------------------------------------------------------

// SearchShows searches for shows matching the query.
// API: GET /v2/search?type=shows&q={query}
func (c *Client) SearchShows(search SearchParams, pagination PaginationParams) (*PaginatedResult[models.Show], error) {
	path := "/search"

	queryParams := search.ToMap()
	queryParams["type"] = "shows"
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.Show](c, path, queryParams)
}

// SearchUserShows searches for shows by a specific user.
// API: GET /v2/search/users/{user_id}?type=shows&q={query}
func (c *Client) SearchUserShows(userID int, search SearchParams, pagination PaginationParams) (*PaginatedResult[models.Show], error) {
	path := fmt.Sprintf("/search/users/%d", userID)

	queryParams := search.ToMap()
	queryParams["type"] = "shows"
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.Show](c, path, queryParams)
}

// -----------------------------------------------------------------------------
// Search Episodes
// -----------------------------------------------------------------------------

// SearchEpisodes searches for episodes matching the query.
// API: GET /v2/search?type=episodes&q={query}
func (c *Client) SearchEpisodes(search SearchParams, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := "/search"

	queryParams := search.ToMap()
	queryParams["type"] = "episodes"
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.Episode](c, path, queryParams)
}

// SearchUserEpisodes searches for episodes by a specific user.
// API: GET /v2/search/users/{user_id}?type=episodes&q={query}
func (c *Client) SearchUserEpisodes(userID int, search SearchParams, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/search/users/%d", userID)

	queryParams := search.ToMap()
	queryParams["type"] = "episodes"
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.Episode](c, path, queryParams)
}

// SearchShowEpisodes searches for episodes within a specific show.
// API: GET /v2/search/shows/{show_id}?type=episodes&q={query}
func (c *Client) SearchShowEpisodes(showID int, search SearchParams, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/search/shows/%d", showID)

	queryParams := search.ToMap()
	queryParams["type"] = "episodes"
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.Episode](c, path, queryParams)
}
