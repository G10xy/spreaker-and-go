package api

import (
	"fmt"
	"net/http"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Episode API Methods
// -----------------------------------------------------------------------------

// GetEpisode retrieves a single episode by ID.
// API: GET /v2/episodes/{episode_id}
func (c *Client) GetEpisode(episodeID int) (*models.Episode, error) {
	path := fmt.Sprintf("/episodes/%d", episodeID)

	var resp models.EpisodeResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Episode, nil
}

type UploadEpisodeParams struct {
	// Required
	Title     string // Episode title
	MediaFile string // Path to the audio file

	// Optional
	Description     string   // Episode description/show notes
	Tags            []string // Tags for the episode
	Explicit        bool     // Contains explicit content
	DownloadEnabled bool     // Allow downloads
	Hidden          bool     // Hidden/private episode
	AutoPublishedAt string   // Schedule publishing (format: "2020-04-20 18:00:00")
}

// UploadEpisode uploads a new episode to a show.
// API: POST /v2/shows/{show_id}/episodes
func (c *Client) UploadEpisode(showID int, params UploadEpisodeParams) (*models.Episode, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	if params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if params.MediaFile == "" {
		return nil, fmt.Errorf("media_file is required")
	}

	path := fmt.Sprintf("/shows/%d/episodes", showID)

	fields := map[string]string{
		"title": params.Title,
	}

	if params.Description != "" {
		fields["description"] = params.Description
	}
	if len(params.Tags) > 0 {
		tagStr := ""
		for i, tag := range params.Tags {
			if i > 0 {
				tagStr += ","
			}
			tagStr += tag
		}
		fields["tags"] = tagStr
	}
	if params.Explicit {
		fields["explicit"] = "true"
	}
	if params.DownloadEnabled {
		fields["download_enabled"] = "true"
	}
	if params.Hidden {
		fields["hidden"] = "true"
	}
	if params.AutoPublishedAt != "" {
		fields["auto_published_at"] = params.AutoPublishedAt
	}

	var resp models.EpisodeResponse
	if err := c.PostFormWithFile(path, fields, "media_file", params.MediaFile, &resp); err != nil {
		return nil, err
	}

	return &resp.Episode, nil
}

// CreateDraftEpisodeParams contains parameters for creating a draft episode.
type CreateDraftEpisodeParams struct {
	// Required
	Title  string // Episode title
	ShowID int    // Show ID where the episode will belong

	// Optional
	Description     string   // Episode description/show notes
	Tags            []string // Tags for the episode
	Explicit        bool     // Contains explicit content
	DownloadEnabled bool     // Allow downloads
	Hidden          bool     // Hidden/private episode
}

// CreateDraftEpisode creates a new draft episode without an audio file.
// The audio file can be uploaded later using UpdateEpisode with a media_file.
// API: POST /v2/episodes/drafts
func (c *Client) CreateDraftEpisode(params CreateDraftEpisodeParams) (*models.Episode, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	if params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if params.ShowID == 0 {
		return nil, fmt.Errorf("show_id is required")
	}

	fields := map[string]string{
		"title":   params.Title,
		"show_id": fmt.Sprintf("%d", params.ShowID),
	}

	if params.Description != "" {
		fields["description"] = params.Description
	}
	if len(params.Tags) > 0 {
		tagStr := ""
		for i, tag := range params.Tags {
			if i > 0 {
				tagStr += ","
			}
			tagStr += tag
		}
		fields["tags"] = tagStr
	}
	if params.Explicit {
		fields["explicit"] = "true"
	}
	if params.DownloadEnabled {
		fields["download_enabled"] = "true"
	}
	if params.Hidden {
		fields["hidden"] = "true"
	}

	var resp models.EpisodeResponse
	if err := c.PostForm("/episodes/drafts", fields, &resp); err != nil {
		return nil, err
	}

	return &resp.Episode, nil
}

type UpdateEpisodeParams struct {
	Title           *string
	Description     *string
	Tags            *[]string
	Explicit        *bool
	DownloadEnabled *bool
	Hidden          *bool
	ShowID          *int    // Move episode to a different show
	AutoPublishedAt *string // Reschedule or unschedule (empty string to unschedule)
}

// UpdateEpisode updates an existing episode.
// API: POST /v2/episodes/{episode_id}
func (c *Client) UpdateEpisode(episodeID int, params UpdateEpisodeParams) (*models.Episode, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d", episodeID)

	// Build form fields only for non-nil parameters
	fields := make(map[string]string)

	if params.Title != nil {
		fields["title"] = *params.Title
	}
	if params.Description != nil {
		fields["description"] = *params.Description
	}
	if params.Tags != nil {
		tagStr := ""
		for i, tag := range *params.Tags {
			if i > 0 {
				tagStr += ","
			}
			tagStr += tag
		}
		fields["tags"] = tagStr
	}
	if params.Explicit != nil {
		if *params.Explicit {
			fields["explicit"] = "true"
		} else {
			fields["explicit"] = "false"
		}
	}
	if params.DownloadEnabled != nil {
		if *params.DownloadEnabled {
			fields["download_enabled"] = "true"
		} else {
			fields["download_enabled"] = "false"
		}
	}
	if params.Hidden != nil {
		if *params.Hidden {
			fields["hidden"] = "true"
		} else {
			fields["hidden"] = "false"
		}
	}
	if params.ShowID != nil {
		fields["show_id"] = fmt.Sprintf("%d", *params.ShowID)
	}
	if params.AutoPublishedAt != nil {
		fields["auto_published_at"] = *params.AutoPublishedAt
	}

	var resp models.EpisodeResponse
	if err := c.PostForm(path, fields, &resp); err != nil {
		return nil, err
	}

	return &resp.Episode, nil
}

// DeleteEpisode deletes an episode.
// API: DELETE /v2/episodes/{episode_id}
func (c *Client) DeleteEpisode(episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d", episodeID)
	return c.Delete(path, nil)
}

// LikeEpisode adds an episode to the user's likes.
// API: PUT /v2/users/{user_id}/likes/{episode_id}
func (c *Client) LikeEpisode(userID, episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%d/likes/%d", userID, episodeID)
	return c.Put(path, nil)
}

// UnlikeEpisode removes an episode from the user's likes.
// Requires authentication.
// API: DELETE /v2/users/{user_id}/likes/{episode_id}
func (c *Client) UnlikeEpisode(userID, episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%d/likes/%d", userID, episodeID)
	return c.Delete(path, nil)
}

// CheckUserLikesEpisode checks if a user has liked a specific episode.
// API: GET /v2/users/{user_id}/likes/{episode_id}
func (c *Client) CheckUserLikesEpisode(userID, episodeID int) (bool, error) {
	path := fmt.Sprintf("/users/%d/likes/%d", userID, episodeID)

	var resp models.EpisodeResponse
	if err := c.Get(path, nil, &resp); err != nil {
		// Check if it's a 404 (not liked)
		if apiErr, ok := err.(*APIError); ok && apiErr.IsNotFound() {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetEpisodeLikes retrieves the list of users who liked an episode.
// API: GET /v2/episodes/{episode_id}/likes
func (c *Client) GetEpisodeLikes(episodeID int, pagination PaginationParams) (*PaginatedResult[models.User], error) {
	path := fmt.Sprintf("/episodes/%d/likes", episodeID)
	return GetPaginated[models.User](c, path, pagination.ToMap())
}

// GetLikedEpisodes retrieves the user's liked episodes.
// API: GET /v2/users/{user_id}/likes
func (c *Client) GetLikedEpisodes(userID int, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/users/%d/likes", userID)
	return GetPaginated[models.Episode](c, path, pagination.ToMap())
}

// BookmarkEpisode adds an episode to the user's bookmarks.
// Note: You can only bookmark episodes on your own account, so userID must match
// the owner of the token used for authentication.
// API: PUT /v2/users/{user_id}/bookmarks/{episode_id}
func (c *Client) BookmarkEpisode(userID, episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%d/bookmarks/%d", userID, episodeID)
	return c.Put(path, nil)
}

// UnbookmarkEpisode removes an episode from the user's bookmarks.
// API: DELETE /v2/users/{user_id}/bookmarks/{episode_id}
func (c *Client) UnbookmarkEpisode(userID, episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%d/bookmarks/%d", userID, episodeID)
	return c.Delete(path, nil)
}

// GetEpisodeDownloadURL retrieves the download URL for an episode.
// API: GET /v2/episodes/{episode_id}/download
func (c *Client) GetEpisodeDownloadURL(episodeID int) (string, error) {
    path := fmt.Sprintf("/episodes/%d/download", episodeID)
    urlStr := c.buildURL(path)

    // Create a client that doesn't follow redirects
    noRedirectClient := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse 
        },
        Timeout: c.HTTPClient.Timeout,
    }

    req, err := c.newRequest(http.MethodGet, urlStr, nil)
    if err != nil {
        return "", err
    }

    resp, err := noRedirectClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 && resp.StatusCode < 400 {
        location := resp.Header.Get("Location")
        if location != "" {
            return location, nil
        }
        return "", fmt.Errorf("redirect response but no Location header")
    }

    if resp.StatusCode == http.StatusOK {
        return urlStr, nil
    }

    return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// GetEpisodePlayURL retrieves the streaming URL for an episode.
// API: GET /v2/episodes/{episode_id}/play
func (c *Client) GetEpisodePlayURL(episodeID int) (string, error) {
	path := fmt.Sprintf("/episodes/%d/play", episodeID)

	var resp struct {
		URL string `json:"url"`
	}
	if err := c.Get(path, nil, &resp); err != nil {
		return "", err
	}

	return resp.URL, nil
}

// GetUserEpisodes retrieves all episodes published by a user.
// API: GET /v2/users/{user_id}/episodes
func (c *Client) GetUserEpisodes(userID int, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/users/%d/episodes", userID)
	return GetPaginated[models.Episode](c, path, pagination.ToMap())
}

// GetShowEpisodes retrieves all episodes of a show.
// API: GET /v2/shows/{show_id}/episodes
func (c *Client) GetShowEpisodes(showID int, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/shows/%d/episodes", showID)
	return GetPaginated[models.Episode](c, path, pagination.ToMap())
}