package api

import (
	"encoding/json"
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Episode Cuepoints API
// -----------------------------------------------------------------------------

// GetEpisodeCuepoints retrieves all cuepoints for an episode.
// API: GET /v2/episodes/{episode_id}/cuepoints
func (c *Client) GetEpisodeCuepoints(episodeID int) ([]models.Cuepoint, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/cuepoints", episodeID)

	var resp models.CuepointsResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Cuepoints, nil
}

// UpdateEpisodeCuepoints replaces all episode cuepoints with the provided ones.
// API: POST /v2/episodes/{episode_id}/cuepoints
// Parameters:
//   - episodeID: The episode ID
//   - cuepoints: Array of cuepoints to set (replaces all existing cuepoints)
func (c *Client) UpdateEpisodeCuepoints(episodeID int, cuepoints []models.Cuepoint) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/cuepoints", episodeID)

	cuepointsJSON, err := json.Marshal(cuepoints)
	if err != nil {
		return fmt.Errorf("failed to encode cuepoints: %w", err)
	}

	params := map[string]string{
		"cuepoints": string(cuepointsJSON),
	}

	return c.Post(path, params, nil)
}

// DeleteEpisodeCuepoints deletes all cuepoints for an episode.
// API: DELETE /v2/episodes/{episode_id}/cuepoints
func (c *Client) DeleteEpisodeCuepoints(episodeID int) error {
	if err := c.CheckAuth(); err != nil {
		return err
	}

	path := fmt.Sprintf("/episodes/%d/cuepoints", episodeID)
	return c.Delete(path, nil)
}
