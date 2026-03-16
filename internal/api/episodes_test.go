package api

import (
	"testing"
)

// ---------------------------------------------------------------------------
// validateDownloadURL
// ---------------------------------------------------------------------------

func TestValidateDownloadURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid cdn.spreaker.com", "https://cdn.spreaker.com/download/episode/123/file.mp3", false},
		{"valid spreaker.com", "https://spreaker.com/path", false},
		{"valid subdomain spreaker.net", "https://foo.spreaker.net/path", false},
		{"valid spreaker.net bare", "https://spreaker.net/path", false},
		{"reject http scheme", "http://cdn.spreaker.com/path", true},
		{"reject wrong domain", "https://evil.com/path", true},
		{"reject suffix match attack", "https://notspreaker.com/path", true},
		{"reject unparseable URL", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDownloadURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDownloadURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// UploadEpisode parameter validation
// ---------------------------------------------------------------------------

func TestUploadEpisode_Validation(t *testing.T) {
	t.Run("missing auth", func(t *testing.T) {
		c := NewClient("") // no token
		_, err := c.UploadEpisode(1, UploadEpisodeParams{Title: "t", MediaFile: "f.mp3"})
		if err == nil {
			t.Fatal("expected auth error")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		c := NewClient("tok")
		_, err := c.UploadEpisode(1, UploadEpisodeParams{MediaFile: "f.mp3"})
		if err == nil {
			t.Fatal("expected title error")
		}
	})

	t.Run("missing media_file", func(t *testing.T) {
		c := NewClient("tok")
		_, err := c.UploadEpisode(1, UploadEpisodeParams{Title: "t"})
		if err == nil {
			t.Fatal("expected media_file error")
		}
	})
}

// ---------------------------------------------------------------------------
// CreateDraftEpisode parameter validation
// ---------------------------------------------------------------------------

func TestCreateDraftEpisode_Validation(t *testing.T) {
	t.Run("missing auth", func(t *testing.T) {
		c := NewClient("")
		_, err := c.CreateDraftEpisode(CreateDraftEpisodeParams{Title: "t", ShowID: 1})
		if err == nil {
			t.Fatal("expected auth error")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		c := NewClient("tok")
		_, err := c.CreateDraftEpisode(CreateDraftEpisodeParams{ShowID: 1})
		if err == nil {
			t.Fatal("expected title error")
		}
	})

	t.Run("missing show_id", func(t *testing.T) {
		c := NewClient("tok")
		_, err := c.CreateDraftEpisode(CreateDraftEpisodeParams{Title: "t"})
		if err == nil {
			t.Fatal("expected show_id error")
		}
	})
}
