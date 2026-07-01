package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tracker/internal/config"
	"tracker/internal/models"
	"tracker/internal/service"
)

var httpClient = &http.Client{Timeout: service.HTTPTimeout}

func isRetryableError(err error, statusCode int) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if urlErr.Timeout() {
				return true
			}
			var opErr *net.OpError
			if errors.As(urlErr.Err, &opErr) {
				return true
			}
		}
		return false
	}

	switch statusCode {
	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusRequestTimeout,
		http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

func doRequestWithCtx(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ошибка сериализации: %w", err)
		}
	}

	var lastErr error
	backoff := service.InitialBackoff

	for attempt := 0; attempt <= service.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > service.MaxBackoff {
				backoff = service.MaxBackoff
			}
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
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
			lastErr = fmt.Errorf("ошибка сети: %w", err)
			if isRetryableError(err, 0) && attempt < service.MaxRetries {
				continue
			}
			return lastErr
		}

		if resp.StatusCode == 401 {
			resp.Body.Close()
			return fmt.Errorf("сессия истекла. Выполните: tracker login")
		}
		if resp.StatusCode == 403 {
			resp.Body.Close()
			return fmt.Errorf("доступ запрещён")
		}

		if resp.StatusCode >= 400 {
			bodyBytesResp, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if isRetryableError(nil, resp.StatusCode) && attempt < service.MaxRetries {
				lastErr = fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
				continue
			}

			var errResp struct {
				Detail string `json:"detail"`
			}
			if err := json.Unmarshal(bodyBytesResp, &errResp); err == nil && errResp.Detail != "" {
				return fmt.Errorf("ошибка %d: %s", resp.StatusCode, errResp.Detail)
			}
			return fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				resp.Body.Close()
				return err
			}
		}
		resp.Body.Close()
		return nil
	}

	return fmt.Errorf("превышено число попыток: %w", lastErr)
}

func doRequest(method, path string, body interface{}, result interface{}) error {
	return doRequestWithCtx(context.Background(), method, path, body, result)
}

func doRawRequestWithCtx(ctx context.Context, method, path string, body interface{}) ([]byte, http.Header, error) {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка сериализации: %w", err)
		}
	}

	var lastErr error
	backoff := service.InitialBackoff

	for attempt := 0; attempt <= service.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > service.MaxBackoff {
				backoff = service.MaxBackoff
			}
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return nil, nil, err
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
			lastErr = fmt.Errorf("ошибка сети: %w", err)
			if isRetryableError(err, 0) && attempt < service.MaxRetries {
				continue
			}
			return nil, nil, lastErr
		}

		if resp.StatusCode == 401 {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("сессия истекла. Выполните: tracker login")
		}
		if resp.StatusCode == 403 {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("доступ запрещён")
		}

		if resp.StatusCode >= 400 {
			bodyBytesResp, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if isRetryableError(nil, resp.StatusCode) && attempt < service.MaxRetries {
				lastErr = fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
				continue
			}

			var errResp struct {
				Detail string `json:"detail"`
			}
			if err := json.Unmarshal(bodyBytesResp, &errResp); err == nil && errResp.Detail != "" {
				return nil, nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, errResp.Detail)
			}
			return nil, nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		return data, resp.Header, nil
	}

	return nil, nil, fmt.Errorf("превышено число попыток: %w", lastErr)
}

func doRawRequest(method, path string, body interface{}) ([]byte, http.Header, error) {
	return doRawRequestWithCtx(context.Background(), method, path, body)
}

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

func CreateTask(payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", "/tasks", payload, &resp)
	return &resp, err
}

func GetTaskByID(id int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("GET", fmt.Sprintf("/tasks/%d", id), nil, &resp)
	return &resp, err
}

func GetTaskByTicket(ticket string) (*models.Task, error) {
	params := map[string]string{
		"ticket": ticket,
		"limit":  "1",
	}

	resp, err := ListTasks(params, 1, 0)
	if err != nil {
		return nil, err
	}

	if len(resp.Tasks) == 0 {
		return nil, fmt.Errorf("тикет %s не найден", ticket)
	}

	return GetTaskWithComments(resp.Tasks[0].ID)
}

func ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error) {
	values := url.Values{}
	for k, v := range params {
		if k == "limit" || k == "offset" {
			continue
		}
		values.Set(k, v)
	}

	if limit > 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		values.Set("offset", fmt.Sprintf("%d", offset))
	}

	path := "/tasks"
	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	data, _, err := doRawRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if len(trimmed) == 0 {
		return &models.TaskListResponse{}, nil
	}

	var resp models.TaskListResponse

	switch trimmed[0] {
	case '[':
		var tasks []models.Task
		if err := json.Unmarshal(data, &tasks); err != nil {
			return nil, fmt.Errorf("ошибка парсинга массива задач: %w", err)
		}
		resp.Tasks = tasks
		resp.Total = len(tasks)
		resp.Limit = limit
		resp.Offset = offset
	case '{':
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("ошибка парсинга структуры задач: %w", err)
		}
		if resp.Total == 0 {
			resp.Total = len(resp.Tasks)
		}
		if resp.Limit == 0 && limit > 0 {
			resp.Limit = limit
		}
		if resp.Offset == 0 && offset > 0 {
			resp.Offset = offset
		}
	default:
		return nil, fmt.Errorf("неожиданный формат ответа сервера")
	}

	return &resp, nil
}

func UpdateTask(taskID int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("PUT", fmt.Sprintf("/tasks/%d", taskID), payload, &resp)
	return &resp, err
}

func DeleteTask(taskID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tasks/%d", taskID), nil, nil)
}

func PauseTask(taskID int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/pause", taskID), payload, &resp)
	return &resp, err
}

func ResumeTask(taskID int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/resume", taskID), payload, &resp)
	return &resp, err
}

func GetTaskSummary(params map[string]string) (*models.TaskSummary, error) {
	path := "/tasks/summary"
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		path += "?" + values.Encode()
	}

	var resp models.TaskSummary
	err := doRequest("GET", path, nil, &resp)
	return &resp, err
}

func ExportTasks(params map[string]string) ([]byte, string, error) {
	values := url.Values{}

	for k, v := range params {
		if v != "" {
			values.Set(k, v)
		}
	}

	path := "/tasks/export?" + values.Encode()

	data, headers, err := doRawRequest("GET", path, nil)
	if err != nil {
		return nil, "", err
	}

	filename := "tasks.csv"
	if contentDisp := headers.Get("Content-Disposition"); strings.Contains(contentDisp, "filename=") {
		parts := strings.Split(contentDisp, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], "\"")
		}
	}

	return data, filename, nil
}

func ListCompanies(limit, offset int) (*models.CompanyListResponse, error) {
	values := url.Values{}
	if limit > 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		values.Set("offset", fmt.Sprintf("%d", offset))
	}

	path := "/companies"
	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	data, _, err := doRawRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if len(trimmed) == 0 {
		return &models.CompanyListResponse{}, nil
	}

	var resp models.CompanyListResponse

	if trimmed[0] == '[' {
		var companies []models.Company
		if err := json.Unmarshal(data, &companies); err != nil {
			return nil, fmt.Errorf("ошибка парсинга массива компаний: %w", err)
		}
		resp.Companies = companies
		resp.Total = len(companies)
		resp.Limit = limit
		resp.Offset = offset
	} else if trimmed[0] == '{' {
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("ошибка парсинга структуры компаний: %w", err)
		}
		if resp.Total == 0 {
			resp.Total = len(resp.Companies)
		}
		if resp.Limit == 0 && limit > 0 {
			resp.Limit = limit
		}
		if resp.Offset == 0 && offset > 0 {
			resp.Offset = offset
		}
	} else {
		return nil, fmt.Errorf("неожиданный формат ответа сервера")
	}

	return &resp, nil
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
	return doRequest("DELETE", fmt.Sprintf("/companies/%s", url.PathEscape(name)), nil, nil)
}

func ListUsers() ([]models.User, error) {
	var resp []models.User
	err := doRequest("GET", "/users", nil, &resp)
	return resp, err
}

func UpdateUserRole(username, role string) error {
	payload := map[string]string{"role": role}
	return doRequest("PUT", fmt.Sprintf("/users/%s/role", url.PathEscape(username)), payload, nil)
}

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
	payload := map[string]string{
		"content": content,
	}

	var resp models.Comment
	err := doRequest("POST", fmt.Sprintf("/tasks/%d/comments", taskID), payload, &resp)
	return &resp, err
}

