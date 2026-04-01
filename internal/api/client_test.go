package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// APIError
// ---------------------------------------------------------------------------

func TestAPIError_Error(t *testing.T) {
	t.Run("with messages", func(t *testing.T) {
		e := &APIError{StatusCode: 400, Messages: []string{"bad request"}}
		want := "spreaker API error 400: bad request"
		if got := e.Error(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("without messages", func(t *testing.T) {
		e := &APIError{StatusCode: 500}
		want := "spreaker API error 500"
		if got := e.Error(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestAPIError_StatusChecks(t *testing.T) {
	tests := []struct {
		name   string
		code   int
		isNF   bool
		isUA   bool
		isRL   bool
	}{
		{"not found", 404, true, false, false},
		{"unauthorized", 401, false, true, false},
		{"rate limited", 429, false, false, true},
		{"other", 500, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &APIError{StatusCode: tt.code}
			if e.IsNotFound() != tt.isNF {
				t.Errorf("IsNotFound() = %v, want %v", e.IsNotFound(), tt.isNF)
			}
			if e.IsUnauthorized() != tt.isUA {
				t.Errorf("IsUnauthorized() = %v, want %v", e.IsUnauthorized(), tt.isUA)
			}
			if e.IsRateLimited() != tt.isRL {
				t.Errorf("IsRateLimited() = %v, want %v", e.IsRateLimited(), tt.isRL)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Client construction
// ---------------------------------------------------------------------------

func TestNewClient(t *testing.T) {
	c := NewClient("tok123")
	if c.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", c.BaseURL, DefaultBaseURL)
	}
	if c.APIVersion != DefaultAPIVersion {
		t.Errorf("APIVersion = %q, want %q", c.APIVersion, DefaultAPIVersion)
	}
	if c.token != "tok123" {
		t.Errorf("token = %q, want %q", c.token, "tok123")
	}
	if c.HTTPClient.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", c.HTTPClient.Timeout, DefaultTimeout)
	}
	if c.UserAgent == "" {
		t.Error("UserAgent is empty")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	t.Run("overrides baseURL and timeout", func(t *testing.T) {
		c := NewClientWithOptions("tok", "https://custom.api", 5*time.Second)
		if c.BaseURL != "https://custom.api" {
			t.Errorf("BaseURL = %q, want %q", c.BaseURL, "https://custom.api")
		}
		if c.HTTPClient.Timeout != 5*time.Second {
			t.Errorf("Timeout = %v, want %v", c.HTTPClient.Timeout, 5*time.Second)
		}
	})

	t.Run("empty/zero values keep defaults", func(t *testing.T) {
		c := NewClientWithOptions("tok", "", 0)
		if c.BaseURL != DefaultBaseURL {
			t.Errorf("BaseURL = %q, want %q", c.BaseURL, DefaultBaseURL)
		}
		if c.HTTPClient.Timeout != DefaultTimeout {
			t.Errorf("Timeout = %v, want %v", c.HTTPClient.Timeout, DefaultTimeout)
		}
	})
}

// ---------------------------------------------------------------------------
// CheckAuth
// ---------------------------------------------------------------------------

func TestCheckAuth(t *testing.T) {
	t.Run("empty token returns error", func(t *testing.T) {
		c := NewClient("")
		if err := c.CheckAuth(); err == nil {
			t.Fatal("expected error for empty token")
		}
	})

	t.Run("non-empty token returns nil", func(t *testing.T) {
		c := NewClient("tok")
		if err := c.CheckAuth(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// PaginationParams.ToMap
// ---------------------------------------------------------------------------

func TestPaginationParams_ToMap(t *testing.T) {
	t.Run("zero values yield empty map", func(t *testing.T) {
		m := PaginationParams{}.ToMap()
		if len(m) != 0 {
			t.Errorf("expected empty map, got %v", m)
		}
	})

	t.Run("non-zero values", func(t *testing.T) {
		m := PaginationParams{Limit: 10, Offset: 20}.ToMap()
		if m["limit"] != "10" {
			t.Errorf("limit = %q, want %q", m["limit"], "10")
		}
		if m["offset"] != "20" {
			t.Errorf("offset = %q, want %q", m["offset"], "20")
		}
	})
}

// ---------------------------------------------------------------------------
// newRequest — Authorization header
// ---------------------------------------------------------------------------

func TestNewRequest_AuthHeader(t *testing.T) {
	t.Run("sets Authorization when token present", func(t *testing.T) {
		c := NewClient("mytoken")
		req, err := c.newRequest(http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer mytoken" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer mytoken")
		}
	})

	t.Run("omits Authorization when token empty", func(t *testing.T) {
		c := NewClient("")
		req, err := c.newRequest(http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization = %q, want empty", got)
		}
	})
}

// ---------------------------------------------------------------------------
// HTTP helpers with httptest
// ---------------------------------------------------------------------------

// helper: create a test server returning a Spreaker-format JSON response
func spreakerServer(t *testing.T, statusCode int, responsePayload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		body := map[string]interface{}{"response": responsePayload}
		json.NewEncoder(w).Encode(body)
	}))
}

func testClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	c := NewClient("test-token")
	c.BaseURL = srv.URL
	c.APIVersion = "v2"
	return c
}

func TestGet_Success(t *testing.T) {
	srv := spreakerServer(t, 200, map[string]interface{}{
		"user": map[string]interface{}{
			"user_id":  42,
			"fullname": "Test User",
		},
	})
	defer srv.Close()

	c := testClient(t, srv)

	var result struct {
		User struct {
			UserID   int    `json:"user_id"`
			Fullname string `json:"fullname"`
		} `json:"user"`
	}
	if err := c.Get("/users/self", nil, &result); err != nil {
		t.Fatal(err)
	}
	if result.User.UserID != 42 {
		t.Errorf("UserID = %d, want 42", result.User.UserID)
	}
}

func TestGet_WithParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("expected limit=5, got %q", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": map[string]string{},
		})
	}))
	defer srv.Close()

	c := testClient(t, srv)
	var result map[string]string
	c.Get("/test", map[string]string{"limit": "5"}, &result)
}

func TestGet_ErrorResponse(t *testing.T) {
	srv := spreakerServer(t, 404, map[string]interface{}{
		"error": map[string]interface{}{
			"code":     1001,
			"messages": []string{"not found"},
		},
	})
	defer srv.Close()

	c := testClient(t, srv)
	err := c.Get("/missing", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Code != 1001 {
		t.Errorf("Code = %d, want 1001", apiErr.Code)
	}
	if len(apiErr.Messages) == 0 || apiErr.Messages[0] != "not found" {
		t.Errorf("Messages = %v, want [\"not found\"]", apiErr.Messages)
	}
}

func TestPost_JSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "hello" {
			t.Errorf("title = %q, want %q", body["title"], "hello")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"response": map[string]string{}})
	}))
	defer srv.Close()

	c := testClient(t, srv)
	c.Post("/test", map[string]string{"title": "hello"}, nil)
}

