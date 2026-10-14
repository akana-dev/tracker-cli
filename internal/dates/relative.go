package dates

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var weekdays = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}

func ParseRelativeDate(s string) (time.Time, error) {
	return ParseRelativeDateAt(s, time.Now())
}

func ParseRelativeDateAt(s string, now time.Time) (time.Time, error) {
	return parseRelativeDate(s, now)
}

func parseRelativeDate(s string, now time.Time) (time.Time, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	case "now":
		return now, nil
	case "today":
		return startOfDay(now), nil
	case "yesterday":
		return startOfDay(now).AddDate(0, 0, -1), nil
	case "tomorrow":
		return startOfDay(now).AddDate(0, 0, 1), nil
	}

	if weekday, ok := weekdays[s]; ok {
		return findWeekday(now, weekday, 0), nil
	}

	weekdayPrefixes := []struct {
		prefix string
		offset int
	}{
		{prefix: "last ", offset: -1},
		{prefix: "next ", offset: 1},
	}
	for _, item := range weekdayPrefixes {
		if strings.HasPrefix(s, item.prefix) {
			dayName := strings.TrimPrefix(s, item.prefix)
			if weekday, ok := weekdays[dayName]; ok {
				return findWeekday(now, weekday, item.offset), nil
			}
		}
	}

	periods := map[string]func(time.Time) time.Time{
		"this week":    func(t time.Time) time.Time { return startOfWeek(t) },
		"last week":    func(t time.Time) time.Time { return startOfWeek(t).AddDate(0, 0, -7) },
		"next week":    func(t time.Time) time.Time { return startOfWeek(t).AddDate(0, 0, 7) },
		"this month":   func(t time.Time) time.Time { return startOfMonth(t) },
		"last month":   func(t time.Time) time.Time { return startOfMonth(t).AddDate(0, -1, 0) },
		"next month":   func(t time.Time) time.Time { return startOfMonth(t).AddDate(0, 1, 0) },
		"this quarter": func(t time.Time) time.Time { return startOfQuarter(t) },
		"last quarter": func(t time.Time) time.Time { return startOfQuarter(t).AddDate(0, -3, 0) },
		"this year":    func(t time.Time) time.Time { return startOfYear(t) },
		"last year":    func(t time.Time) time.Time { return startOfYear(t).AddDate(-1, 0, 0) },
	}
	if calc, ok := periods[s]; ok {
		return calc(now), nil
	}

	for _, prefix := range []string{"last ", "past "} {
		if strings.HasPrefix(s, prefix) {
			parts := strings.Fields(strings.TrimPrefix(s, prefix))
			if len(parts) == 2 {
				n, err := strconv.Atoi(parts[0])
				if err != nil {
					return time.Time{}, fmt.Errorf("некорректное число периодов %q: %w", parts[0], err)
				}
				if result, ok := addUnits(now, -n, parts[1]); ok {
					return result, nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить относительную дату: %s", s)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func startOfWeek(t time.Time) time.Time {
	return findWeekday(t, time.Monday, 0)
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func startOfQuarter(t time.Time) time.Time {
	month := time.Month(((int(t.Month()) - 1) / 3 * 3) + 1)
	return time.Date(t.Year(), month, 1, 0, 0, 0, 0, t.Location())
}

func startOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

func addUnits(t time.Time, n int, unit string) (time.Time, bool) {
	switch unit {
	case "day", "days":
		return t.AddDate(0, 0, n), true
	case "week", "weeks":
		return t.AddDate(0, 0, n*7), true
	case "month", "months":
		return t.AddDate(0, n, 0), true
	case "year", "years":
		return t.AddDate(n, 0, 0), true
	default:
		return time.Time{}, false
	}
}

func findWeekday(t time.Time, weekday time.Weekday, offset int) time.Time {
	currentWeekday := t.Weekday()
	diff := int(weekday - currentWeekday)

	if offset < 0 && diff >= 0 {
		diff -= 7
	} else if offset > 0 && diff <= 0 {
		diff += 7
	} else if offset == 0 && diff < 0 {
		// targetDay is already in the current week, so no adjustment is needed.
	} else if offset == 0 && diff > 0 {
		diff -= 7
	}

	result := t.AddDate(0, 0, diff)
	return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, t.Location())
}

func ParsePeriod(s string) (dateFrom, dateTo time.Time, err error) {
	return ParsePeriodAt(s, time.Now())
}

func ParsePeriodAt(s string, now time.Time) (dateFrom, dateTo time.Time, err error) {
	return parsePeriod(s, now)
}

func parsePeriod(s string, now time.Time) (dateFrom, dateTo time.Time, err error) {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	case "today":
		return startOfDay(now), now, nil
	case "yesterday":
		yesterday := startOfDay(now).AddDate(0, 0, -1)
		return yesterday, yesterday.Add(24*time.Hour - time.Nanosecond), nil
	case "this week":
		return startOfWeek(now), now, nil
	case "last week":
		lastMonday := startOfWeek(now).AddDate(0, 0, -7)
		return lastMonday, lastMonday.Add(7*time.Hour*24 - time.Nanosecond), nil
	case "this month":
		return startOfMonth(now), now, nil
	case "last month":
		firstOfLastMonth := startOfMonth(now).AddDate(0, -1, 0)
		lastOfLastMonth := firstOfLastMonth.AddDate(0, 1, -1)
		lastOfLastMonth = time.Date(lastOfLastMonth.Year(), lastOfLastMonth.Month(), lastOfLastMonth.Day(), 23, 59, 59, 0, now.Location())
		return firstOfLastMonth, lastOfLastMonth, nil
	case "this quarter":
		return startOfQuarter(now), now, nil
	case "this year":
		return startOfYear(now), now, nil
	}

	parts := strings.Fields(s)
	if len(parts) == 3 && (parts[0] == "last" || parts[0] == "past") {
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("некорректное число периодов %q: %w", parts[1], err)
		}
		dateFrom, ok := addUnits(now, -n, parts[2])
		if !ok {
			return time.Time{}, time.Time{}, fmt.Errorf("неизвестная единица: %s", parts[2])
		}
		return dateFrom, now, nil
	}

	dateFrom, err = parseRelativeDate(s, now)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return dateFrom, now, nil
}
