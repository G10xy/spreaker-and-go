package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Episode Messages API
// -----------------------------------------------------------------------------

// GetEpisodeMessages retrieves all messages for an episode.
// API: GET /v2/episodes/{episode_id}/messages
func (c *Client) GetEpisodeMessages(episodeID int, pagination PaginationParams) (*PaginatedResult[models.Message], error) {
	path := fmt.Sprintf("/episodes/%d/messages", episodeID)
	return GetPaginated[models.Message](c, path, pagination.ToMap())
}

// CreateMessage leaves a new message on an episode.
// API: POST /v2/episodes/{episode_id}/messages
// Parameters:
//   - episodeID: The episode ID to leave the message on
//   - text: The message text (max 4000 characters)
func (c *Client) CreateMessage(episodeID int, text string) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	if text == "" {
		return fmt.Errorf("text is required")
	}
	if len(text) > 4000 {
		return fmt.Errorf("text exceeds maximum length of 4000 characters")
	}

	path := fmt.Sprintf("/episodes/%d/messages", episodeID)

	params := map[string]string{
		"text": text,
	}

	return c.Post(path, params, nil)
}

// DeleteMessage deletes a message from an episode.
// The user must be either the message author or the episode owner.
// API: DELETE /v2/episodes/{episode_id}/messages/{message_id}
func (c *Client) DeleteMessage(episodeID, messageID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/messages/%d", episodeID, messageID)
	return c.Delete(path, nil)
}

// ReportMessageAbuse reports a message as spam or as violating Spreaker's terms.
// API: POST /v2/episodes/{episode_id}/messages/{message_id}/report-abuse
func (c *Client) ReportMessageAbuse(episodeID, messageID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/messages/%d/report-abuse", episodeID, messageID)
	return c.Post(path, nil, nil)
}