func TestPostForm_MultipartFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("failed to parse multipart: %v", err)
		}
		if r.FormValue("title") != "ep1" {
			t.Errorf("title = %q, want %q", r.FormValue("title"), "ep1")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"response": map[string]string{}})
	}))
	defer srv.Close()

	c := testClient(t, srv)
	c.PostForm("/test", map[string]string{"title": "ep1"}, nil)
}

func TestDelete_Method(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"response": nil})
	}))
	defer srv.Close()

	c := testClient(t, srv)
	c.Delete("/test", nil)
}

func TestPut_Method(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"response": nil})
	}))
	defer srv.Close()

	c := testClient(t, srv)
	c.Put("/test", nil)
}

func TestGetPaginated(t *testing.T) {
	srv := spreakerServer(t, 200, map[string]interface{}{
		"items": []map[string]interface{}{
			{"user_id": 1, "fullname": "Alice"},
			{"user_id": 2, "fullname": "Bob"},
		},
		"next_url": "https://api.spreaker.com/v2/next?page=2",
	})
	defer srv.Close()

	c := testClient(t, srv)

	type simpleUser struct {
		UserID   int    `json:"user_id"`
		Fullname string `json:"fullname"`
	}

	result, err := GetPaginated[simpleUser](c, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if !result.HasMore {
		t.Error("HasMore should be true when next_url is present")
	}
	if result.NextURL == "" {
		t.Error("NextURL should not be empty")
	}
}

func TestGetPaginated_NoMore(t *testing.T) {
	srv := spreakerServer(t, 200, map[string]interface{}{
		"items":    []map[string]interface{}{},
		"next_url": "",
	})
	defer srv.Close()

	c := testClient(t, srv)

	type dummy struct{}
	result, err := GetPaginated[dummy](c, "/empty", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasMore {
		t.Error("HasMore should be false when next_url is empty")
	}
}
