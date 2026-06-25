package timeparse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse парсит гибкое время из строки
// Поддерживаемые форматы:
//   - "now" или "" - текущее время
//   - "+30m", "-1h", "+2d" - относительное время
//   - "15:00", "09:30" - только время (сегодня)
//   - "23.06 15:00" - дата и время (текущий год)
//   - "23.06.2026" - только дата
//   - ISO формат (RFC3339)
func Parse(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	now := time.Now()

	if input == "" || strings.ToLower(input) == "now" {
		return now, nil
	}

	if strings.HasPrefix(input, "+") || strings.HasPrefix(input, "-") {
		return parseRelative(input, now)
	}

	if len(input) == 5 && strings.Count(input, ":") == 1 {
		return parseTimeOnly(input, now)
	}

	if len(input) == 11 && strings.Count(input, ".") == 1 && strings.Count(input, ":") == 1 {
		return parseDateTime(input, now)
	}

	if len(input) == 10 && strings.Count(input, ".") == 2 {
		return parseDateOnly(input, now)
	}

	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	if t, err := time.Parse("2006-01-02 15:04:05", input); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить время: %s", input)
}

func parseRelative(input string, now time.Time) (time.Time, error) {
	sign := 1
	if strings.HasPrefix(input, "-") {
		sign = -1
		input = input[1:]
	} else {
		input = input[1:]
	}

	if len(input) < 2 {
		return time.Time{}, fmt.Errorf("неверный формат: %s", input)
	}

	value, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		return time.Time{}, fmt.Errorf("неверное число: %s", input)
	}

	unit := strings.ToLower(input[len(input)-1:])
	switch unit {
	case "m":
		return now.Add(time.Duration(sign*value) * time.Minute), nil
	case "h":
		return now.Add(time.Duration(sign*value) * time.Hour), nil
	case "d":
		return now.AddDate(0, 0, sign*value), nil
	case "w":
		return now.AddDate(0, 0, sign*value*7), nil
	default:
		return time.Time{}, fmt.Errorf("неверная единица: %s (допустимы: m, h, d, w)", unit)
	}
}

func parseTimeOnly(input string, now time.Time) (time.Time, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("неверный формат времени: %s", input)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("неверный час: %s (должен быть 0-23)", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("неверная минута: %s (должна быть 0-59)", parts[1])
	}

	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
}

func parseDateTime(input string, now time.Time) (time.Time, error) {
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("неверный формат: %s", input)
	}

	datePart := parts[0]
	timePart := parts[1]

	dateParts := strings.Split(datePart, ".")
	if len(dateParts) != 2 {
		return time.Time{}, fmt.Errorf("неверная дата: %s", datePart)
	}

	day, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный день: %s", dateParts[0])
	}

	month, err := strconv.Atoi(dateParts[1])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("неверный месяц: %s (должен быть 1-12)", dateParts[1])
	}

	timeParts := strings.Split(timePart, ":")
	if len(timeParts) != 2 {
		return time.Time{}, fmt.Errorf("неверное время: %s", timePart)
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil || hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("неверный час: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil || minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("неверная минута: %s", timeParts[1])
	}

	return time.Date(now.Year(), time.Month(month), day, hour, minute, 0, 0, now.Location()), nil
}

func parseDateOnly(input string, now time.Time) (time.Time, error) {
	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("неверная дата: %s", input)
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("неверный день: %s (должен быть 1-31)", parts[0])
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("неверный месяц: %s (должен быть 1-12)", parts[1])
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный год: %s", parts[2])
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location()), nil
}

func TodayStart() string {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
}

func WeekAgo() string {
	return time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
}

func MonthAgo() string {
	return time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
}

func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

func FormatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

func FormatTime(t time.Time) string {
	return t.Format("15:04")
}
