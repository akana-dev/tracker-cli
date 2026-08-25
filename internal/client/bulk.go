package client

import "tracker/internal/models"

func BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	payload := map[string]interface{}{"task_ids": taskIDs}
	err := doRequest("POST", "/tasks/bulk/close", payload, &resp)
	return &resp, err
}

func BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	payload := map[string]interface{}{
		"task_ids":          taskIDs,
		"assignee_username": assignee,
	}
	err := doRequest("POST", "/tasks/bulk/assign", payload, &resp)
	return &resp, err
}

func BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	payload := map[string]interface{}{"task_ids": taskIDs}
	err := doRequest("POST", "/tasks/bulk/delete", payload, &resp)
	return &resp, err
}
