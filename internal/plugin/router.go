package plugin

import (
	"tracker/internal/client"
	"tracker/internal/models"
)

type Router struct {
	nativeClient client.Client
	pluginClient TrackerPlugin
	usePlugin    bool
}

func NewRouter(native client.Client) *Router {
	return &Router{
		nativeClient: native,
		usePlugin:    false,
	}
}

func (r *Router) AttachPlugin(p TrackerPlugin) {
	r.pluginClient = p
	r.usePlugin = true
}

func (r *Router) DetachPlugin() {
	r.pluginClient = nil
	r.usePlugin = false
}

func (r *Router) IsPluginActive() bool {
	return r.usePlugin && r.pluginClient != nil
}

func (r *Router) ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListTasks(params, limit, offset)
	}
	return r.nativeClient.ListTasks(params, limit, offset)
}

func (r *Router) GetTaskByTicket(ticket string) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.GetTaskByTicket(ticket)
	}
	return r.nativeClient.GetTaskByTicket(ticket)
}

func (r *Router) CreateTask(payload map[string]interface{}) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.CreateTask(payload)
	}
	return r.nativeClient.CreateTask(payload)
}

func (r *Router) UpdateTask(id int, payload map[string]interface{}) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.UpdateTask(id, payload)
	}
	return r.nativeClient.UpdateTask(id, payload)
}

func (r *Router) DeleteTask(id int) error {
	if r.IsPluginActive() {
		return r.pluginClient.DeleteTask(id)
	}
	return r.nativeClient.DeleteTask(id)
}

func (r *Router) PauseTask(id int, payload map[string]interface{}) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.PauseTask(id, payload)
	}
	return r.nativeClient.PauseTask(id, payload)
}

func (r *Router) ResumeTask(id int, payload map[string]interface{}) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ResumeTask(id, payload)
	}
	return r.nativeClient.ResumeTask(id, payload)
}

func (r *Router) BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error) {
	if r.IsPluginActive() {
		return r.pluginClient.BulkCloseTasks(taskIDs)
	}
	return r.nativeClient.BulkCloseTasks(taskIDs)
}

func (r *Router) BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error) {
	if r.IsPluginActive() {
		return r.pluginClient.BulkAssignTasks(taskIDs, assignee)
	}
	return r.nativeClient.BulkAssignTasks(taskIDs, assignee)
}

func (r *Router) BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error) {
	if r.IsPluginActive() {
		return r.pluginClient.BulkDeleteTasks(taskIDs)
	}
	return r.nativeClient.BulkDeleteTasks(taskIDs)
}

func (r *Router) ListComments(taskID, limit, offset int) ([]models.Comment, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListComments(taskID, limit, offset)
	}
	return r.nativeClient.ListComments(taskID, limit, offset)
}

func (r *Router) CreateComment(taskID int, content string) (*models.Comment, error) {
	if r.IsPluginActive() {
		return r.pluginClient.CreateComment(taskID, content)
	}
	return r.nativeClient.CreateComment(taskID, content)
}

func (r *Router) UpdateComment(taskID, commentID int, content string) (*models.Comment, error) {
	if r.IsPluginActive() {
		return r.pluginClient.UpdateComment(taskID, commentID, content)
	}
	return r.nativeClient.UpdateComment(taskID, commentID, content)
}

func (r *Router) DeleteComment(taskID, commentID int) error {
	if r.IsPluginActive() {
		return r.pluginClient.DeleteComment(taskID, commentID)
	}
	return r.nativeClient.DeleteComment(taskID, commentID)
}

func (r *Router) ListCompanies(limit, offset int) (*models.CompanyListResponse, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListCompanies(limit, offset)
	}
	return r.nativeClient.ListCompanies(limit, offset)
}

func (r *Router) CreateCompany(name, description string) (*models.Company, error) {
	if r.IsPluginActive() {
		return r.pluginClient.CreateCompany(name, description)
	}
	return r.nativeClient.CreateCompany(name, description)
}

func (r *Router) DeleteCompany(name string) error {
	if r.IsPluginActive() {
		return r.pluginClient.DeleteCompany(name)
	}
	return r.nativeClient.DeleteCompany(name)
}

func (r *Router) ListUsers() ([]models.User, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListUsers()
	}
	return r.nativeClient.ListUsers()
}

func (r *Router) UpdateUserRole(username, role string) error {
	if r.IsPluginActive() {
		return r.pluginClient.UpdateUserRole(username, role)
	}
	return r.nativeClient.UpdateUserRole(username, role)
}

func (r *Router) ListTags(search string) ([]models.Tag, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListTags(search)
	}
	return r.nativeClient.ListTags(search)
}

func (r *Router) CreateTag(name, color string) (*models.Tag, error) {
	if r.IsPluginActive() {
		return r.pluginClient.CreateTag(name, color)
	}
	return r.nativeClient.CreateTag(name, color)
}

func (r *Router) UpdateTag(id int, payload map[string]interface{}) (*models.Tag, error) {
	if r.IsPluginActive() {
		return r.pluginClient.UpdateTag(id, payload)
	}
	return r.nativeClient.UpdateTag(id, payload)
}

func (r *Router) DeleteTag(id int) error {
	if r.IsPluginActive() {
		return r.pluginClient.DeleteTag(id)
	}
	return r.nativeClient.DeleteTag(id)
}

func (r *Router) ListTemplates(includeAll bool) ([]models.Template, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ListTemplates(includeAll)
	}
	return r.nativeClient.ListTemplates(includeAll)
}

func (r *Router) GetTemplateByID(id int) (*models.Template, error) {
	if r.IsPluginActive() {
		return r.pluginClient.GetTemplateByID(id)
	}
	return r.nativeClient.GetTemplateByID(id)
}

func (r *Router) GetTemplateByName(name string) (*models.Template, error) {
	if r.IsPluginActive() {
		return r.pluginClient.GetTemplateByName(name)
	}
	return r.nativeClient.GetTemplateByName(name)
}

func (r *Router) CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error) {
	if r.IsPluginActive() {
		return r.pluginClient.CreateTemplate(name, title, description, company, solution, isPublic)
	}
	return r.nativeClient.CreateTemplate(name, title, description, company, solution, isPublic)
}

func (r *Router) UpdateTemplate(id int, payload map[string]interface{}) (*models.Template, error) {
	if r.IsPluginActive() {
		return r.pluginClient.UpdateTemplate(id, payload)
	}
	return r.nativeClient.UpdateTemplate(id, payload)
}

func (r *Router) DeleteTemplate(id int) error {
	if r.IsPluginActive() {
		return r.pluginClient.DeleteTemplate(id)
	}
	return r.nativeClient.DeleteTemplate(id)
}

func (r *Router) UseTemplate(id int) (*models.Task, error) {
	if r.IsPluginActive() {
		return r.pluginClient.UseTemplate(id)
	}
	return r.nativeClient.UseTemplate(id)
}

func (r *Router) ExportTasks(params map[string]string) ([]byte, string, error) {
	if r.IsPluginActive() {
		return r.pluginClient.ExportTasks(params)
	}
	return r.nativeClient.ExportTasks(params)
}

func (r *Router) GetMe() (*models.User, error) {
	if r.IsPluginActive() {
		return r.pluginClient.GetMe()
	}
	return r.nativeClient.GetMe()
}
