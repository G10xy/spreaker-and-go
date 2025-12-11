package models

type Statistics struct {
	Plays int `json:"plays"`

	Downloads int `json:"downloads"`

	Likes int `json:"likes"`

	Messages int `json:"messages"`
}

type StatisticsResponse struct {
	Statistics Statistics `json:"statistics"`
}

type PlayStatistics struct {
	Date string `json:"date"`

	Plays int `json:"plays"`

	Downloads int `json:"downloads"`
}

type GeographicStatistics struct {
	CountryCode string `json:"country_code"`

	CountryName string `json:"country_name"`

	Plays int `json:"plays"`

	Percentage float64 `json:"percentage"`
}

type DeviceStatistics struct {
	Device string `json:"device"`

	Plays int `json:"plays"`

	Percentage float64 `json:"percentage"`
}

type SourceStatistics struct {
	Source string `json:"source"`

	Plays int `json:"plays"`

	Percentage float64 `json:"percentage"`
}
