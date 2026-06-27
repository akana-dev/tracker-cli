package presets

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"tracker/internal/config"
)

var presetsFile = filepath.Join(config.ConfigDir, "export-presets.yaml")

type ExportPreset struct {
	// Name — имя пресета (идентификатор)
	Name string `yaml:"name" json:"name"`

	// Format — формат экспорта (csv, json, xlsx)
	Format string `yaml:"format" json:"format"`

	// Period — относительный период (опционально)
	// Примеры: "today", "week", "month", "last 7 days"
	Period string `yaml:"period,omitempty" json:"period,omitempty"`

	// DateFrom — дата начала (опционально, переопределяет period)
	DateFrom string `yaml:"date_from,omitempty" json:"date_from,omitempty"`

	// DateTo — дата конца (опционально)
	DateTo string `yaml:"date_to,omitempty" json:"date_to,omitempty"`

	// Timezone — часовой пояс
	Timezone string `yaml:"timezone,omitempty" json:"timezone,omitempty"`

	// Filters
	Company  string `yaml:"company,omitempty" json:"company,omitempty"`
	Solution string `yaml:"solution,omitempty" json:"solution,omitempty"`
	Assignee string `yaml:"assignee,omitempty" json:"assignee,omitempty"`
	Search   string `yaml:"search,omitempty" json:"search,omitempty"`
	Ticket   string `yaml:"ticket,omitempty" json:"ticket,omitempty"`

	// Status filters
	OpenOnly   bool `yaml:"open_only,omitempty" json:"open_only,omitempty"`
	ClosedOnly bool `yaml:"closed_only,omitempty" json:"closed_only,omitempty"`
	PausedOnly bool `yaml:"paused_only,omitempty" json:"paused_only,omitempty"`
	ActiveOnly bool `yaml:"active_only,omitempty" json:"active_only,omitempty"`
	AllUsers   bool `yaml:"all_users,omitempty" json:"all_users,omitempty"`

	// Column selection
	Fields []string `yaml:"fields,omitempty" json:"fields,omitempty"`

	// Output settings
	Filename string `yaml:"filename,omitempty" json:"filename,omitempty"`

	// Description — описание пресета (опционально)
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

func LoadPresets() (map[string]*ExportPreset, error) {
	data, err := os.ReadFile(presetsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*ExportPreset), nil
		}
		return nil, fmt.Errorf("ошибка чтения файла пресетов: %w", err)
	}

	var presets map[string]*ExportPreset
	if err := yaml.Unmarshal(data, &presets); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	if presets == nil {
		presets = make(map[string]*ExportPreset)
	}

	return presets, nil
}

func SavePresets(presets map[string]*ExportPreset) error {
	if err := os.MkdirAll(config.ConfigDir, 0700); err != nil {
		return fmt.Errorf("не удалось создать директорию: %w", err)
	}

	data, err := yaml.Marshal(presets)
	if err != nil {
		return fmt.Errorf("ошибка сериализации YAML: %w", err)
	}

	if err := os.WriteFile(presetsFile, data, 0600); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

func Get(name string) (*ExportPreset, error) {
	presets, err := LoadPresets()
	if err != nil {
		return nil, err
	}

	preset, ok := presets[name]
	if !ok {
		return nil, fmt.Errorf("пресет %q не найден", name)
	}

	return preset, nil
}

func Save(preset *ExportPreset) error {
	presets, err := LoadPresets()
	if err != nil {
		return err
	}

	if preset.Name == "" {
		return fmt.Errorf("имя пресета не может быть пустым")
	}

	presets[preset.Name] = preset
	return SavePresets(presets)
}

func Delete(name string) error {
	presets, err := LoadPresets()
	if err != nil {
		return err
	}

	if _, ok := presets[name]; !ok {
		return fmt.Errorf("пресет %q не найден", name)
	}

	delete(presets, name)
	return SavePresets(presets)
}

func List() ([]*ExportPreset, error) {
	presets, err := LoadPresets()
	if err != nil {
		return nil, err
	}

	result := make([]*ExportPreset, 0, len(presets))
	for _, preset := range presets {
		result = append(result, preset)
	}

	return result, nil
}

func Exists(name string) bool {
	presets, err := LoadPresets()
	if err != nil {
		return false
	}
	_, ok := presets[name]
	return ok
}
