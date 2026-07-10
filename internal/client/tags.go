package client

import (
	"fmt"
	"net/url"

	"tracker/internal/models"
)

func CreateTag(name, color string) (*models.Tag, error) {
	var resp models.Tag
	payload := map[string]interface{}{"name": name}
	if color != "" {
		payload["color"] = color
	}
	err := doRequest("POST", "/tags", payload, &resp)
	return &resp, err
}

func ListTags(search string) ([]models.Tag, error) {
	params := url.Values{}
	if search != "" {
		params.Set("search", search)
	}

	var resp []models.Tag
	err := doRequest("GET", "/tags?"+params.Encode(), nil, &resp)
	return resp, err
}

func UpdateTag(tagID int, payload map[string]interface{}) (*models.Tag, error) {
	var resp models.Tag
	err := doRequest("PUT", fmt.Sprintf("/tags/%d", tagID), payload, &resp)
	return &resp, err
}

func DeleteTag(tagID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tags/%d", tagID), nil, nil)
}
