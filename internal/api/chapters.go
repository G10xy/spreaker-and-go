package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Episode Chapters API
// -----------------------------------------------------------------------------

// GetEpisodeChapters retrieves all chapters for an episode.
// API: GET /v2/episodes/{episode_id}/chapters
func (c *Client) GetEpisodeChapters(episodeID int, pagination PaginationParams) (*PaginatedResult[models.Chapter], error) {
	path := fmt.Sprintf("/episodes/%d/chapters", episodeID)
	return GetPaginated[models.Chapter](c, path, pagination.ToMap())
}

// ChapterParams contains the parameters for creating or updating a chapter.
type ChapterParams struct {
	StartsAt *int

	Title string

	ExternalURL string

	ImageFile string

	ImageCrop string
}

func (p ChapterParams) ToMap() map[string]string {
	params := make(map[string]string)

	if p.StartsAt != nil {
		params["starts_at"] = fmt.Sprintf("%d", *p.StartsAt)
	}
	if p.Title != "" {
		params["title"] = p.Title
	}
	if p.ExternalURL != "" {
		params["external_url"] = p.ExternalURL
	}
	if p.ImageFile != "" {
		params["image_file"] = p.ImageFile
	}
	if p.ImageCrop != "" {
		params["image_crop"] = p.ImageCrop
	}

	return params
}

// AddChapter adds a new chapter to an episode.
// API: POST /v2/episodes/{episode_id}/chapters
func (c *Client) AddChapter(episodeID int, params ChapterParams) (*models.Chapter, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	if params.StartsAt == nil {
		return nil, fmt.Errorf("starts_at is required")
	}
	if params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	path := fmt.Sprintf("/episodes/%d/chapters", episodeID)

	var resp models.ChapterResponse
	if err := c.Post(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Chapter, nil
}

// UpdateChapter updates an existing chapter.
// API: POST /v2/episodes/{episode_id}/chapters/{chapter_id}
func (c *Client) UpdateChapter(episodeID, chapterID int, params ChapterParams) (*models.Chapter, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/chapters/%d", episodeID, chapterID)

	var resp models.ChapterResponse
	if err := c.Post(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Chapter, nil
}

// DeleteChapter deletes a single chapter from an episode.
// API: DELETE /v2/episodes/{episode_id}/chapters/{chapter_id}
func (c *Client) DeleteChapter(episodeID, chapterID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/chapters/%d", episodeID, chapterID)
	return c.Delete(path, nil)
}

// DeleteAllChapters deletes all chapters from an episode.
// API: DELETE /v2/episodes/{episode_id}/chapters
func (c *Client) DeleteAllChapters(episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/chapters", episodeID)
	return c.Delete(path, nil)
}
