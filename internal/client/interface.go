package client

import "tracker/internal/models"

type Client interface {
	ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error)
	GetTaskByTicket(ticket string) (*models.Task, error)
	CreateTask(payload map[string]interface{}) (*models.Task, error)
	UpdateTask(id int, payload map[string]interface{}) (*models.Task, error)
	DeleteTask(id int) error
	PauseTask(id int, payload map[string]interface{}) (*models.Task, error)
	ResumeTask(id int, payload map[string]interface{}) (*models.Task, error)

	BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error)
	BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error)
	BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error)

	ListComments(taskID, limit, offset int) ([]models.Comment, error)
	CreateComment(taskID int, content string) (*models.Comment, error)
	UpdateComment(taskID, commentID int, content string) (*models.Comment, error)
	DeleteComment(taskID, commentID int) error

	ListCompanies(limit, offset int) (*models.CompanyListResponse, error)
	CreateCompany(name, description string) (*models.Company, error)
	DeleteCompany(name string) error

	ListUsers() ([]models.User, error)
	UpdateUserRole(username, role string) error

	ListTags(search string) ([]models.Tag, error)
	CreateTag(name, color string) (*models.Tag, error)
	UpdateTag(id int, payload map[string]interface{}) (*models.Tag, error)
	DeleteTag(id int) error

	ListTemplates(includeAll bool) ([]models.Template, error)
	GetTemplateByID(id int) (*models.Template, error)
	GetTemplateByName(name string) (*models.Template, error)
	CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error)
	UpdateTemplate(id int, payload map[string]interface{}) (*models.Template, error)
	DeleteTemplate(id int) error
	UseTemplate(id int) (*models.Task, error)

	ExportTasks(params map[string]string) ([]byte, string, error)

	GetMe() (*models.User, error)
}
