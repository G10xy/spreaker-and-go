package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// -----------------------------------------------------------------------------
// Statistics Query Parameters
// -----------------------------------------------------------------------------

type StatisticsParams struct {
	From string
	To string
	Group string
	Precision int
}

func (p StatisticsParams) ToMap() map[string]string {
	params := make(map[string]string)
	if p.From != "" {
		params["from"] = p.From
	}
	if p.To != "" {
		params["to"] = p.To
	}
	if p.Group != "" {
		params["group"] = p.Group
	}
	if p.Precision > 0 {
		params["precision"] = fmt.Sprintf("%d", p.Precision)
	}
	return params
}

// -----------------------------------------------------------------------------
// Overall Statistics
// -----------------------------------------------------------------------------

// GetUserStatistics retrieves a user's all-time overall statistics.
// API: GET /v2/users/{user_id}/statistics
func (c *Client) GetUserStatistics(userID int) (*models.UserOverallStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics", userID)

	var resp models.UserOverallStatisticsResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetMyStatistics is a convenience method to get the authenticated user's overall statistics.
func (c *Client) GetMyStatistics() (*models.UserOverallStatistics, error) {
	me, err := c.GetMe()
	if err != nil {
		return nil, err
	}
	return c.GetUserStatistics(me.UserID)
}

// GetShowStatistics retrieves a show's all-time overall statistics.
// API: GET /v2/shows/{show_id}/statistics
func (c *Client) GetShowStatistics(showID int) (*models.ShowOverallStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics", showID)

	var resp models.ShowOverallStatisticsResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetEpisodeStatistics retrieves an episode's all-time overall statistics.
// API: GET /v2/episodes/{episode_id}/statistics
func (c *Client) GetEpisodeStatistics(episodeID int) (*models.EpisodeOverallStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics", episodeID)

	var resp models.EpisodeOverallStatisticsResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Play Statistics (Time-series)
// -----------------------------------------------------------------------------

// GetUserPlayStatistics retrieves a user's play statistics for a date range.
// API: GET /v2/users/{user_id}/statistics/plays
// Required params: From, To, Group
func (c *Client) GetUserPlayStatistics(userID int, params StatisticsParams) ([]models.PlayStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/plays", userID)

	var resp models.PlayStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetShowPlayStatistics retrieves a show's play statistics for a date range.
// API: GET /v2/shows/{show_id}/statistics/plays
// Required params: From, To, Group
func (c *Client) GetShowPlayStatistics(showID int, params StatisticsParams) ([]models.PlayStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/plays", showID)

	var resp models.PlayStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetEpisodePlayStatistics retrieves an episode's play statistics for a date range.
// API: GET /v2/episodes/{episode_id}/statistics/plays
// Required params: From, To, Group
func (c *Client) GetEpisodePlayStatistics(episodeID int, params StatisticsParams) ([]models.PlayStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics/plays", episodeID)

	var resp models.PlayStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetUserShowsPlayTotals retrieves total play statistics for each show owned by a user.
// API: GET /v2/users/{user_id}/shows/statistics/plays/totals
// Required params: From, To
// Returns a paginated list.
func (c *Client) GetUserShowsPlayTotals(userID int, params StatisticsParams, pagination PaginationParams) (*PaginatedResult[models.ShowPlayTotals], error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/shows/statistics/plays/totals", userID)

	// Merge statistics params with pagination params
	queryParams := params.ToMap()
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.ShowPlayTotals](c, path, queryParams)
}

// GetShowEpisodesPlayTotals retrieves total play statistics for each episode in a show.
// API: GET /v2/shows/{show_id}/episodes/statistics/plays/totals
// Required params: From, To
// Returns a paginated list.
func (c *Client) GetShowEpisodesPlayTotals(showID int, params StatisticsParams, pagination PaginationParams) (*PaginatedResult[models.EpisodePlayTotals], error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/episodes/statistics/plays/totals", showID)

	// Merge statistics params with pagination params
	queryParams := params.ToMap()
	for k, v := range pagination.ToMap() {
		queryParams[k] = v
	}

	return GetPaginated[models.EpisodePlayTotals](c, path, queryParams)
}

// -----------------------------------------------------------------------------
// Likes Statistics (Time-series)
// -----------------------------------------------------------------------------

// GetUserLikesStatistics retrieves a user's likes statistics for a date range.
// API: GET /v2/users/{user_id}/statistics/likes
// Required params: From, To, Group
func (c *Client) GetUserLikesStatistics(userID int, params StatisticsParams) ([]models.LikesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/likes", userID)

	var resp models.LikesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetShowLikesStatistics retrieves a show's likes statistics for a date range.
// API: GET /v2/shows/{show_id}/statistics/likes
// Required params: From, To, Group
func (c *Client) GetShowLikesStatistics(showID int, params StatisticsParams) ([]models.LikesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/likes", showID)

	var resp models.LikesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetEpisodeLikesStatistics retrieves an episode's likes statistics for a date range.
// API: GET /v2/episodes/{episode_id}/statistics/likes
// Required params: From, To, Group
func (c *Client) GetEpisodeLikesStatistics(episodeID int, params StatisticsParams) ([]models.LikesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics/likes", episodeID)

	var resp models.LikesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Followers Statistics (Time-series, User only)
// -----------------------------------------------------------------------------

// GetUserFollowersStatistics retrieves a user's followers statistics for a date range.
// API: GET /v2/users/{user_id}/statistics/followers
// Required params: From, To, Group
func (c *Client) GetUserFollowersStatistics(userID int, params StatisticsParams) ([]models.FollowersStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/followers", userID)

	var resp models.FollowersStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Sources Statistics
// -----------------------------------------------------------------------------

// GetUserSourcesStatistics retrieves a user's play/download sources statistics.
// API: GET /v2/users/{user_id}/statistics/sources
// Required params: From, To, Group
func (c *Client) GetUserSourcesStatistics(userID int, params StatisticsParams) (*models.SourcesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/sources", userID)

	var resp models.SourcesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetShowSourcesStatistics retrieves a show's play/download sources statistics.
// API: GET /v2/shows/{show_id}/statistics/sources
// Required params: From, To, Group
func (c *Client) GetShowSourcesStatistics(showID int, params StatisticsParams) (*models.SourcesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/sources", showID)

	var resp models.SourcesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetEpisodeSourcesStatistics retrieves an episode's play/download sources statistics.
// API: GET /v2/episodes/{episode_id}/statistics/sources
// Required params: From, To, Group
func (c *Client) GetEpisodeSourcesStatistics(episodeID int, params StatisticsParams) (*models.SourcesStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics/sources", episodeID)

	var resp models.SourcesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Devices Statistics
// -----------------------------------------------------------------------------

// GetUserDevicesStatistics retrieves a user's device type statistics.
// API: GET /v2/users/{user_id}/statistics/devices
// Required params: From, To
func (c *Client) GetUserDevicesStatistics(userID int, params StatisticsParams) ([]models.DeviceStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/devices", userID)

	var resp models.DevicesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetShowDevicesStatistics retrieves a show's device type statistics.
// API: GET /v2/shows/{show_id}/statistics/devices
// Required params: From, To
func (c *Client) GetShowDevicesStatistics(showID int, params StatisticsParams) ([]models.DeviceStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/devices", showID)

	var resp models.DevicesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// GetEpisodeDevicesStatistics retrieves an episode's device type statistics.
// API: GET /v2/episodes/{episode_id}/statistics/devices
// Required params: From, To
func (c *Client) GetEpisodeDevicesStatistics(episodeID int, params StatisticsParams) ([]models.DeviceStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics/devices", episodeID)

	var resp models.DevicesStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Operating Systems Statistics
// -----------------------------------------------------------------------------

// GetUserOSStatistics retrieves a user's operating system statistics.
// API: GET /v2/users/{user_id}/statistics/os
// Required params: From, To
func (c *Client) GetUserOSStatistics(userID int, params StatisticsParams) (*models.OSStatisticsBreakdown, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/os", userID)

	var resp models.OSStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetShowOSStatistics retrieves a show's operating system statistics.
// API: GET /v2/shows/{show_id}/statistics/os
// Required params: From, To
// Optional: Precision
func (c *Client) GetShowOSStatistics(showID int, params StatisticsParams) (*models.OSStatisticsBreakdown, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/os", showID)

	var resp models.OSStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetEpisodeOSStatistics retrieves an episode's operating system statistics.
// API: GET /v2/episodes/{episode_id}/statistics/os
// Required params: From, To
// Optional: Precision
func (c *Client) GetEpisodeOSStatistics(episodeID int, params StatisticsParams) (*models.OSStatisticsBreakdown, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/episodes/%d/statistics/os", episodeID)

	var resp models.OSStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Geographic Statistics
// -----------------------------------------------------------------------------

// GetUserGeographicStatistics retrieves a user's geographic statistics.
// API: GET /v2/users/{user_id}/statistics/geographics
// Required params: From, To
func (c *Client) GetUserGeographicStatistics(userID int, params StatisticsParams) (*models.GeographicStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/users/%d/statistics/geographics", userID)

	var resp models.GeographicStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// GetShowGeographicStatistics retrieves a show's geographic statistics.
// API: GET /v2/shows/{show_id}/statistics/geographics
// Required params: From, To
func (c *Client) GetShowGeographicStatistics(showID int, params StatisticsParams) (*models.GeographicStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/geographics", showID)

	var resp models.GeographicStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return &resp.Statistics, nil
}

// -----------------------------------------------------------------------------
// Listeners Statistics (Show only)
// -----------------------------------------------------------------------------

// GetShowListenersStatistics retrieves a show's listeners statistics for a date range.
// API: GET /v2/shows/{show_id}/statistics/listeners
// Required params: From, To, Group
func (c *Client) GetShowListenersStatistics(showID int, params StatisticsParams) ([]models.ListenersStatistics, error) {
	if err := c.CheckAuth(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/shows/%d/statistics/listeners", showID)

	var resp models.ListenersStatisticsResponse
	if err := c.Get(path, params.ToMap(), &resp); err != nil {
		return nil, err
	}

	return resp.Statistics, nil
}