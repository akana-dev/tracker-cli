package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"tracker/internal/config"
	"tracker/internal/models"
	"tracker/internal/service"
)

func LoginPassword(username, password string) (*models.TokenResponse, error) {
	apiURL := config.GetAPIURL()
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	ctx, cancel := context.WithTimeout(context.Background(), service.HTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL+"/auth/login", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка авторизации: %s", string(bodyBytes))
	}

	var tokenResp models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func LoginAD(username, password string) (*models.TokenResponse, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
		"method":   "password",
	}
	var resp models.TokenResponse
	err := doRequest("POST", "/auth/login/ad", payload, &resp)
	return &resp, err
}

func GetMe() (*models.User, error) {
	var resp models.User
	err := doRequest("GET", "/auth/me", nil, &resp)
	return &resp, err
}

func RegisterUser(username, email, password string) error {
	payload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}
	return doRequest("POST", "/auth/register", payload, nil)
}
