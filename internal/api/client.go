package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Default values for the client
const (
	DefaultBaseURL    = "https://api.spreaker.com"
	DefaultAPIVersion = "v2"
	DefaultTimeout    = 30 * time.Second
)


type Client struct {
	BaseURL string
	APIVersion string
	Token string
	HTTPClient *http.Client
	UserAgent string
}

// NewClient creates a new Spreaker API client with the given OAuth token.
// If token is empty, only public (unauthenticated) endpoints will work.
func NewClient(token string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		APIVersion: DefaultAPIVersion,
		Token:      token,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		UserAgent: "spreaker-cli/1.0",
	}
}

// NewClientWithOptions creates a client with custom configuration.
// Useful for testing or pointing to different environments.
func NewClientWithOptions(token, baseURL string, timeout time.Duration) *Client {
	client := NewClient(token)
	if baseURL != "" {
		client.BaseURL = baseURL
	}
	if timeout > 0 {
		client.HTTPClient.Timeout = timeout
	}
	return client
}

// -----------------------------------------------------------------------------
// API Error Handling
// -----------------------------------------------------------------------------

// APIError represents an error response from the Spreaker API.
type APIError struct {
	StatusCode int      // HTTP status code
	Code       int      // Spreaker error code
	Messages   []string // Error messages from the API
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Messages) > 0 {
		return fmt.Sprintf("spreaker API error %d: %s", e.StatusCode, e.Messages[0])
	}
	return fmt.Sprintf("spreaker API error %d", e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// -----------------------------------------------------------------------------
// API Response Wrapper
// -----------------------------------------------------------------------------

// apiResponse wraps the Spreaker API response format.
// All Spreaker responses are wrapped in a "response" object.
type apiResponse struct {
	Response json.RawMessage `json:"response"`
}

// errorResponse represents the error format in Spreaker responses.
type errorResponse struct {
	Error struct {
		Messages []string `json:"messages"`
		Code     int      `json:"code"`
	} `json:"error"`
}

// paginatedResponse represents a paginated list response.
type paginatedResponse struct {
	Items   json.RawMessage `json:"items"`
	NextURL string          `json:"next_url"`
}

// -----------------------------------------------------------------------------
// HTTP Request Methods
// -----------------------------------------------------------------------------

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s/%s%s", c.BaseURL, c.APIVersion, path)
}

// newRequest creates a new HTTP request with common headers set.
func (c *Client) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	// Set authorization header if we have a token
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}

// do executes an HTTP request and handles the response.
// It unmarshals the response into the provided result pointer.
func (c *Client) do(req *http.Request, result interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error responses (4xx, 5xx)
	if resp.StatusCode >= 400 {
		return c.parseErrorResponse(resp.StatusCode, body)
	}

	// If no result is expected, we're done
	if result == nil {
		return nil
	}

	// Parse the response wrapper
	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Unmarshal the inner response into the result
	if err := json.Unmarshal(apiResp.Response, result); err != nil {
		return fmt.Errorf("failed to parse response data: %w", err)
	}

	return nil
}

// parseErrorResponse extracts error information from an API error response.
func (c *Client) parseErrorResponse(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}

	// Try to parse the error response
	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		var errResp errorResponse
		if err := json.Unmarshal(apiResp.Response, &errResp); err == nil {
			apiErr.Code = errResp.Error.Code
			apiErr.Messages = errResp.Error.Messages
		}
	}

	return apiErr
}

// -----------------------------------------------------------------------------
// HTTP Verb Helpers
// -----------------------------------------------------------------------------

func (c *Client) Get(path string, params map[string]string, result interface{}) error {
	// Build URL with query parameters
	urlStr := c.buildURL(path)
	if len(params) > 0 {
		query := url.Values{}
		for k, v := range params {
			query.Set(k, v)
		}
		urlStr = urlStr + "?" + query.Encode()
	}

	req, err := c.newRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

	return c.do(req, result)
}

func (c *Client) Post(path string, body interface{}, result interface{}) error {
	// Serialize body to JSON
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to serialize request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := c.newRequest(http.MethodPost, c.buildURL(path), bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.do(req, result)
}

// This is used for endpoints that accept form fields data (multipart/form-data), like episode uploads.
func (c *Client) PostForm(path string, fields map[string]string, result interface{}) error {
	// Create a buffer to write the multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, c.buildURL(path), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.do(req, result)
}

// PostFormWithFile performs a POST request with form data including a file upload.
// This is used for uploading episode audio files.
func (c *Client) PostFormWithFile(path string, fields map[string]string, fileField, filePath string, result interface{}) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a buffer to write the multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file field
	part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file contents to the form
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file to form: %w", err)
	}

	// Add other form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, c.buildURL(path), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.do(req, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string, result interface{}) error {
	req, err := c.newRequest(http.MethodDelete, c.buildURL(path), nil)
	if err != nil {
		return err
	}

	return c.do(req, result)
}

// Put performs a PUT request (used for some Spreaker endpoints like follow/favorite).
func (c *Client) Put(path string, result interface{}) error {
	req, err := c.newRequest(http.MethodPut, c.buildURL(path), nil)
	if err != nil {
		return err
	}

	return c.do(req, result)
}

// -----------------------------------------------------------------------------
// Pagination Helper
// -----------------------------------------------------------------------------

// PaginatedResult holds a page of results plus the URL for the next page.
type PaginatedResult[T any] struct {
	Items   []T
	NextURL string
	HasMore bool
}

// GetPaginated performs a GET request and parses a paginated response.
// T is the type of items in the list.
func GetPaginated[T any](c *Client, path string, params map[string]string) (*PaginatedResult[T], error) {
	// Build URL with query parameters
	urlStr := c.buildURL(path)
	if len(params) > 0 {
		query := url.Values{}
		for k, v := range params {
			query.Set(k, v)
		}
		urlStr = urlStr + "?" + query.Encode()
	}

	req, err := c.newRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseErrorResponse(resp.StatusCode, body)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var paginated paginatedResponse
	if err := json.Unmarshal(apiResp.Response, &paginated); err != nil {
		return nil, fmt.Errorf("failed to parse paginated response: %w", err)
	}

	var items []T
	if err := json.Unmarshal(paginated.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to parse items: %w", err)
	}

	return &PaginatedResult[T]{
		Items:   items,
		NextURL: paginated.NextURL,
		HasMore: paginated.NextURL != "",
	}, nil
}

// -----------------------------------------------------------------------------
// Convenience: Pagination Parameters
// -----------------------------------------------------------------------------

type PaginationParams struct {
	Limit  int 
	Offset int 
}

func (p PaginationParams) ToMap() map[string]string {
	params := make(map[string]string)
	if p.Limit > 0 {
		params["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Offset > 0 {
		params["offset"] = strconv.Itoa(p.Offset)
	}
	return params
}


func (c *Client) CheckAuth() error {
    if c.token == "" {
        return fmt.Errorf("authentication required: this endpoint requires an OAuth token")
    }
    return nil
}
