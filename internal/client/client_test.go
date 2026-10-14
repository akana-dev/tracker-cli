package client

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type testNetError struct {
	err       error
	timeout   bool
	temporary bool
}

func (e *testNetError) Error() string   { return e.err.Error() }
func (e *testNetError) Timeout() bool   { return e.timeout }
func (e *testNetError) Temporary() bool { return e.temporary }

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		want       bool
	}{
		{
			name:       "timeout net.Error",
			err:        &testNetError{err: errors.New("timeout"), timeout: true},
			statusCode: 0,
			want:       true,
		},
		{
			name:       "non-timeout net.Error",
			err:        &testNetError{err: errors.New("some net error"), timeout: false},
			statusCode: 0,
			want:       false,
		},
		{
			name:       "url.Error timeout",
			err:        &url.Error{Op: "Post", URL: "http://example.com", Err: &testNetError{err: errors.New("timeout"), timeout: true}},
			statusCode: 0,
			want:       true,
		},
		{
			name:       "url.Error with net.OpError",
			err:        &url.Error{Op: "Post", URL: "http://example.com", Err: &net.OpError{Op: "dial", Err: errors.New("connection refused")}},
			statusCode: 0,
			want:       true,
		},
		{
			name:       "url.Error non-retryable",
			err:        &url.Error{Op: "Post", URL: "http://example.com", Err: errors.New("some error")},
			statusCode: 0,
			want:       false,
		},
		{
			name:       "HTTP 408",
			err:        nil,
			statusCode: http.StatusRequestTimeout,
			want:       true,
		},
		{
			name:       "HTTP 429",
			err:        nil,
			statusCode: http.StatusTooManyRequests,
			want:       true,
		},
		{
			name:       "HTTP 502",
			err:        nil,
			statusCode: http.StatusBadGateway,
			want:       true,
		},
		{
			name:       "HTTP 503",
			err:        nil,
			statusCode: http.StatusServiceUnavailable,
			want:       true,
		},
		{
			name:       "HTTP 504",
			err:        nil,
			statusCode: http.StatusGatewayTimeout,
			want:       true,
		},
		{
			name:       "HTTP 400",
			err:        nil,
			statusCode: http.StatusBadRequest,
			want:       false,
		},
		{
			name:       "HTTP 401",
			err:        nil,
			statusCode: http.StatusUnauthorized,
			want:       false,
		},
		{
			name:       "HTTP 500",
			err:        nil,
			statusCode: http.StatusInternalServerError,
			want:       false,
		},
		{
			name:       "no error success status",
			err:        nil,
			statusCode: http.StatusOK,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryableError(tt.err, tt.statusCode); got != tt.want {
				t.Errorf("isRetryableError(%v, %d) = %v, want %v", tt.err, tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestParseRetryAfter(t *testing.T) {
	const tolerance = 500 * time.Millisecond

	tests := []struct {
		name       string
		value      string
		expected   time.Duration
		isHTTPDate bool
	}{
		{"empty", "", 0, false},
		{"zero seconds", "0", 0, false},
		{"negative seconds", "-5", 0, false},
		{"positive seconds", "42", 42 * time.Second, false},
		{"invalid string", "abc", 0, false},
		{"future HTTP date", time.Now().Add(2 * time.Hour).UTC().Format(http.TimeFormat), 2 * time.Hour, true},
		{"past HTTP date", time.Now().Add(-2 * time.Hour).UTC().Format(http.TimeFormat), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.expected
			if tt.isHTTPDate {
				parsed, err := http.ParseTime(tt.value)
				if err != nil {
					t.Fatalf("failed to parse HTTP date %q: %v", tt.value, err)
				}
				expected = time.Until(parsed)
				if expected <= 0 {
					expected = 0
				}
			}

			got := parseRetryAfter(tt.value)
			if expected == 0 {
				if got != 0 {
					t.Errorf("parseRetryAfter(%q) = %v, want 0", tt.value, got)
				}
				return
			}
			if got < expected-tolerance || got > expected+tolerance {
				t.Errorf("parseRetryAfter(%q) = %v, want ~%v", tt.value, got, expected)
			}
		})
	}
}

func TestParseRetryAfterInvalid(t *testing.T) {
	cases := []string{
		"not-a-date",
		"123abc",
		"Mon, 99 Jan 2026 00:00:00 GMT",
	}

	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			if got := parseRetryAfter(value); got != 0 {
				t.Errorf("parseRetryAfter(%q) = %v, want 0", value, got)
			}
		})
	}
}

func TestParseRetryAfterSecondsBoundary(t *testing.T) {
	// Ensure strconv parsing returns 0 for non-numeric strings and correct values for numbers.
	if got := parseRetryAfter(strconv.Itoa(0)); got != 0 {
		t.Errorf("parseRetryAfter(\"0\") = %v, want 0", got)
	}
	if got := parseRetryAfter(strconv.Itoa(1)); got != 1*time.Second {
		t.Errorf("parseRetryAfter(\"1\") = %v, want %v", got, 1*time.Second)
	}
}
