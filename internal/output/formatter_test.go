package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

// ---------------------------------------------------------------------------
// truncate
// ---------------------------------------------------------------------------

func TestTruncate(t *testing.T) {
	tests := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{"within limit", "hi", 10, "hi"},
		{"at limit", "hello", 5, "hello"},
		{"over limit", "hello world", 8, "hello..."},
		{"max<=3 edge", "hello", 3, "hel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.s, tt.max)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.s, tt.max, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// formatDuration
// ---------------------------------------------------------------------------

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		ms   int
		want string
	}{
		{"with hours", 3661000, "1:01:01"},
		{"minutes only", 125000, "2:05"},
		{"zero", 0, "0:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.ms)
			if got != tt.want {
				t.Errorf("formatDuration(%d) = %q, want %q", tt.ms, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Formatter creation
// ---------------------------------------------------------------------------

func TestNew_FormatSelection(t *testing.T) {
	tests := []struct {
		input string
		want  Format
	}{
		{"json", FormatJSON},
		{"table", FormatTable},
		{"plain", FormatPlain},
		{"INVALID", FormatTable},
		{"  JSON  ", FormatJSON},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := New(tt.input)
			if f.format != tt.want {
				t.Errorf("New(%q).format = %q, want %q", tt.input, f.format, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Helper to create a formatter writing to a buffer
// ---------------------------------------------------------------------------

func newTestFormatter(format string) (*Formatter, *bytes.Buffer) {
	f := New(format)
	buf := &bytes.Buffer{}
	f.writer = buf
	return f, buf
}

// ---------------------------------------------------------------------------
// PrintUser
// ---------------------------------------------------------------------------

func TestPrintUser_JSON(t *testing.T) {
	f, buf := newTestFormatter("json")
	user := &models.User{UserID: 1, Fullname: "Alice", Username: "alice"}
	f.PrintUser(user)

	var decoded map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v\noutput: %s", err, buf.String())
	}
	if int(decoded["user_id"].(float64)) != 1 {
		t.Errorf("user_id = %v, want 1", decoded["user_id"])
	}
}

func TestPrintUser_Plain(t *testing.T) {
	f, buf := newTestFormatter("plain")
	user := &models.User{UserID: 1, Fullname: "Alice"}
	f.PrintUser(user)

	out := buf.String()
	if !strings.Contains(out, "1") || !strings.Contains(out, "Alice") {
		t.Errorf("plain output missing expected content: %q", out)
	}
}

func TestPrintUser_Table(t *testing.T) {
	f, buf := newTestFormatter("table")
	user := &models.User{UserID: 1, Fullname: "Alice", Username: "alice"}
	f.PrintUser(user)

	out := buf.String()
	if !strings.Contains(out, "ID:") || !strings.Contains(out, "1") {
		t.Errorf("table output missing ID: %q", out)
	}
}

// ---------------------------------------------------------------------------
// PrintShows
// ---------------------------------------------------------------------------

func TestPrintShows_Table(t *testing.T) {
	f, buf := newTestFormatter("table")
	shows := []models.Show{
		{ShowID: 1, Title: "Show1"},
		{ShowID: 2, Title: "Show2"},
	}
	f.PrintShows(shows)

	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "TITLE") {
		t.Error("table output missing header row")
	}
	if !strings.Contains(out, "Show1") || !strings.Contains(out, "Show2") {
		t.Error("table output missing show data")
	}
}

func TestPrintShows_JSON(t *testing.T) {
	f, buf := newTestFormatter("json")
	shows := []models.Show{{ShowID: 1}, {ShowID: 2}}
	f.PrintShows(shows)

	var decoded []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Errorf("got %d items, want 2", len(decoded))
	}
}

// ---------------------------------------------------------------------------
// PrintMessage / PrintSuccess / PrintError
// ---------------------------------------------------------------------------

func TestPrintMessage(t *testing.T) {
	f, buf := newTestFormatter("table")
	f.PrintMessage("hello world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("output = %q", buf.String())
	}
}

func TestPrintSuccess(t *testing.T) {
	f, buf := newTestFormatter("table")
	f.PrintSuccess("done")
	out := buf.String()
	if !strings.HasPrefix(out, "✓") {
		t.Errorf("expected ✓ prefix, got %q", out)
	}
}
