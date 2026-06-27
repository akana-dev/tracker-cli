package service

import (
	"fmt"
	"strings"
	"time"

	"tracker/internal/models"
	"tracker/internal/ui"
)

func FormatEndTime(startTime, endTime time.Time) string {
	startLocal := startTime.Local()
	endLocal := endTime.Local()

	if endLocal.Year() != startLocal.Year() || endLocal.YearDay() != startLocal.YearDay() {
		return endLocal.Format("02.01 15:04")
	}
	return endLocal.Format("15:04")
}

func FormatSessions(t models.Task) string {
	if len(t.Sessions) == 0 {
		return ui.Dim("—")
	}

	var lines []string
	for _, s := range t.Sessions {
		startLocal := s.StartTime.Time.Local()
		startStr := startLocal.Format("02.01 15:04")

		if s.EndTime != nil && !s.EndTime.IsZero() {
			endStr := FormatEndTime(s.StartTime.Time, s.EndTime.Time)
			duration := s.EndTime.Time.UTC().Sub(s.StartTime.Time.UTC()).Hours()

			if duration < 0 {
				lines = append(lines, ui.Error(
					fmt.Sprintf("%s - %s (ошибка времени!)", startStr, endStr)))
			} else {
				lines = append(lines, fmt.Sprintf("%s - %s (%.1fч)",
					startStr, endStr, duration))
			}
		} else {
			if t.IsPaused() {
				pauseDuration := t.PausedAt.Time.UTC().Sub(s.StartTime.Time.UTC()).Hours()
				if pauseDuration < 0 {
					pauseDuration = 0
				}
				lines = append(lines, ui.Paused(
					fmt.Sprintf("%s - на паузе (%.1fч)", startStr, pauseDuration)))
			} else {
				hoursWorking := time.Since(startLocal).Hours()
				if hoursWorking < 0 {
					hoursWorking = 0
				}
				if hoursWorking > ErrorHoursWorked {
					lines = append(lines, ui.Error(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				} else if hoursWorking > WarningHoursWorked {
					lines = append(lines, ui.Warning(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				} else {
					lines = append(lines, ui.InProgress(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

func FormatStatus(t models.Task) string {
	solution := "—"
	if t.Solution != nil {
		solution = *t.Solution
	}

	switch GetTaskStatus(&t) {
	case TaskStatusPaused:
		return ui.Paused(fmt.Sprintf("%s (на паузе)", solution))
	case TaskStatusClosed:
		return ui.Closed(solution)
	default:
		return ui.InProgress(solution)
	}
}

func FormatTaskCell(t models.Task) string {
	result := ui.Bold(t.Title)

	if t.Comment != nil && *t.Comment != "" {
		comment := strings.TrimSpace(*t.Comment)
		if comment != "" {
			result += "\n" + ui.Dim(comment)
		}
	}

	return result
}

func PrintIndented(text, indent string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		wrappedLines := WrapText(line, WrapWidth)
		for _, wl := range wrappedLines {
			fmt.Println(indent + wl)
		}
	}
}

func WrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var result []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result = append(result, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		result = append(result, currentLine)
	}

	return result
}
