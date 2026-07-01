package task

import (
	"fmt"
	"strings"

	"tracker/internal/client"
)

func resolveTagNamesToIDs(names []string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}

	allTags, err := client.ListTags("")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список тегов: %w", err)
	}

	tagMap := make(map[string]int, len(allTags))
	for _, t := range allTags {
		tagMap[strings.ToLower(t.Name)] = t.ID
	}

	ids := make([]int, 0, len(names))
	for _, name := range names {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			continue
		}
		id, ok := tagMap[key]
		if !ok {
			return nil, fmt.Errorf("тег '%s' не найден. Сначала создайте его: tracker tag add %s", name, name)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