func UpdateComment(taskID, commentID int, content string) (*models.Comment, error) {
	payload := map[string]string{
		"content": content,
	}

	var resp models.Comment
	err := doRequest("PUT", fmt.Sprintf("/tasks/%d/comments/%d", taskID, commentID), payload, &resp)
	return &resp, err
}

func DeleteComment(taskID, commentID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tasks/%d/comments/%d", taskID, commentID), nil, nil)
}

func GetTaskWithComments(taskID int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("GET", fmt.Sprintf("/tasks/%d", taskID), nil, &resp)
	return &resp, err
}

type BulkResult struct {
	TaskID int    `json:"task_id"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type BulkResponse struct {
	Total     int          `json:"total"`
	Succeeded int          `json:"succeeded"`
	Failed    int          `json:"failed"`
	Results   []BulkResult `json:"results"`
}

func BulkCloseTasks(taskIDs []int) (*BulkResponse, error) {
	var resp BulkResponse
	payload := map[string]interface{}{
		"task_ids": taskIDs,
	}
	err := doRequest("POST", "/tasks/bulk/close", payload, &resp)
	return &resp, err
}

func BulkAssignTasks(taskIDs []int, assignee string) (*BulkResponse, error) {
	var resp BulkResponse
	payload := map[string]interface{}{
		"task_ids":          taskIDs,
		"assignee_username": assignee,
	}
	err := doRequest("POST", "/tasks/bulk/assign", payload, &resp)
	return &resp, err
}

func BulkDeleteTasks(taskIDs []int) (*BulkResponse, error) {
	var resp BulkResponse
	payload := map[string]interface{}{
		"task_ids": taskIDs,
	}
	err := doRequest("POST", "/tasks/bulk/delete", payload, &resp)
	return &resp, err
}

type Tag struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Color             string `json:"color,omitempty"`
	CreatedByID       int    `json:"created_by_id"`
	CreatedByUsername string `json:"created_by_username"`
	CreatedAt         string `json:"created_at"`
}

func CreateTag(name, color string) (*models.Tag, error) {
	var resp models.Tag
	payload := map[string]interface{}{
		"name": name,
	}
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
	if err == nil {
		return resp, nil
	}

	var respObj struct {
		Items []models.Tag `json:"items"`
	}
	err = doRequest("GET", "/tags?"+params.Encode(), nil, &respObj)
	if err != nil {
		return nil, err
	}

	return respObj.Items, nil
}

func UpdateTag(tagID int, payload map[string]interface{}) (*models.Tag, error) {
	var resp models.Tag
	err := doRequest("PUT", fmt.Sprintf("/tags/%d", tagID), payload, &resp)
	return &resp, err
}

func DeleteTag(tagID int) error {
	return doRequest("DELETE", fmt.Sprintf("/tags/%d", tagID), nil, nil)
}

type Template struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	CompanyName     string `json:"company_name,omitempty"`
	DefaultSolution string `json:"default_solution,omitempty"`
	IsPublic        bool   `json:"is_public"`
	OwnerID         int    `json:"owner_id"`
	OwnerUsername   string `json:"owner_username"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	CanEdit         bool   `json:"can_edit"`
	CanDelete       bool   `json:"can_delete"`
}

func CreateTemplate(name, title, description, company, solution string, isPublic bool) (*Template, error) {
	var resp Template
	payload := map[string]interface{}{
		"name":             name,
		"title":            title,
		"description":      description,
		"company_name":     company,
		"default_solution": solution,
		"is_public":        isPublic,
	}
	err := doRequest("POST", "/templates", payload, &resp)
	return &resp, err
}

func ListTemplates(includeAll bool) ([]Template, error) {
	var resp []Template
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

func UpdateTemplate(templateID int, payload map[string]interface{}) (*Template, error) {
	var resp Template
	err := doRequest("PUT", fmt.Sprintf("/templates/%d", templateID), payload, &resp)
	return &resp, err
}

func DeleteTemplate(templateID int) error {
	return doRequest("DELETE", fmt.Sprintf("/templates/%d", templateID), nil, nil)
}
