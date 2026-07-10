package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	CurrentTicket string       `json:"current_ticket"`
	StartedAt     time.Time    `json:"started_at"`
	PausedTasks   []PausedTask `json:"paused_tasks,omitempty"`
}

type PausedTask struct {
	Ticket    string    `json:"ticket"`
	TaskID    int       `json:"task_id"`
	PausedAt  time.Time `json:"paused_at"`
	StartedAt time.Time `json:"started_at"`
	Duration  int64     `json:"duration_seconds"`
}

func GetStatePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
	}

	stateDir := filepath.Join(homeDir, ".tracker")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию %s: %w", stateDir, err)
	}

	return filepath.Join(stateDir, "state.json"), nil
}

func Load() (*State, error) {
	statePath, err := GetStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка чтения state.json: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("ошибка парсинга state.json: %w", err)
	}

	if state.PausedTasks == nil {
		state.PausedTasks = []PausedTask{}
	}

	return &state, nil
}

func Save(state *State) error {
	statePath, err := GetStatePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи state.json: %w", err)
	}

	return nil
}

func Clear() error {
	statePath, err := GetStatePath()
	if err != nil {
		return err
	}

	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ошибка удаления state.json: %w", err)
	}

	return nil
}

func CleanupOldPausedTasks() error {
	s, err := Load()
	if err != nil || s == nil {
		return nil
	}

	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	var filtered []PausedTask
	removedCount := 0
	for _, task := range s.PausedTasks {
		if task.PausedAt.After(cutoff) {
			filtered = append(filtered, task)
		} else {
			removedCount++
		}
	}

	if removedCount > 0 {
		s.PausedTasks = filtered
		if err := Save(s); err != nil {
			return err
		}
		fmt.Printf("Удалено %d задач старше 7 дней из состояния\n", removedCount)
	}
	return nil
}

func (s *State) FindPausedTask(ticket string) *PausedTask {
	for i := range s.PausedTasks {
		if s.PausedTasks[i].Ticket == ticket {
			return &s.PausedTasks[i]
		}
	}
	return nil
}

func (s *State) RemovePausedTask(ticket string) {
	filtered := make([]PausedTask, 0, len(s.PausedTasks))
	for _, task := range s.PausedTasks {
		if task.Ticket != ticket {
			filtered = append(filtered, task)
		}
	}
	s.PausedTasks = filtered
}

func (s *State) AddPausedTask(task PausedTask) {
	for i := range s.PausedTasks {
		if s.PausedTasks[i].Ticket == task.Ticket {
			s.PausedTasks[i] = task
			return
		}
	}
	s.PausedTasks = append(s.PausedTasks, task)
}
