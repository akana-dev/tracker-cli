package export

import (
	"fmt"
	"strings"
	"time"

	"tracker/internal/dates"
)

func ResolveDates(period, dateFrom, dateTo string) (string, string, error) {
	if period != "" {
		from, to, err := dates.ParsePeriod(period)
		if err != nil {
			return "", "", fmt.Errorf("ошибка в периоде %q: %w", period, err)
		}
		return from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339), nil
	}

	var resolvedFrom string
	if dateFrom != "" {
		t, err := dates.ParseRelativeDate(dateFrom)
		if err != nil {
			if parsed, err2 := time.Parse(time.RFC3339, dateFrom); err2 == nil {
				t = parsed
			} else {
				return "", "", fmt.Errorf("ошибка в date-from %q: %w", dateFrom, err)
			}
		}
		resolvedFrom = t.UTC().Format(time.RFC3339)
	}

	var resolvedTo string
	if dateTo != "" {
		t, err := dates.ParseRelativeDate(dateTo)
		if err != nil {
			if parsed, err2 := time.Parse(time.RFC3339, dateTo); err2 == nil {
				t = parsed
			} else {
				return "", "", fmt.Errorf("ошибка в date-to %q: %w", dateTo, err)
			}
		}
		resolvedTo = t.UTC().Format(time.RFC3339)
	}

	return resolvedFrom, resolvedTo, nil
}

func ResolveFields(fieldsStr, fieldsPreset string) []string {
	if fieldsStr != "" {
		var fields []string
		for _, f := range strings.Split(fieldsStr, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				fields = append(fields, f)
			}
		}
		return fields
	}

	switch strings.ToLower(fieldsPreset) {
	case "minimal":
		return []string{"ticket", "title", "total_hours"}
	case "standard":
		return []string{"ticket", "title", "start_time", "end_time", "total_hours", "company_name", "assignee_username", "solution"}
	case "full":
		return nil
	default:
		return nil
	}
}
