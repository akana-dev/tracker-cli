package task

import (
	"reflect"
	"testing"
)

func TestSplitTagNames(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single tag",
			input:    []string{"bug"},
			expected: []string{"bug"},
		},
		{
			name:     "multiple tags in one string",
			input:    []string{"bug, feature, docs"},
			expected: []string{"bug", "feature", "docs"},
		},
		{
			name:     "multiple strings with spaces",
			input:    []string{"  bug  ", "feature,  docs", "ui "},
			expected: []string{"bug", "feature", "docs", "ui"},
		},
		{
			name:     "empty elements discarded",
			input:    []string{"bug,,feature", "", "  ", "docs"},
			expected: []string{"bug", "feature", "docs"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "only empty strings",
			input:    []string{"", "  ", ","},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTagNames(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("splitTagNames(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
