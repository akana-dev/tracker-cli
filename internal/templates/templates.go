package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"tracker/internal/config"
	"tracker/internal/models"
	"tracker/internal/service"
)

var templatesDir = filepath.Join(config.ConfigDir, "templates")

var templateNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func templateFilePath(name string) string {
	return filepath.Join(templatesDir, name+".yaml")
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("имя шаблона не может быть пустым")
	}
	if !templateNameRegex.MatchString(name) {
		return fmt.Errorf("имя шаблона может содержать только буквы, цифры, дефисы, подчёркивания и точки")
	}
	return nil
}

func Load(name string) (*models.Template, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	path := templateFilePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("шаблон %q не найден", name)
		}
		return nil, fmt.Errorf("ошибка чтения шаблона: %w", err)
	}

	var tmpl models.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	tmpl.Name = name
	return &tmpl, nil
}

func Save(tmpl *models.Template) error {
	if err := validateName(tmpl.Name); err != nil {
		return err
	}
	if err := service.ValidateTitle(tmpl.Title); err != nil {
		return fmt.Errorf("шаблон невалиден: %w", err)
	}

	if err := os.MkdirAll(templatesDir, 0700); err != nil {
		return fmt.Errorf("не удалось создать директорию шаблонов: %w", err)
	}

	data, err := yaml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("ошибка сериализации YAML: %w", err)
	}

	path := templateFilePath(tmpl.Name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("ошибка записи шаблона: %w", err)
	}

	return nil
}

func Delete(name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	path := templateFilePath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("шаблон %q не найден", name)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("ошибка удаления шаблона: %w", err)
	}

	return nil
}

func List() ([]*models.Template, error) {
	if err := os.MkdirAll(templatesDir, 0700); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории шаблонов: %w", err)
	}

	var templates []*models.Template
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		tmpl, err := Load(name)
		if err != nil {
			continue
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}

func Exists(name string) bool {
	path := templateFilePath(name)
	_, err := os.Stat(path)
	return err == nil
}

func GenerateExample() string {
	return `# Имя шаблона (используется для вызова)
name: example

# Название задачи (обязательно)
title: "Ежедневная планёрка"

# Компания (опционально)
company: COMP1

# Исполнитель (опционально)
assignee: ""

# Статус решения по умолчанию (опционально)
solution: ""

# Комментарий (опционально)
comment: "Регулярная встреча команды"

# Теги (опционально)
tags:
  - meeting
  - daily
`
}

func ParseYAML(data []byte) (*models.Template, error) {
	var tmpl models.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func MarshalYAML(tmpl *models.Template) ([]byte, error) {
	return yaml.Marshal(tmpl)
}
