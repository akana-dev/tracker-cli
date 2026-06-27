package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"tracker/internal/config"
)

var tagsFile = filepath.Join(config.ConfigDir, "tags.json")

type TagsStore map[string]map[string][]string

var (
	tagsCache TagsStore
	tagsMu    sync.RWMutex
)

func init() {
	tagsCache = nil
}

func loadTags() (TagsStore, error) {
	tagsMu.RLock()
	if tagsCache != nil {
		cached := tagsCache
		tagsMu.RUnlock()
		return cached, nil
	}
	tagsMu.RUnlock()

	data, err := os.ReadFile(tagsFile)
	if err != nil {
		if os.IsNotExist(err) {
			store := make(TagsStore)
			setTagsCache(store)
			return store, nil
		}
		return nil, err
	}

	var store TagsStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("ошибка парсинга тегов: %w", err)
	}

	if store == nil {
		store = make(TagsStore)
	}

	setTagsCache(store)
	return store, nil
}

func setTagsCache(store TagsStore) {
	tagsMu.Lock()
	tagsCache = store
	tagsMu.Unlock()
}

func saveTags(store TagsStore) error {
	if err := os.MkdirAll(config.ConfigDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tagsFile, data, 0600); err != nil {
		return err
	}

	setTagsCache(store)
	return nil
}

func getCurrentServerName() (string, error) {
	server, err := config.GetCurrentServer()
	if err != nil {
		return "", err
	}
	return server.Name, nil
}

func normalizeTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(tag))
}

func Get(ticket string) ([]string, error) {
	serverName, err := getCurrentServerName()
	if err != nil {
		return nil, err
	}

	store, err := loadTags()
	if err != nil {
		return nil, err
	}

	serverTags, ok := store[serverName]
	if !ok {
		return []string{}, nil
	}

	tags, ok := serverTags[ticket]
	if !ok {
		return []string{}, nil
	}

	result := make([]string, len(tags))
	copy(result, tags)
	sort.Strings(result)
	return result, nil
}

func Set(ticket string, tags []string) error {
	serverName, err := getCurrentServerName()
	if err != nil {
		return err
	}

	store, err := loadTags()
	if err != nil {
		return err
	}

	if _, ok := store[serverName]; !ok {
		store[serverName] = make(map[string][]string)
	}

	normalized := normalizeTags(tags)

	if len(normalized) == 0 {
		delete(store[serverName], ticket)
	} else {
		store[serverName][ticket] = normalized
	}

	return saveTags(store)
}

func Add(ticket string, tags []string) error {
	existing, err := Get(ticket)
	if err != nil {
		return err
	}

	combined := append(existing, tags...)
	return Set(ticket, combined)
}

func Remove(ticket string, tags []string) error {
	existing, err := Get(ticket)
	if err != nil {
		return err
	}

	toRemove := make(map[string]bool)
	for _, t := range tags {
		toRemove[normalizeTag(t)] = true
	}

	var result []string
	for _, t := range existing {
		if !toRemove[t] {
			result = append(result, t)
		}
	}

	return Set(ticket, result)
}

func ListAll() ([]string, error) {
	serverName, err := getCurrentServerName()
	if err != nil {
		return nil, err
	}

	store, err := loadTags()
	if err != nil {
		return nil, err
	}

	serverTags, ok := store[serverName]
	if !ok {
		return []string{}, nil
	}

	uniqueTags := make(map[string]bool)
	for _, tags := range serverTags {
		for _, tag := range tags {
			uniqueTags[tag] = true
		}
	}

	result := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		result = append(result, tag)
	}
	sort.Strings(result)
	return result, nil
}

func ListByTicket() (map[string][]string, error) {
	serverName, err := getCurrentServerName()
	if err != nil {
		return nil, err
	}

	store, err := loadTags()
	if err != nil {
		return nil, err
	}

	serverTags, ok := store[serverName]
	if !ok {
		return make(map[string][]string), nil
	}

	result := make(map[string][]string, len(serverTags))
	for ticket, tags := range serverTags {
		tagsCopy := make([]string, len(tags))
		copy(tagsCopy, tags)
		sort.Strings(tagsCopy)
		result[ticket] = tagsCopy
	}

	return result, nil
}

func HasTag(ticket, tag string) (bool, error) {
	tags, err := Get(ticket)
	if err != nil {
		return false, err
	}
	normalized := normalizeTag(tag)
	for _, t := range tags {
		if t == normalized {
			return true, nil
		}
	}
	return false, nil
}

func FilterTicketsByTag(requiredTags []string) ([]string, error) {
	allTags, err := ListByTicket()
	if err != nil {
		return nil, err
	}

	required := make(map[string]bool)
	for _, t := range requiredTags {
		required[normalizeTag(t)] = true
	}

	var result []string
	for ticket, tags := range allTags {
		for _, tag := range tags {
			if required[tag] {
				result = append(result, ticket)
				break
			}
		}
	}

	sort.Strings(result)
	return result, nil
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, t := range tags {
		normalized := normalizeTag(t)
		if normalized == "" {
			continue
		}
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	sort.Strings(result)
	return result
}
