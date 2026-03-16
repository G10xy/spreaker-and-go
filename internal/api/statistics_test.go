package api

import (
	"testing"
)

func TestStatisticsParams_ToMap(t *testing.T) {
	t.Run("empty params yield empty map", func(t *testing.T) {
		m := StatisticsParams{}.ToMap()
		if len(m) != 0 {
			t.Errorf("expected empty map, got %v", m)
		}
	})

	t.Run("all fields set", func(t *testing.T) {
		m := StatisticsParams{
			From:      "2024-01-01",
			To:        "2024-01-31",
			Group:     "day",
			Precision: 2,
		}.ToMap()

		if m["from"] != "2024-01-01" {
			t.Errorf("from = %q, want %q", m["from"], "2024-01-01")
		}
		if m["to"] != "2024-01-31" {
			t.Errorf("to = %q, want %q", m["to"], "2024-01-31")
		}
		if m["group"] != "day" {
			t.Errorf("group = %q, want %q", m["group"], "day")
		}
		if m["precision"] != "2" {
			t.Errorf("precision = %q, want %q", m["precision"], "2")
		}
	})

	t.Run("partial fields", func(t *testing.T) {
		m := StatisticsParams{From: "2024-01-01"}.ToMap()
		if len(m) != 1 {
			t.Errorf("expected 1 entry, got %d", len(m))
		}
		if m["from"] != "2024-01-01" {
			t.Errorf("from = %q, want %q", m["from"], "2024-01-01")
		}
		if _, ok := m["to"]; ok {
			t.Error("to should not be present")
		}
	})
}
