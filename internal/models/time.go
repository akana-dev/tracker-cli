package models

import (
	"strings"
	"time"
)

var timeFormats = []string{
	"2006-01-02T15:04:05.999999999",       // с микросекундами, без TZ
	"2006-01-02T15:04:05.999999999Z07:00", // с микросекундами и TZ
	"2006-01-02T15:04:05",                 // без микросекунд, без TZ
	"2006-01-02T15:04:05Z07:00",           // RFC3339
	time.RFC3339,
	time.RFC3339Nano,
}

type FlexibleTime struct {
	time.Time
}

func (t *FlexibleTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		t.Time = time.Time{}
		return nil
	}

	for _, format := range timeFormats {
		if parsed, err := time.Parse(format, s); err == nil {
			t.Time = parsed
			return nil
		}
	}

	for _, format := range timeFormats {
		if parsed, err := time.ParseInLocation(format, s, time.Local); err == nil {
			t.Time = parsed
			return nil
		}
	}

	return nil
}

func (t FlexibleTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format("2006-01-02T15:04:05.999999") + `"`), nil
}

func (t FlexibleTime) IsZero() bool {
	return t.Time.IsZero()
}

type FlexibleTimePtr struct {
	*FlexibleTime
}

func (t *FlexibleTimePtr) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		t.FlexibleTime = nil
		return nil
	}

	ft := &FlexibleTime{}
	for _, format := range timeFormats {
		if parsed, err := time.Parse(format, s); err == nil {
			ft.Time = parsed
			t.FlexibleTime = ft
			return nil
		}
	}

	for _, format := range timeFormats {
		if parsed, err := time.ParseInLocation(format, s, time.Local); err == nil {
			ft.Time = parsed
			t.FlexibleTime = ft
			return nil
		}
	}

	return nil
}
