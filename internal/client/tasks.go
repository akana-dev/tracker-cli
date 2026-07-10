package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"tracker/internal/models"
)

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

func GetTaskWithComments(taskID int) (*models.Task, error) {
	var resp models.Task
	err := doRequest("GET", fmt.Sprintf("/tasks/%d", taskID), nil, &resp)
	return &resp, err
}
