package client

import (
	"fmt"
	"net/url"

	"tracker/internal/models"
)

func ListUsers() ([]models.User, error) {
	var resp []models.User
	err := doRequest("GET", "/users", nil, &resp)
	return resp, err
}

func UpdateUserRole(username, role string) error {
	payload := map[string]string{"role": role}
	return doRequest("PUT", fmt.Sprintf("/users/%s/role", url.PathEscape(username)), payload, nil)
}
