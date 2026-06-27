package dates

import (
	"fmt"
	"strings"
	"time"
)

func ParseRelativeDate(s string) (time.Time, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	now := time.Now()

	switch s {
	case "now":
		return now, nil
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	case "yesterday":
		return time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()), nil
	case "tomorrow":
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()), nil
	}

	weekdays := map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}

	if weekday, ok := weekdays[s]; ok {
		return findWeekday(now, weekday, 0), nil
	}

	for prefix, offset := range map[string]int{"last ": -1, "next ": 1} {
		if strings.HasPrefix(s, prefix) {
			dayName := strings.TrimPrefix(s, prefix)
			if weekday, ok := weekdays[dayName]; ok {
				return findWeekday(now, weekday, offset), nil
			}
		}
	}

	periods := []struct {
		prefix string
		calc   func(time.Time, int) time.Time
	}{
		{"this week", func(t time.Time, _ int) time.Time {
			return findWeekday(t, time.Monday, 0)
		}},
		{"last week", func(t time.Time, _ int) time.Time {
			return findWeekday(t, time.Monday, -1)
		}},
		{"next week", func(t time.Time, _ int) time.Time {
			return findWeekday(t, time.Monday, 1)
		}},
		{"this month", func(t time.Time, _ int) time.Time {
			return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}},
		{"last month", func(t time.Time, _ int) time.Time {
			return time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
		}},
		{"next month", func(t time.Time, _ int) time.Time {
			return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		}},
		{"this quarter", func(t time.Time, _ int) time.Time {
			quarter := (int(t.Month())-1)/3 + 1
			month := time.Month((quarter-1)*3 + 1)
			return time.Date(t.Year(), month, 1, 0, 0, 0, 0, t.Location())
		}},
		{"last quarter", func(t time.Time, _ int) time.Time {
			quarter := (int(t.Month()) - 1) / 3
			if quarter == 0 {
				quarter = 4
			}
			month := time.Month((quarter-1)*3 + 1)
			year := t.Year()
			if month == time.January && int(t.Month()) <= 3 {
				year--
			}
			return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
		}},
		{"this year", func(t time.Time, _ int) time.Time {
			return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
		}},
		{"last year", func(t time.Time, _ int) time.Time {
			return time.Date(t.Year()-1, time.January, 1, 0, 0, 0, 0, t.Location())
		}},
	}

	for _, p := range periods {
		if s == p.prefix {
			return p.calc(now, 0), nil
		}
	}

	for _, pattern := range []string{"last ", "past "} {
		if strings.HasPrefix(s, pattern) {
			rest := strings.TrimPrefix(s, pattern)
			parts := strings.Fields(rest)
			if len(parts) == 2 {
				var n int
				fmt.Sscanf(parts[0], "%d", &n)
				unit := parts[1]

				switch unit {
				case "day", "days":
					return now.AddDate(0, 0, -n), nil
				case "week", "weeks":
					return now.AddDate(0, 0, -n*7), nil
				case "month", "months":
					return now.AddDate(0, -n, 0), nil
				case "year", "years":
					return now.AddDate(-n, 0, 0), nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить относительную дату: %s", s)
}

func findWeekday(t time.Time, weekday time.Weekday, offset int) time.Time {
	currentWeekday := t.Weekday()
	diff := int(weekday - currentWeekday)

	if offset < 0 && diff >= 0 {
		diff -= 7
	} else if offset > 0 && diff <= 0 {
		diff += 7
	} else if offset == 0 && diff < 0 {
	} else if offset == 0 && diff > 0 {
		diff -= 7
	}

	result := t.AddDate(0, 0, diff)
	return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, t.Location())
}

func ParsePeriod(s string) (dateFrom, dateTo time.Time, err error) {
	s = strings.ToLower(strings.TrimSpace(s))
	now := time.Now()

	switch s {
	case "today":
		dateFrom = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dateTo = now
		return dateFrom, dateTo, nil

	case "yesterday":
		yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		return yesterday, yesterday.Add(24*time.Hour - time.Nanosecond), nil

	case "this week":
		dateFrom = findWeekday(now, time.Monday, 0)
		dateTo = now
		return dateFrom, dateTo, nil

	case "last week":
		lastMonday := findWeekday(now, time.Monday, -1)
		lastSunday := lastMonday.Add(7*24*time.Hour - time.Nanosecond)
		return lastMonday, lastSunday, nil

	case "this month":
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		dateTo = now
		return dateFrom, dateTo, nil

	case "last month":
		firstOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -1)
		lastOfLastMonth = time.Date(lastOfLastMonth.Year(), lastOfLastMonth.Month(), lastOfLastMonth.Day(), 23, 59, 59, 0, now.Location())
		return firstOfLastMonth, lastOfLastMonth, nil

	case "this quarter":
		quarter := (int(now.Month())-1)/3 + 1
		firstMonth := time.Month((quarter-1)*3 + 1)
		dateFrom = time.Date(now.Year(), firstMonth, 1, 0, 0, 0, 0, now.Location())
		dateTo = now
		return dateFrom, dateTo, nil

	case "this year":
		dateFrom = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		dateTo = now
		return dateFrom, dateTo, nil
	}

	parts := strings.Fields(s)
	if len(parts) == 3 && (parts[0] == "last" || parts[0] == "past") {
		var n int
		fmt.Sscanf(parts[1], "%d", &n)
		unit := parts[2]

		dateTo = now
		switch unit {
		case "day", "days":
			dateFrom = now.AddDate(0, 0, -n)
		case "week", "weeks":
			dateFrom = now.AddDate(0, 0, -n*7)
		case "month", "months":
			dateFrom = now.AddDate(0, -n, 0)
		case "year", "years":
			dateFrom = now.AddDate(-n, 0, 0)
		default:
			return time.Time{}, time.Time{}, fmt.Errorf("неизвестная единица: %s", unit)
		}
		return dateFrom, dateTo, nil
	}

	dateFrom, err = ParseRelativeDate(s)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	dateTo = now
	return dateFrom, dateTo, nil
}
