package cmd

import "testing"

func TestExtractTimezone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2026-01-08T11:00:00-05:00", "America/New_York"},
		{"2026-07-08T11:00:00-04:00", "America/New_York"},
		{"2026-01-08T11:00:00-06:00", "America/Chicago"},
		{"2026-01-08T11:00:00-07:00", "America/Denver"},
		{"2026-01-08T11:00:00-08:00", "America/Los_Angeles"},
		{"2026-01-08T16:00:00Z", "UTC"},
		{"2026-01-08T11:00:00+00:00", "UTC"},
		{"invalid", ""},
		{"2026-01-08T11:00:00+05:30", ""}, // India - not mapped
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := extractTimezone(tc.input)
			if got != tc.expected {
				t.Errorf("extractTimezone(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
