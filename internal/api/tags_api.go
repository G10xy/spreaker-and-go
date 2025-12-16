package api

import (
	"fmt"
	"net/url"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Tags API
// -----------------------------------------------------------------------------

// GetEpisodesByTag retrieves the latest episodes with a specific tag.
// API: GET /v2/tags/{tag_name}/episodes
// Parameters:
//   - tagName: The tag name to search for (can contain spaces, will be URL encoded)
//   - pagination: Pagination parameters
func (c *Client) GetEpisodesByTag(tagName string, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	// URL encode the tag name to handle spaces and special characters
	encodedTag := url.PathEscape(tagName)
	path := fmt.Sprintf("/tags/%s/episodes", encodedTag)

	return GetPaginated[models.Episode](c, path, pagination.ToMap())
}
