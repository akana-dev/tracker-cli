package client

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"tracker/internal/models"
)

var (
	templateCache     []models.Template
	templateCacheTime time.Time
	templateCacheMu   sync.RWMutex
	templateCacheTTL  = 5 * time.Minute
)

func listTemplatesCached(includeAll bool) ([]models.Template, error) {
	templateCacheMu.RLock()
	if time.Since(templateCacheTime) < templateCacheTTL && templateCache != nil {
		defer templateCacheMu.RUnlock()
		return templateCache, nil
	}
	templateCacheMu.RUnlock()

	templates, err := ListTemplates(includeAll)
	if err != nil {
		return nil, err
	}

	templateCacheMu.Lock()
	templateCache = templates
	templateCacheTime = time.Now()
	templateCacheMu.Unlock()

	return templates, nil
}

func InvalidateTemplateCache() {
	templateCacheMu.Lock()
	templateCache = nil
	templateCacheTime = time.Time{}
	templateCacheMu.Unlock()
}

func CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error) {
	var resp models.Template
	payload := map[string]interface{}{
		"name":             name,
		"title":            title,
		"description":      description,
		"company_name":     company,
		"default_solution": solution,
		"is_public":        isPublic,
	}
	err := doRequest("POST", "/templates", payload, &resp)
	if err == nil {
		InvalidateTemplateCache()
	}
	return &resp, err
}

func ListTemplates(includeAll bool) ([]models.Template, error) {
	var resp []models.Template
	params := url.Values{}
	if includeAll {
		params.Set("include_all", "true")
	}
	err := doRequest("GET", "/templates?"+params.Encode(), nil, &resp)
	return resp, err
}

func UseTemplate(templateID int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", fmt.Sprintf("/templates/%d/use", templateID), nil, &resp)
	return &resp, err
}

func UpdateTemplate(templateID int, payload map[string]interface{}) (*models.Template, error) {
	var resp models.Template
	err := doRequest("PUT", fmt.Sprintf("/templates/%d", templateID), payload, &resp)
	if err == nil {
		InvalidateTemplateCache()
	}
	return &resp, err
}

func DeleteTemplate(templateID int) error {
	err := doRequest("DELETE", fmt.Sprintf("/templates/%d", templateID), nil, nil)
	if err == nil {
		InvalidateTemplateCache()
	}
	return err
}

func GetTemplateByName(name string) (models.Template, error) {
	templates, err := listTemplatesCached(false)
	if err != nil {
		return models.Template{}, err
	}

	for _, tmpl := range templates {
		if tmpl.Name == name {
			return tmpl, nil
		}
	}

	return models.Template{}, fmt.Errorf("шаблон %q не найден", name)
}

func GetTemplateByID(templateID int) (*models.Template, error) {
	var resp models.Template
	err := doRequest("GET", fmt.Sprintf("/templates/%d", templateID), nil, &resp)
	return &resp, err
}
