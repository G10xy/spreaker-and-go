package models

// -----------------------------------------------------------------------------
// Overall Statistics Models
// -----------------------------------------------------------------------------

type UserOverallStatistics struct {
	PlaysCount         int   `json:"plays_count"`
	PlaysOndemandCount int   `json:"plays_ondemand_count"`
	PlaysLiveCount     int   `json:"plays_live_count"`
	ShowsCount         int   `json:"shows_count"`
	EpisodesCount      int   `json:"episodes_count"`
	LikesCount         int   `json:"likes_count"`
	DownloadsCount     int   `json:"downloads_count"`
	FollowersCount     int   `json:"followers_count"`
	User               *User `json:"user,omitempty"`
}

type UserOverallStatisticsResponse struct {
	Statistics UserOverallStatistics `json:"statistics"`
}

type ShowOverallStatistics struct {
	Title              string `json:"title,omitempty"`
	PlaysCount         int    `json:"plays_count"`
	PlaysOndemandCount int    `json:"plays_ondemand_count"`
	PlaysLiveCount     int    `json:"plays_live_count"`
	EpisodesCount      int    `json:"episodes_count"`
	DownloadsCount     int    `json:"downloads_count"`
	LikesCount         int    `json:"likes_count"`
	Show               *Show  `json:"show,omitempty"`
}

type ShowOverallStatisticsResponse struct {
	Statistics ShowOverallStatistics `json:"statistics"`
}

type EpisodeOverallStatistics struct {
	PlaysCount         int      `json:"plays_count"`
	PlaysOndemandCount int      `json:"plays_ondemand_count"`
	PlaysLiveCount     int      `json:"plays_live_count"`
	ChaptersCount      int      `json:"chapters_count"`
	MessagesCount      int      `json:"messages_count"`
	LikesCount         int      `json:"likes_count"`
	DownloadsCount     int      `json:"downloads_count"`
	Episode            *Episode `json:"episode,omitempty"`
}

type EpisodeOverallStatisticsResponse struct {
	Statistics EpisodeOverallStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Play Statistics Models (Time-series)
// -----------------------------------------------------------------------------


type PlayStatistics struct {
	Date               string `json:"date"` // Format: YYYY-MM-DD
	PlaysCount         int    `json:"plays_count"`
	PlaysLiveCount     int    `json:"plays_live_count"`
	PlaysOndemandCount int    `json:"plays_ondemand_count"`
	DownloadsCount     int    `json:"downloads_count"`
}

type PlayStatisticsResponse struct {
	Statistics []PlayStatistics `json:"statistics"`
}

type ShowPlayTotals struct {
	ShowID             int    `json:"show_id"`
	Title              string `json:"title"`
	IsDeleted          bool   `json:"is_deleted"`
	IsTransferred      bool   `json:"is_transferred"`
	PlaysCount         int    `json:"plays_count"`
	PlaysLiveCount     int    `json:"plays_live_count"`
	PlaysOndemandCount int    `json:"plays_ondemand_count"`
	DownloadsCount     int    `json:"downloads_count"`
}

type EpisodePlayTotals struct {
	EpisodeID          int    `json:"episode_id"`
	Title              string `json:"title"`
	IsDeleted          bool   `json:"is_deleted"`
	IsTransferred      bool   `json:"is_transferred"`
	PlaysCount         int    `json:"plays_count"`
	PlaysLiveCount     int    `json:"plays_live_count"`
	PlaysOndemandCount int    `json:"plays_ondemand_count"`
	DownloadsCount     int    `json:"downloads_count"`
}

// -----------------------------------------------------------------------------
// Likes Statistics Models (Time-series)
// -----------------------------------------------------------------------------

type LikesStatistics struct {
	Date       string `json:"date"` // Format: YYYY-MM-DD
	LikesCount int    `json:"likes_count"`
}

type LikesStatisticsResponse struct {
	Statistics []LikesStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Followers Statistics Models (Time-series, User only)
// -----------------------------------------------------------------------------

type FollowersStatistics struct {
	Date           string `json:"date"` // Format: YYYY-MM-DD
	FollowersCount int    `json:"followers_count"`
}

type FollowersStatisticsResponse struct {
	Statistics []FollowersStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Sources Statistics Models
// -----------------------------------------------------------------------------

type SourceOverall struct {
	Name       string `json:"name"`
	PlaysCount int    `json:"plays_count"`
	Percentage int    `json:"percentage"`
}

type SourceDetail map[string]interface{}

type SourcesStatistics struct {
	Overall []SourceOverall `json:"overall"`
	Details []SourceDetail  `json:"details"`
}

type SourcesStatisticsResponse struct {
	Statistics SourcesStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Devices Statistics Models
// -----------------------------------------------------------------------------

type DeviceStatistics struct {
	Name       string  `json:"name"` // Desktop, Mobile, Tablet, Others
	Percentage float64 `json:"percentage"`
}

type DevicesStatisticsResponse struct {
	Statistics []DeviceStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Operating Systems Statistics Models
// -----------------------------------------------------------------------------

type OSStatistics struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}

type OSStatisticsBreakdown struct {
	Desktop []OSStatistics `json:"desktop"`
	Mobile  []OSStatistics `json:"mobile"`
}

type OSStatisticsResponse struct {
	Statistics OSStatisticsBreakdown `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Geographic Statistics Models
// -----------------------------------------------------------------------------

type GeoStatistics struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}

type GeographicStatistics struct {
	Country []GeoStatistics `json:"country"`
	City    []GeoStatistics `json:"city"`
}

type GeographicStatisticsResponse struct {
	Statistics GeographicStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Listeners Statistics Models (Show only)
// -----------------------------------------------------------------------------

type ListenersStatistics struct {
	Date           string `json:"date"` // Format: YYYY-MM-DD
	ListenersCount int    `json:"listeners_count"`
}

type ListenersStatisticsResponse struct {
	Statistics []ListenersStatistics `json:"statistics"`
}

// -----------------------------------------------------------------------------
// Legacy/Simplified Models (for backwards compatibility)
// ----------------------------------------------------------------------------

type Statistics struct {
	Plays     int `json:"plays"`
	Downloads int `json:"downloads"`
	Likes     int `json:"likes"`
	Messages  int `json:"messages"`
}

type StatisticsResponse struct {
	Statistics Statistics `json:"statistics"`
}