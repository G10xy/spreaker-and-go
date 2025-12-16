package models

// -----------------------------------------------------------------------------
// Episode Cuepoint Models
// -----------------------------------------------------------------------------

type Cuepoint struct {
	Timecode int `json:"timecode"`
	AdsMaxCount int `json:"ads_max_count"`
}


type CuepointsResponse struct {
	Cuepoints []Cuepoint `json:"cuepoints"`
}
