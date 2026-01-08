package cmd

import (
	"testing"
	"time"
)

func TestParseTimeExpr_WeekdaySameDay(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 1, 5, 12, 0, 0, 0, loc) // Monday

	got, err := parseTimeExpr("monday", now, loc)
	if err != nil {
		t.Fatalf("parseTimeExpr: %v", err)
	}

	want := time.Date(2026, 1, 5, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestParseTimeExpr_WeekdayNextWeek(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 1, 5, 12, 0, 0, 0, loc) // Monday

	got, err := parseTimeExpr("next monday", now, loc)
	if err != nil {
		t.Fatalf("parseTimeExpr: %v", err)
	}

	want := time.Date(2026, 1, 12, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestStartOfWeek_WeekStartSunday(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 1, 7, 12, 0, 0, 0, loc) // Wednesday

	got := startOfWeek(now, time.Sunday)
	want := time.Date(2026, 1, 4, 0, 0, 0, 0, loc) // Sunday
	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
