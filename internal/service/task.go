package service

import (
	"fmt"
	"time"

	"tracker/internal/models"
)

type TaskStatus int

const (
	TaskStatusActive TaskStatus = iota // В работе
	TaskStatusPaused                   // На паузе
	TaskStatusClosed                   // Закрыта
)

func (s TaskStatus) String() string {
	switch s {
	case TaskStatusActive:
		return "active"
	case TaskStatusPaused:
		return "paused"
	case TaskStatusClosed:
		return "closed"
	default:
		return "unknown"
	}
}

func GetTaskStatus(t *models.Task) TaskStatus {
	if t.IsClosed() {
		return TaskStatusClosed
	}
	if t.IsPaused() {
		return TaskStatusPaused
	}
	return TaskStatusActive
}

type TaskStats struct {
	Total      int
	Active     int
	Paused     int
	Closed     int
	TotalHours float64
}

func CalculateTasksStats(tasks []models.Task) TaskStats {
	stats := TaskStats{Total: len(tasks)}
	for _, t := range tasks {
		stats.TotalHours += CalculateTaskHours(t)
		switch GetTaskStatus(&t) {
		case TaskStatusActive:
			stats.Active++
		case TaskStatusPaused:
			stats.Paused++
		case TaskStatusClosed:
			stats.Closed++
		}
	}
	return stats
}

func CalculateTaskHours(t models.Task) float64 {
	totalHours := 0.0
	for _, s := range t.Sessions {
		startUTC := s.StartTime.Time.UTC()

		if s.EndTime != nil && !s.EndTime.IsZero() {
			endUTC := s.EndTime.Time.UTC()
			duration := endUTC.Sub(startUTC).Hours()
			if duration > 0 {
				totalHours += duration
			}
		} else if t.IsPaused() {
			pausedUTC := t.PausedAt.Time.UTC()
			duration := pausedUTC.Sub(startUTC).Hours()
			if duration > 0 {
				totalHours += duration
			}
		} else {
			nowUTC := time.Now().UTC()
			duration := nowUTC.Sub(startUTC).Hours()
			if duration > 0 {
				totalHours += duration
			}
		}
	}
	return totalHours
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	totalMinutes := int(d.Minutes())
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours == 0 {
		return fmt.Sprintf("%d мин", minutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%d ч", hours)
	}
	return fmt.Sprintf("%d ч %d мин", hours, minutes)
}
