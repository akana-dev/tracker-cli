package client

import "tracker/internal/models"

type NativeClient struct{}

func NewNativeClient() *NativeClient {
	return &NativeClient{}
}

func (c *NativeClient) ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error) {
	return ListTasks(params, limit, offset)
}

func (c *NativeClient) GetTaskByTicket(ticket string) (*models.Task, error) {
	return GetTaskByTicket(ticket)
}

func (c *NativeClient) CreateTask(payload map[string]interface{}) (*models.Task, error) {
	return CreateTask(payload)
}

func (c *NativeClient) UpdateTask(id int, payload map[string]interface{}) (*models.Task, error) {
	return UpdateTask(id, payload)
}

func (c *NativeClient) DeleteTask(id int) error {
	return DeleteTask(id)
}

func (c *NativeClient) PauseTask(id int, payload map[string]interface{}) (*models.Task, error) {
	return PauseTask(id, payload)
}

func (c *NativeClient) ResumeTask(id int, payload map[string]interface{}) (*models.Task, error) {
	return ResumeTask(id, payload)
}

func (c *NativeClient) BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error) {
	return BulkCloseTasks(taskIDs)
}

func (c *NativeClient) BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error) {
	return BulkAssignTasks(taskIDs, assignee)
}

func (c *NativeClient) BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error) {
	return BulkDeleteTasks(taskIDs)
}

func (c *NativeClient) ListComments(taskID, limit, offset int) ([]models.Comment, error) {
	return ListComments(taskID, limit, offset)
}

func (c *NativeClient) CreateComment(taskID int, content string) (*models.Comment, error) {
	return CreateComment(taskID, content)
}

func (c *NativeClient) UpdateComment(taskID, commentID int, content string) (*models.Comment, error) {
	return UpdateComment(taskID, commentID, content)
}

func (c *NativeClient) DeleteComment(taskID, commentID int) error {
	return DeleteComment(taskID, commentID)
}

func (c *NativeClient) ListCompanies(limit, offset int) (*models.CompanyListResponse, error) {
	return ListCompanies(limit, offset)
}

func (c *NativeClient) CreateCompany(name, description string) (*models.Company, error) {
	return CreateCompany(name, description)
}

func (c *NativeClient) DeleteCompany(name string) error {
	return DeleteCompany(name)
}

func (c *NativeClient) ListUsers() ([]models.User, error) {
	return ListUsers()
}

func (c *NativeClient) UpdateUserRole(username, role string) error {
	return UpdateUserRole(username, role)
}

func (c *NativeClient) ListTags(search string) ([]models.Tag, error) {
	return ListTags(search)
}

func (c *NativeClient) CreateTag(name, color string) (*models.Tag, error) {
	return CreateTag(name, color)
}

func (c *NativeClient) UpdateTag(id int, payload map[string]interface{}) (*models.Tag, error) {
	return UpdateTag(id, payload)
}

func (c *NativeClient) DeleteTag(id int) error {
	return DeleteTag(id)
}

func (c *NativeClient) ListTemplates(includeAll bool) ([]models.Template, error) {
	return ListTemplates(includeAll)
}

func (c *NativeClient) GetTemplateByID(id int) (*models.Template, error) {
	return GetTemplateByID(id)
}

func (c *NativeClient) GetTemplateByName(name string) (*models.Template, error) {
	tmpl, err := GetTemplateByName(name)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (c *NativeClient) CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error) {
	return CreateTemplate(name, title, description, company, solution, isPublic)
}

func (c *NativeClient) UpdateTemplate(id int, payload map[string]interface{}) (*models.Template, error) {
	return UpdateTemplate(id, payload)
}

func (c *NativeClient) DeleteTemplate(id int) error {
	return DeleteTemplate(id)
}

func (c *NativeClient) UseTemplate(id int) (*models.Task, error) {
	return UseTemplate(id)
}

func (c *NativeClient) ExportTasks(params map[string]string) ([]byte, string, error) {
	return ExportTasks(params)
}

func (c *NativeClient) GetMe() (*models.User, error) {
	return GetMe()
}
