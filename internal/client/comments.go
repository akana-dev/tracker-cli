package client

import (
	"fmt"
	"net/url"

	"tracker/internal/models"
)

func ListComments(taskID int, limit, offset int) ([]models.Comment, error) {
	path := fmt.Sprintf("/tasks/%d/comments", taskID)

	values := url.Values{}
	if limit > 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		values.Set("offset", fmt.Sprintf("%d", offset))
	}

	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	var resp []models.Comment
	err := doRequest("GET", path, nil, &resp)
	return resp, err
}

func CreateComment(taskID int, content string) (*models.Comment, error) {
	payload := map[string]string{"content": content}
	var resp models.Comment
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/comments", taskID), payload, &resp)
	return &resp, err
}

func UpdateComment(taskID, commentID int, content string) (*models.Comment, error) {
	payload := map[string]string{"content": content}
	var resp models.Comment
	err := doRequest("PUT", fmt.Sprintf("/tasks/%d/comments/%d", taskID, commentID), payload, &resp)
	return &resp, err
}

func DeleteComment(taskID, commentID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tasks/%d/comments/%d", taskID, commentID), nil, nil)
}
