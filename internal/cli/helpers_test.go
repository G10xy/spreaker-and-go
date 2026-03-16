package cli

import (
	"testing"
)

func TestParseIntArg(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    int
		wantErr bool
	}{
		{"valid int", "42", 42, false},
		{"invalid string", "abc", 0, true},
		{"whitespace trimming", "  99  ", 99, false},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIntArg(tt.arg, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIntArg(%q) error = %v, wantErr %v", tt.arg, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseIntArg(%q) = %d, want %d", tt.arg, got, tt.want)
			}
		})
	}
}

func TestParseShowID(t *testing.T) {
	id, err := parseShowID("123")
	if err != nil {
		t.Fatal(err)
	}
	if id != 123 {
		t.Errorf("got %d, want 123", id)
	}

	_, err = parseShowID("abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEpisodeID(t *testing.T) {
	id, err := parseEpisodeID("456")
	if err != nil {
		t.Fatal(err)
	}
	if id != 456 {
		t.Errorf("got %d, want 456", id)
	}
}

func TestParseUserID(t *testing.T) {
	id, err := parseUserID("789")
	if err != nil {
		t.Fatal(err)
	}
	if id != 789 {
		t.Errorf("got %d, want 789", id)
	}
}
