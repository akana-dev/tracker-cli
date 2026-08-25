package plugin

import (
	"net/rpc"

	"tracker/internal/models"

	"github.com/hashicorp/go-plugin"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TRACKER_PLUGIN",
	MagicCookieValue: "tracker-cli-plugin-v1",
}

type TrackerPlugin interface {
	Name() string
	Version() string

	// Задачи
	ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error)
	GetTaskByTicket(ticket string) (*models.Task, error)
	CreateTask(payload map[string]interface{}) (*models.Task, error)
	UpdateTask(id int, payload map[string]interface{}) (*models.Task, error)
	DeleteTask(id int) error
	PauseTask(id int, payload map[string]interface{}) (*models.Task, error)
	ResumeTask(id int, payload map[string]interface{}) (*models.Task, error)

	// Массовые операции — ИСПОЛЬЗУЕМ models.BulkResponse
	BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error)
	BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error)
	BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error)

	// Комментарии
	ListComments(taskID, limit, offset int) ([]models.Comment, error)
	CreateComment(taskID int, content string) (*models.Comment, error)
	UpdateComment(taskID, commentID int, content string) (*models.Comment, error)
	DeleteComment(taskID, commentID int) error

	// Компании
	ListCompanies(limit, offset int) (*models.CompanyListResponse, error)
	CreateCompany(name, description string) (*models.Company, error)
	DeleteCompany(name string) error

	// Пользователи
	ListUsers() ([]models.User, error)
	UpdateUserRole(username, role string) error

	// Теги
	ListTags(search string) ([]models.Tag, error)
	CreateTag(name, color string) (*models.Tag, error)
	UpdateTag(id int, payload map[string]interface{}) (*models.Tag, error)
	DeleteTag(id int) error

	// Шаблоны
	ListTemplates(includeAll bool) ([]models.Template, error)
	GetTemplateByID(id int) (*models.Template, error)
	GetTemplateByName(name string) (*models.Template, error)
	CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error)
	UpdateTemplate(id int, payload map[string]interface{}) (*models.Template, error)
	DeleteTemplate(id int) error
	UseTemplate(id int) (*models.Task, error)

	// Экспорт
	ExportTasks(params map[string]string) ([]byte, string, error)

	// Аутентификация
	GetMe() (*models.User, error)

	// Конфигурация
	Configure(config map[string]string) error
}

type TrackerPluginPlugin struct {
	Impl TrackerPlugin
}

func (p *TrackerPluginPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TrackerPluginRPCServer{Impl: p.Impl}, nil
}

func (p *TrackerPluginPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TrackerPluginRPCClient{client: c}, nil
}
