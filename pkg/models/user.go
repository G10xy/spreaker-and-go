package models

type User struct {
	UserID int `json:"user_id"`

	Fullname string `json:"fullname"`

	Username string `json:"username"`

	Description string `json:"description"`

	SiteURL string `json:"site_url"`

	ImageURL string `json:"image_url"`

	ImageOriginalURL string `json:"image_original_url"`

	Kind string `json:"kind"`

	Plan string `json:"plan"`

	FollowersCount int `json:"followers_count"`

	FollowingsCount int `json:"followings_count"`

	ContactEmail string `json:"contact_email,omitempty"`

	Gender string `json:"gender,omitempty"`

	Birthday string `json:"birthday,omitempty"`

	Location string `json:"location,omitempty"`

	LocationLatitude float64 `json:"location_latitude,omitempty"`
	LocationLongitude float64 `json:"location_longitude,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}

// UpdateUserParams contains parameters for updating a user profile.
type UpdateUserParams struct {
	Fullname         *string  
	Description      *string  
	Gender           *string  
	Birthday         *string 
	ShowAge          *bool    
	Location         *string 
	LocationLatitude *float64
	LocationLongitude *float64 
	ContentLanguages *string  
	Username         *string  
	ContactEmail     *string  
}