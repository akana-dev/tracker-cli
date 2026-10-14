package dates

import (
	"testing"
	"time"
)

func TestParseRelativeDateAt(t *testing.T) {
	now := time.Date(2026, time.June, 15, 14, 30, 0, 0, time.UTC) // Monday

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"today", "today", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)},
		{"yesterday", "yesterday", time.Date(2026, time.June, 14, 0, 0, 0, 0, time.UTC)},
		{"tomorrow", "tomorrow", time.Date(2026, time.June, 16, 0, 0, 0, 0, time.UTC)},

		{"weekday monday", "monday", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)},
		{"weekday tuesday", "tuesday", time.Date(2026, time.June, 9, 0, 0, 0, 0, time.UTC)},
		{"weekday wednesday", "wednesday", time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC)},
		{"weekday thursday", "thursday", time.Date(2026, time.June, 11, 0, 0, 0, 0, time.UTC)},
		{"weekday friday", "friday", time.Date(2026, time.June, 12, 0, 0, 0, 0, time.UTC)},
		{"weekday saturday", "saturday", time.Date(2026, time.June, 13, 0, 0, 0, 0, time.UTC)},
		{"weekday sunday", "sunday", time.Date(2026, time.June, 14, 0, 0, 0, 0, time.UTC)},

		{"last monday", "last monday", time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)},
		{"last friday", "last friday", time.Date(2026, time.June, 12, 0, 0, 0, 0, time.UTC)},
		{"next friday", "next friday", time.Date(2026, time.June, 19, 0, 0, 0, 0, time.UTC)},
		{"next monday", "next monday", time.Date(2026, time.June, 22, 0, 0, 0, 0, time.UTC)},

		{"this week", "this week", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)},
		{"last week", "last week", time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)},
		{"next week", "next week", time.Date(2026, time.June, 22, 0, 0, 0, 0, time.UTC)},

		{"this month", "this month", time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)},
		{"last month", "last month", time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)},
		{"next month", "next month", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)},

		{"this quarter", "this quarter", time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)},
		{"last quarter", "last quarter", time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},

		{"this year", "this year", time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"last year", "last year", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)},

		{"last 7 days", "last 7 days", time.Date(2026, time.June, 8, 14, 30, 0, 0, time.UTC)},
		{"last 2 weeks", "last 2 weeks", time.Date(2026, time.June, 1, 14, 30, 0, 0, time.UTC)},
		{"past 3 months", "past 3 months", time.Date(2026, time.March, 15, 14, 30, 0, 0, time.UTC)},
		{"last 1 year", "last 1 year", time.Date(2025, time.June, 15, 14, 30, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRelativeDateAt(tt.input, now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tt.expected) {
				t.Errorf("ParseRelativeDateAt(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseRelativeDateAtErrors(t *testing.T) {
	now := time.Date(2026, time.June, 15, 14, 30, 0, 0, time.UTC)

	cases := []string{
		"abc",
		"last abc days",
		"past xyz months",
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseRelativeDateAt(input, now); err == nil {
				t.Errorf("ParseRelativeDateAt(%q) expected error, got nil", input)
			}
		})
	}
}

func TestParsePeriodAt(t *testing.T) {
	now := time.Date(2026, time.June, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name         string
		input        string
		expectedFrom time.Time
		expectedTo   time.Time
	}{
		{
			name:         "today",
			input:        "today",
			expectedFrom: time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "yesterday",
			input:        "yesterday",
			expectedFrom: time.Date(2026, time.June, 14, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2026, time.June, 14, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:         "this week",
			input:        "this week",
			expectedFrom: time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "last week",
			input:        "last week",
			expectedFrom: time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2026, time.June, 14, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:         "this month",
			input:        "this month",
			expectedFrom: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "last month",
			input:        "last month",
			expectedFrom: time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2026, time.May, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			name:         "this quarter",
			input:        "this quarter",
			expectedFrom: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "this year",
			input:        "this year",
			expectedFrom: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "last 7 days",
			input:        "last 7 days",
			expectedFrom: time.Date(2026, time.June, 8, 14, 30, 0, 0, time.UTC),
			expectedTo:   now,
		},
		{
			name:         "past 2 months",
			input:        "past 2 months",
			expectedFrom: time.Date(2026, time.April, 15, 14, 30, 0, 0, time.UTC),
			expectedTo:   now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to, err := ParsePeriodAt(tt.input, now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !from.Equal(tt.expectedFrom) {
				t.Errorf("ParsePeriodAt(%q) from = %v, want %v", tt.input, from, tt.expectedFrom)
			}
			if !to.Equal(tt.expectedTo) {
				t.Errorf("ParsePeriodAt(%q) to = %v, want %v", tt.input, to, tt.expectedTo)
			}
		})
	}
}

func TestParsePeriodAtErrors(t *testing.T) {
	now := time.Date(2026, time.June, 15, 14, 30, 0, 0, time.UTC)

	cases := []struct {
		name  string
		input string
	}{
		{"unknown unit", "last 3 abc"},
		{"invalid number", "last abc days"},
		{"unknown keyword", "unknown"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := ParsePeriodAt(tt.input, now); err == nil {
				t.Errorf("ParsePeriodAt(%q) expected error, got nil", tt.input)
			}
		})
	}
}
