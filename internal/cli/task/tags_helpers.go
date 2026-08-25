package task

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"tracker/internal/app"
	"tracker/internal/models"
)

var tagCache struct {
	sync.RWMutex
	tags   []models.Tag
	expire time.Time
}

const tagCacheTTL = 5 * time.Minute

func getTagsFromCache() ([]models.Tag, error) {
	tagCache.RLock()
	if time.Now().Before(tagCache.expire) && len(tagCache.tags) > 0 {
		tags := tagCache.tags
		tagCache.RUnlock()
		return tags, nil
	}
	tagCache.RUnlock()

	tags, err := app.GetClient().ListTags("")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список тегов: %w", err)
	}

	tagCache.Lock()
	tagCache.tags = tags
	tagCache.expire = time.Now().Add(tagCacheTTL)
	tagCache.Unlock()

	return tags, nil
}

func resolveTagNamesToIDs(names []string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}

	allTags, err := getTagsFromCache()
	if err != nil {
		return nil, err
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

func ClearTagCache() {
	tagCache.Lock()
	tagCache.tags = nil
	tagCache.expire = time.Time{}
	tagCache.Unlock()
}
