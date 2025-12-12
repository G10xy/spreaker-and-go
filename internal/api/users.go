package api

import (
	"fmt"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)


// GetMe retrieves the authenticated user's profile.
// API: GET /v2/me
func (c *Client) GetMe() (*models.User, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("authentication required: call login first")
	}

	var resp models.UserResponse
	if err := c.Get("/me", nil, &resp); err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// GetUser retrieves a user's public profile by ID.
// API: GET /v2/users/{user_id}
func (c *Client) GetUser(userID int) (*models.User, error) {
	path := fmt.Sprintf("/users/%d", userID)

	var resp models.UserResponse
	if err := c.Get(path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// GetUserShows retrieves all shows belonging to a user.
// API: GET /v2/users/{user_id}/shows
//func (c *Client) GetUserShows(userID int, pagination PaginationParams) (*PaginatedResult[models.Show], error) {
//	path := fmt.Sprintf("/users/%d/shows", userID)
//}

// GetMyShows is a convenience method to get the authenticated user's shows.
// It first retrieves the current user's ID, then fetches their shows.
func (c *Client) GetMyShows(pagination PaginationParams) (*PaginatedResult[models.Show], error) {
	me, err := c.GetMe()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/users/%d/shows", me.UserID)
	return GetPaginated[models.Show](c, path, pagination.ToMap())
}

// GetUserFollowers retrieves a user's followers.
// API: GET /v2/users/{user_id}/followers
func (c *Client) GetUserFollowers(userID int, pagination PaginationParams) (*PaginatedResult[models.User], error) {
	path := fmt.Sprintf("/users/%d/followers", userID)
	return GetPaginated[models.User](c, path, pagination.ToMap())
}

// GetUserFollowings retrieves who a user follows.
// API: GET /v2/users/{user_id}/followings
func (c *Client) GetUserFollowings(userID int, pagination PaginationParams) (*PaginatedResult[models.User], error) {
	path := fmt.Sprintf("/users/%d/followings", userID)
	return GetPaginated[models.User](c, path, pagination.ToMap())
}

// FollowUser follows a user.
// API: PUT /v2/users/{user_id}/followings/{following_id}
// Parameters:
//   - userID: The ID of the authenticated user (the one who wants to follow)
//   - followingID: The ID of the user to follow
func (c *Client) FollowUser(userID, followingID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/followings/%d", userID, followingID)
	return c.Put(path, nil)
}

// UnfollowUser unfollows a user.
// API: DELETE /v2/users/{user_id}/followings/{following_id}
// Parameters:
//   - userID: The ID of the authenticated user (the one who wants to unfollow)
//   - followingID: The ID of the user to unfollow
func (c *Client) UnfollowUser(userID, followingID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/followings/%d", userID, followingID)
	return c.Delete(path, nil)
}

// GetUserEpisodes retrieves all episodes published by a user.
// API: GET /v2/users/{user_id}/episodes
func (c *Client) GetUserEpisodes(userID int, pagination PaginationParams) (*PaginatedResult[models.Episode], error) {
	path := fmt.Sprintf("/users/%d/episodes", userID)
	return GetPaginated[models.Episode](c, path, pagination.ToMap())
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

// UpdateUser updates a user's profile.
// API: POST /v2/users/{user_id}
func (c *Client) UpdateUser(userID int, params UpdateUserParams) (*models.User, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d", userID)

	fields := make(map[string]string)

	if params.Fullname != nil {
		fields["fullname"] = *params.Fullname
	}
	if params.Description != nil {
		fields["description"] = *params.Description
	}
	if params.Gender != nil {
		fields["gender"] = *params.Gender
	}
	if params.Birthday != nil {
		fields["birthday"] = *params.Birthday
	}
	if params.ShowAge != nil {
		if *params.ShowAge {
			fields["show_age"] = "true"
		} else {
			fields["show_age"] = "false"
		}
	}
	if params.Location != nil {
		fields["location"] = *params.Location
	}
	if params.LocationLatitude != nil {
		fields["location_latitude"] = fmt.Sprintf("%f", *params.LocationLatitude)
	}
	if params.LocationLongitude != nil {
		fields["location_longitude"] = fmt.Sprintf("%f", *params.LocationLongitude)
	}
	if params.ContentLanguages != nil {
		fields["content_languages"] = *params.ContentLanguages
	}
	if params.Username != nil {
		fields["username"] = *params.Username
	}
	if params.ContactEmail != nil {
		fields["contact_email"] = *params.ContactEmail
	}

	var resp models.UserResponse
	if err := c.PostForm(path, fields, &resp); err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// GetUserBlocks retrieves a user's blocked users list.
// API: GET /v2/users/{user_id}/blocks
func (c *Client) GetUserBlocks(userID int, pagination PaginationParams) (*PaginatedResult[models.User], error) {
	if c.Token == "" {
		return nil, fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/blocks", userID)
	return GetPaginated[models.User](c, path, pagination.ToMap())
}

// BlockUser blocks a user.
// API: PUT /v2/users/{user_id}/blocks/{blocked_id}
// Parameters:
//   - userID: The ID of the authenticated user (the one who wants to block)
//   - blockedID: The ID of the user to block
func (c *Client) BlockUser(userID, blockedID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/blocks/%d", userID, blockedID)
	return c.Put(path, nil)
}

// UnblockUser unblocks a user.
// API: DELETE /v2/users/{user_id}/blocks/{blocked_id}
// Parameters:
//   - userID: The ID of the authenticated user (the one who wants to unblock)
//   - blockedID: The ID of the user to unblock
func (c *Client) UnblockUser(userID, blockedID int) error {
	if c.Token == "" {
		return fmt.Errorf("authentication required")
	}

	path := fmt.Sprintf("/users/%d/blocks/%d", userID, blockedID)
	return c.Delete(path, nil)
}