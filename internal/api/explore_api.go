package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)


// GetCategoryShows retrieves shows in a specific category.
// API: GET /v2/explore/categories/{category_id}/items
// Parameters:
//   - categoryID: The category ID to explore (use GetShowCategories to list available categories)
//   - pagination: Pagination parameters
func (c *Client) GetCategoryShows(categoryID int, pagination PaginationParams) (*PaginatedResult[models.ExploreShow], error) {
	path := fmt.Sprintf("/explore/categories/%d/items", categoryID)
	return GetPaginated[models.ExploreShow](c, path, pagination.ToMap())
}
