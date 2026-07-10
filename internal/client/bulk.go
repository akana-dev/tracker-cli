package client

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
	payload := map[string]interface{}{"task_ids": taskIDs}
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
	payload := map[string]interface{}{"task_ids": taskIDs}
	err := doRequest("POST", "/tasks/bulk/delete", payload, &resp)
	return &resp, err
}
