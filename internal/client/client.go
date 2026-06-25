package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tracker/internal/config"
	"tracker/internal/models"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func doRequest(method, path string, body interface{}, result interface{}) error {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ошибка сериализации: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	token := config.LoadToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("сессия истекла. Выполните: tracker login")
	}
	if resp.StatusCode == 403 {
		return fmt.Errorf("доступ запрещён")
	}
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var errResp struct {
			Detail string `json:"detail"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Detail != "" {
			return fmt.Errorf("ошибка %d: %s", resp.StatusCode, errResp.Detail)
		}
		return fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func LoginPassword(username, password string) (*models.TokenResponse, error) {
	apiURL := config.GetAPIURL()
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", apiURL+"/auth/login", strings.NewReader(data.Encode()))
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

func CreateTask(payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", "/tasks", payload, &resp)
	return &resp, err
}

func ListTasks(params map[string]string) ([]models.Task, error) {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + "/tasks"

	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		fullURL += "?" + values.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	token := config.LoadToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tasks []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func UpdateTask(taskID int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("PUT", fmt.Sprintf("/tasks/%d", taskID), payload, &resp)
	return &resp, err
}

func DeleteTask(taskID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tasks/%d", taskID), nil, nil)
}

func PauseTask(taskID int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/pause", taskID), nil, &resp)
	return &resp, err
}

func ResumeTask(taskID int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/resume", taskID), nil, &resp)
	return &resp, err
}

func GetTaskSummary(params map[string]string) (*models.TaskSummary, error) {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + "/tasks/summary"

	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		fullURL += "?" + values.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	token := config.LoadToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var summary models.TaskSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func ListCompanies() ([]models.Company, error) {
	var resp []models.Company
	err := doRequest("GET", "/companies", nil, &resp)
	return resp, err
}

func CreateCompany(name, description string) (*models.Company, error) {
	payload := map[string]string{"name": name}
	if description != "" {
		payload["description"] = description
	}
	var resp models.Company
	err := doRequest("POST", "/companies", payload, &resp)
	return &resp, err
}

func DeleteCompany(name string) error {
	return doRequest("DELETE", fmt.Sprintf("/companies/%s", name), nil, nil)
}

func ListUsers() ([]models.User, error) {
	var resp []models.User
	err := doRequest("GET", "/tasks/users", nil, &resp)
	return resp, err
}

func UpdateUserRole(username, role string) error {
	payload := map[string]string{"role": role}
	return doRequest("PUT", fmt.Sprintf("/auth/users/%s/role", username), payload, nil)
}

func ExportTasks(format string, params map[string]string) ([]byte, string, error) {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + "/tasks/export"

	values := url.Values{}
	values.Set("format", format)
	for k, v := range params {
		values.Set(k, v)
	}
	fullURL += "?" + values.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, "", err
	}

	token := config.LoadToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytes))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	contentDisp := resp.Header.Get("Content-Disposition")
	filename := fmt.Sprintf("tasks.%s", format)
	if strings.Contains(contentDisp, "filename=") {
		parts := strings.Split(contentDisp, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], "\"")
		}
	}

	return data, filename, nil
}
