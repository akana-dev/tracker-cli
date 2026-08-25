package plugin

import (
	"net/rpc"

	"tracker/internal/models"
)

type TrackerPluginRPCServer struct {
	Impl TrackerPlugin
}

func (s *TrackerPluginRPCServer) Name(args interface{}, resp *string) error {
	*resp = s.Impl.Name()
	return nil
}

func (s *TrackerPluginRPCServer) Version(args interface{}, resp *string) error {
	*resp = s.Impl.Version()
	return nil
}

func (s *TrackerPluginRPCServer) ListTasks(args ListTasksArgs, resp *models.TaskListResponse) error {
	result, err := s.Impl.ListTasks(args.Params, args.Limit, args.Offset)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) GetTaskByTicket(args GetTaskByTicketArgs, resp *models.Task) error {
	result, err := s.Impl.GetTaskByTicket(args.Ticket)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) CreateTask(args CreateTaskArgs, resp *models.Task) error {
	result, err := s.Impl.CreateTask(args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) UpdateTask(args UpdateTaskArgs, resp *models.Task) error {
	result, err := s.Impl.UpdateTask(args.ID, args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) DeleteTask(args DeleteTaskArgs, resp *interface{}) error {
	return s.Impl.DeleteTask(args.ID)
}

func (s *TrackerPluginRPCServer) PauseTask(args PauseTaskArgs, resp *models.Task) error {
	result, err := s.Impl.PauseTask(args.ID, args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) ResumeTask(args ResumeTaskArgs, resp *models.Task) error {
	result, err := s.Impl.ResumeTask(args.ID, args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) BulkCloseTasks(args BulkCloseTasksArgs, resp *models.BulkResponse) error {
	result, err := s.Impl.BulkCloseTasks(args.TaskIDs)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) BulkAssignTasks(args BulkAssignTasksArgs, resp *models.BulkResponse) error {
	result, err := s.Impl.BulkAssignTasks(args.TaskIDs, args.Assignee)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) BulkDeleteTasks(args BulkDeleteTasksArgs, resp *models.BulkResponse) error {
	result, err := s.Impl.BulkDeleteTasks(args.TaskIDs)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) ListComments(args ListCommentsArgs, resp *[]models.Comment) error {
	result, err := s.Impl.ListComments(args.TaskID, args.Limit, args.Offset)
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

func (s *TrackerPluginRPCServer) CreateComment(args CreateCommentArgs, resp *models.Comment) error {
	result, err := s.Impl.CreateComment(args.TaskID, args.Content)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) UpdateComment(args UpdateCommentArgs, resp *models.Comment) error {
	result, err := s.Impl.UpdateComment(args.TaskID, args.CommentID, args.Content)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) DeleteComment(args DeleteCommentArgs, resp *interface{}) error {
	return s.Impl.DeleteComment(args.TaskID, args.CommentID)
}

func (s *TrackerPluginRPCServer) ListCompanies(args ListCompaniesArgs, resp *models.CompanyListResponse) error {
	result, err := s.Impl.ListCompanies(args.Limit, args.Offset)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) CreateCompany(args CreateCompanyArgs, resp *models.Company) error {
	result, err := s.Impl.CreateCompany(args.Name, args.Description)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) DeleteCompany(args DeleteCompanyArgs, resp *interface{}) error {
	return s.Impl.DeleteCompany(args.Name)
}

func (s *TrackerPluginRPCServer) ListUsers(args interface{}, resp *[]models.User) error {
	result, err := s.Impl.ListUsers()
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

func (s *TrackerPluginRPCServer) UpdateUserRole(args UpdateUserRoleArgs, resp *interface{}) error {
	return s.Impl.UpdateUserRole(args.Username, args.Role)
}

func (s *TrackerPluginRPCServer) ListTags(args ListTagsArgs, resp *[]models.Tag) error {
	result, err := s.Impl.ListTags(args.Search)
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

func (s *TrackerPluginRPCServer) CreateTag(args CreateTagArgs, resp *models.Tag) error {
	result, err := s.Impl.CreateTag(args.Name, args.Color)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) UpdateTag(args UpdateTagArgs, resp *models.Tag) error {
	result, err := s.Impl.UpdateTag(args.ID, args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) DeleteTag(args DeleteTagArgs, resp *interface{}) error {
	return s.Impl.DeleteTag(args.ID)
}

func (s *TrackerPluginRPCServer) ListTemplates(args ListTemplatesArgs, resp *[]models.Template) error {
	result, err := s.Impl.ListTemplates(args.IncludeAll)
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

func (s *TrackerPluginRPCServer) GetTemplateByID(args GetTemplateByIDArgs, resp *models.Template) error {
	result, err := s.Impl.GetTemplateByID(args.ID)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) GetTemplateByName(args GetTemplateByNameArgs, resp *models.Template) error {
	result, err := s.Impl.GetTemplateByName(args.Name)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) CreateTemplate(args CreateTemplateArgs, resp *models.Template) error {
	result, err := s.Impl.CreateTemplate(args.Name, args.Title, args.Description, args.Company, args.Solution, args.IsPublic)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) UpdateTemplate(args UpdateTemplateArgs, resp *models.Template) error {
	result, err := s.Impl.UpdateTemplate(args.ID, args.Payload)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) DeleteTemplate(args DeleteTemplateArgs, resp *interface{}) error {
	return s.Impl.DeleteTemplate(args.ID)
}

func (s *TrackerPluginRPCServer) UseTemplate(args UseTemplateArgs, resp *models.Task) error {
	result, err := s.Impl.UseTemplate(args.ID)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) ExportTasks(args ExportTasksArgs, resp *ExportTasksResult) error {
	data, filename, err := s.Impl.ExportTasks(args.Params)
	if err != nil {
		return err
	}
	resp.Data = data
	resp.Filename = filename
	return nil
}

func (s *TrackerPluginRPCServer) GetMe(args interface{}, resp *models.User) error {
	result, err := s.Impl.GetMe()
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

func (s *TrackerPluginRPCServer) Configure(args ConfigureArgs, resp *interface{}) error {
	return s.Impl.Configure(args.Config)
}

type TrackerPluginRPCClient struct {
	client *rpc.Client
}

func (c *TrackerPluginRPCClient) Name() (string, error) {
	var resp string
	err := c.client.Call("Plugin.Name", new(interface{}), &resp)
	return resp, err
}

func (c *TrackerPluginRPCClient) Version() (string, error) {
	var resp string
	err := c.client.Call("Plugin.Version", new(interface{}), &resp)
	return resp, err
}

func (c *TrackerPluginRPCClient) ListTasks(params map[string]string, limit, offset int) (*models.TaskListResponse, error) {
	var resp models.TaskListResponse
	err := c.client.Call("Plugin.ListTasks", ListTasksArgs{Params: params, Limit: limit, Offset: offset}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) GetTaskByTicket(ticket string) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.GetTaskByTicket", GetTaskByTicketArgs{Ticket: ticket}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) CreateTask(payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.CreateTask", CreateTaskArgs{Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) UpdateTask(id int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.UpdateTask", UpdateTaskArgs{ID: id, Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) DeleteTask(id int) error {
	return c.client.Call("Plugin.DeleteTask", DeleteTaskArgs{ID: id}, new(interface{}))
}

func (c *TrackerPluginRPCClient) PauseTask(id int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.PauseTask", PauseTaskArgs{ID: id, Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) ResumeTask(id int, payload map[string]interface{}) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.ResumeTask", ResumeTaskArgs{ID: id, Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) BulkCloseTasks(taskIDs []int) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	err := c.client.Call("Plugin.BulkCloseTasks", BulkCloseTasksArgs{TaskIDs: taskIDs}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) BulkAssignTasks(taskIDs []int, assignee string) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	err := c.client.Call("Plugin.BulkAssignTasks", BulkAssignTasksArgs{TaskIDs: taskIDs, Assignee: assignee}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) BulkDeleteTasks(taskIDs []int) (*models.BulkResponse, error) {
	var resp models.BulkResponse
	err := c.client.Call("Plugin.BulkDeleteTasks", BulkDeleteTasksArgs{TaskIDs: taskIDs}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) ListComments(taskID, limit, offset int) ([]models.Comment, error) {
	var resp []models.Comment
	err := c.client.Call("Plugin.ListComments", ListCommentsArgs{TaskID: taskID, Limit: limit, Offset: offset}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TrackerPluginRPCClient) CreateComment(taskID int, content string) (*models.Comment, error) {
	var resp models.Comment
	err := c.client.Call("Plugin.CreateComment", CreateCommentArgs{TaskID: taskID, Content: content}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) UpdateComment(taskID, commentID int, content string) (*models.Comment, error) {
	var resp models.Comment
	err := c.client.Call("Plugin.UpdateComment", UpdateCommentArgs{TaskID: taskID, CommentID: commentID, Content: content}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) DeleteComment(taskID, commentID int) error {
	return c.client.Call("Plugin.DeleteComment", DeleteCommentArgs{TaskID: taskID, CommentID: commentID}, new(interface{}))
}

func (c *TrackerPluginRPCClient) ListCompanies(limit, offset int) (*models.CompanyListResponse, error) {
	var resp models.CompanyListResponse
	err := c.client.Call("Plugin.ListCompanies", ListCompaniesArgs{Limit: limit, Offset: offset}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) CreateCompany(name, description string) (*models.Company, error) {
	var resp models.Company
	err := c.client.Call("Plugin.CreateCompany", CreateCompanyArgs{Name: name, Description: description}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) DeleteCompany(name string) error {
	return c.client.Call("Plugin.DeleteCompany", DeleteCompanyArgs{Name: name}, new(interface{}))
}

func (c *TrackerPluginRPCClient) ListUsers() ([]models.User, error) {
	var resp []models.User
	err := c.client.Call("Plugin.ListUsers", new(interface{}), &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TrackerPluginRPCClient) UpdateUserRole(username, role string) error {
	return c.client.Call("Plugin.UpdateUserRole", UpdateUserRoleArgs{Username: username, Role: role}, new(interface{}))
}

func (c *TrackerPluginRPCClient) ListTags(search string) ([]models.Tag, error) {
	var resp []models.Tag
	err := c.client.Call("Plugin.ListTags", ListTagsArgs{Search: search}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TrackerPluginRPCClient) CreateTag(name, color string) (*models.Tag, error) {
	var resp models.Tag
	err := c.client.Call("Plugin.CreateTag", CreateTagArgs{Name: name, Color: color}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) UpdateTag(id int, payload map[string]interface{}) (*models.Tag, error) {
	var resp models.Tag
	err := c.client.Call("Plugin.UpdateTag", UpdateTagArgs{ID: id, Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) DeleteTag(id int) error {
	return c.client.Call("Plugin.DeleteTag", DeleteTagArgs{ID: id}, new(interface{}))
}

func (c *TrackerPluginRPCClient) ListTemplates(includeAll bool) ([]models.Template, error) {
	var resp []models.Template
	err := c.client.Call("Plugin.ListTemplates", ListTemplatesArgs{IncludeAll: includeAll}, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TrackerPluginRPCClient) GetTemplateByID(id int) (*models.Template, error) {
	var resp models.Template
	err := c.client.Call("Plugin.GetTemplateByID", GetTemplateByIDArgs{ID: id}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) GetTemplateByName(name string) (*models.Template, error) {
	var resp models.Template
	err := c.client.Call("Plugin.GetTemplateByName", GetTemplateByNameArgs{Name: name}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) CreateTemplate(name, title, description, company, solution string, isPublic bool) (*models.Template, error) {
	var resp models.Template
	err := c.client.Call("Plugin.CreateTemplate", CreateTemplateArgs{
		Name: name, Title: title, Description: description,
		Company: company, Solution: solution, IsPublic: isPublic,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) UpdateTemplate(id int, payload map[string]interface{}) (*models.Template, error) {
	var resp models.Template
	err := c.client.Call("Plugin.UpdateTemplate", UpdateTemplateArgs{ID: id, Payload: payload}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) DeleteTemplate(id int) error {
	return c.client.Call("Plugin.DeleteTemplate", DeleteTemplateArgs{ID: id}, new(interface{}))
}

func (c *TrackerPluginRPCClient) UseTemplate(id int) (*models.Task, error) {
	var resp models.Task
	err := c.client.Call("Plugin.UseTemplate", UseTemplateArgs{ID: id}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) ExportTasks(params map[string]string) ([]byte, string, error) {
	var resp ExportTasksResult
	err := c.client.Call("Plugin.ExportTasks", ExportTasksArgs{Params: params}, &resp)
	if err != nil {
		return nil, "", err
	}
	return resp.Data, resp.Filename, nil
}

func (c *TrackerPluginRPCClient) GetMe() (*models.User, error) {
	var resp models.User
	err := c.client.Call("Plugin.GetMe", new(interface{}), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *TrackerPluginRPCClient) Configure(config map[string]string) error {
	return c.client.Call("Plugin.Configure", ConfigureArgs{Config: config}, new(interface{}))
}

type ListTasksArgs struct {
	Params map[string]string
	Limit  int
	Offset int
}

type GetTaskByTicketArgs struct {
	Ticket string
}

type CreateTaskArgs struct {
	Payload map[string]interface{}
}

type UpdateTaskArgs struct {
	ID      int
	Payload map[string]interface{}
}

type DeleteTaskArgs struct {
	ID int
}

type PauseTaskArgs struct {
	ID      int
	Payload map[string]interface{}
}

type ResumeTaskArgs struct {
	ID      int
	Payload map[string]interface{}
}

type BulkCloseTasksArgs struct {
	TaskIDs []int
}

type BulkAssignTasksArgs struct {
	TaskIDs  []int
	Assignee string
}

type BulkDeleteTasksArgs struct {
	TaskIDs []int
}

type ListCommentsArgs struct {
	TaskID int
	Limit  int
	Offset int
}

type CreateCommentArgs struct {
	TaskID  int
	Content string
}

type UpdateCommentArgs struct {
	TaskID    int
	CommentID int
	Content   string
}

type DeleteCommentArgs struct {
	TaskID    int
	CommentID int
}

type ListCompaniesArgs struct {
	Limit  int
	Offset int
}

type CreateCompanyArgs struct {
	Name        string
	Description string
}

type DeleteCompanyArgs struct {
	Name string
}

type UpdateUserRoleArgs struct {
	Username string
	Role     string
}

type ListTagsArgs struct {
	Search string
}

type CreateTagArgs struct {
	Name  string
	Color string
}

type UpdateTagArgs struct {
	ID      int
	Payload map[string]interface{}
}

type DeleteTagArgs struct {
	ID int
}

type ListTemplatesArgs struct {
	IncludeAll bool
}

type GetTemplateByIDArgs struct {
	ID int
}

type GetTemplateByNameArgs struct {
	Name string
}

type CreateTemplateArgs struct {
	Name        string
	Title       string
	Description string
	Company     string
	Solution    string
	IsPublic    bool
}

type UpdateTemplateArgs struct {
	ID      int
	Payload map[string]interface{}
}

type DeleteTemplateArgs struct {
	ID int
}

type UseTemplateArgs struct {
	ID int
}

type ExportTasksArgs struct {
	Params map[string]string
}

type ExportTasksResult struct {
	Data     []byte
	Filename string
}

type ConfigureArgs struct {
	Config map[string]string
}
