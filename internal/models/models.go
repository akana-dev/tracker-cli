package models

import "fmt"

type User struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	FullName  *string      `json:"full_name,omitempty"`
	Role      string       `json:"role"`
	IsActive  bool         `json:"is_active"`
	CreatedAt FlexibleTime `json:"created_at"`
}

func (u *User) GetFullName() string {
	if u.FullName != nil && *u.FullName != "" {
		return *u.FullName
	}
	return u.Username
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Task struct {
	ID                 int           `json:"id"`
	Ticket             string        `json:"ticket"`
	Title              string        `json:"title"`
	StartTime          FlexibleTime  `json:"start_time"`
	EndTime            *FlexibleTime `json:"end_time,omitempty"`
	Comment            *string       `json:"comment,omitempty"`
	Solution           *string       `json:"solution,omitempty"`
	UserID             int           `json:"user_id"`
	OwnerUsername      string        `json:"owner_username"`
	OwnerFullName      *string       `json:"owner_full_name,omitempty"`
	AssigneeID         *int          `json:"assignee_id,omitempty"`
	AssigneeUsername   *string       `json:"assignee_username,omitempty"`
	AssigneeFullName   *string       `json:"assignee_full_name,omitempty"`
	CompanyName        string        `json:"company_name"`
	PausedAt           *FlexibleTime `json:"paused_at,omitempty"`
	TotalWorkedSeconds int           `json:"total_worked_seconds"`
	TotalHours         float64       `json:"total_hours"`
	Sessions           []TaskSession `json:"sessions"`
	CanEdit            bool          `json:"can_edit"`
	CanDelete          bool          `json:"can_delete"`
	Comments           []Comment     `json:"comments,omitempty"`
}

func (t *Task) GetAssigneeDisplay() string {
	if t.AssigneeUsername == nil || *t.AssigneeUsername == "" || *t.AssigneeUsername == t.OwnerUsername {
		return formatUserName(t.OwnerUsername, t.OwnerFullName)
	}
	return formatUserName(*t.AssigneeUsername, t.AssigneeFullName)
}

func (t *Task) GetOwnerDisplay() string {
	return formatUserName(t.OwnerUsername, t.OwnerFullName)
}

func (t *Task) IsAssignedToSomeone() bool {
	return t.AssigneeUsername != nil &&
		*t.AssigneeUsername != "" &&
		*t.AssigneeUsername != t.OwnerUsername
}

func (t *Task) IsClosed() bool {
	return t.EndTime != nil && !t.EndTime.IsZero()
}

func (t *Task) IsPaused() bool {
	return t.PausedAt != nil && !t.PausedAt.IsZero()
}

func (t *Task) IsActive() bool {
	return !t.IsClosed() && !t.IsPaused()
}

func formatUserName(username string, fullName *string) string {
	if fullName != nil && *fullName != "" && *fullName != username {
		return fmt.Sprintf("%s (%s)", *fullName, username)
	}
	return username
}

type TaskSession struct {
	ID              int           `json:"id"`
	StartTime       FlexibleTime  `json:"start_time"`
	EndTime         *FlexibleTime `json:"end_time,omitempty"`
	DurationSeconds int           `json:"duration_seconds"`
	DurationHours   float64       `json:"duration_hours"`
}

type Company struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Description    *string      `json:"description,omitempty"`
	NextTaskNumber int          `json:"next_task_number"`
	CreatedAt      FlexibleTime `json:"created_at"`
}

type Server struct {
	Name            string   `json:"name"`
	APIURL          string   `json:"api_url"`
	Token           string   `json:"token,omitempty"`
	UserRole        string   `json:"user_role,omitempty"`
	CurrentUsername string   `json:"current_username,omitempty"`
	AuthMethods     []string `json:"auth_methods"`
	ADDomain        string   `json:"ad_domain,omitempty"`
	DefaultCompany  string   `json:"default_company,omitempty"`
}

type ServersConfig struct {
	Current string             `json:"current"`
	Servers map[string]*Server `json:"servers"`
}

type TaskSummary struct {
	Total      int            `json:"total"`
	TotalHours float64        `json:"total_hours"`
	BySolution map[string]int `json:"by_solution"`
	ByCompany  map[string]int `json:"by_company"`
	ByAssignee map[string]int `json:"by_assignee"`
}

type TaskListResponse struct {
	Tasks  []Task `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (r *TaskListResponse) Pages() int {
	if r.Limit <= 0 {
		return 1
	}
	pages := r.Total / r.Limit
	if r.Total%r.Limit > 0 {
		pages++
	}
	if pages == 0 {
		pages = 1
	}
	return pages
}

func (r *TaskListResponse) CurrentPage() int {
	if r.Limit <= 0 {
		return 1
	}
	return (r.Offset / r.Limit) + 1
}

func (r *TaskListResponse) HasNext() bool {
	return r.Offset+len(r.Tasks) < r.Total
}

func (r *TaskListResponse) HasPrev() bool {
	return r.Offset > 0
}

type CompanyListResponse struct {
	Companies []Company `json:"items"`
	Total     int       `json:"total"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

func (r *CompanyListResponse) Pages() int {
	if r.Limit <= 0 {
		return 1
	}
	pages := r.Total / r.Limit
	if r.Total%r.Limit > 0 {
		pages++
	}
	if pages == 0 {
		pages = 1
	}
	return pages
}

func (r *CompanyListResponse) CurrentPage() int {
	if r.Limit <= 0 {
		return 1
	}
	return (r.Offset / r.Limit) + 1
}

func (r *CompanyListResponse) HasNext() bool {
	return r.Offset+len(r.Companies) < r.Total
}

func (r *CompanyListResponse) HasPrev() bool {
	return r.Offset > 0
}

type Template struct {
	// Name — имя шаблона (используется как идентификатор, не хранится в YAML)
	Name string `yaml:"-" json:"name"`

	// Title — название задачи (обязательное поле)
	Title string `yaml:"title" json:"title"`

	// Company — название компании (опционально)
	Company string `yaml:"company,omitempty" json:"company,omitempty"`

	// Assignee — исполнитель (опционально)
	Assignee string `yaml:"assignee,omitempty" json:"assignee,omitempty"`

	// Solution — статус решения по умолчанию (опционально)
	Solution string `yaml:"solution,omitempty" json:"solution,omitempty"`

	// Comment — комментарий по умолчанию (опционально)
	Comment string `yaml:"comment,omitempty" json:"comment,omitempty"`

	// Tags — список тегов (опционально)
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type Comment struct {
	ID               int           `json:"id"`
	TaskID           int           `json:"task_id"`
	Content          string        `json:"content"`
	ContentHTML      string        `json:"content_html"`
	User             CommentUser   `json:"user"`
	MentionedUserIDs []int         `json:"mentioned_user_ids"`
	CreatedAt        FlexibleTime  `json:"created_at"`
	UpdatedAt        *FlexibleTime `json:"updated_at,omitempty"`
	IsEdited         bool          `json:"is_edited"`
	CanEdit          bool          `json:"can_edit"`
	CanDelete        bool          `json:"can_delete"`
}

type CommentUser struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	FullName *string `json:"full_name,omitempty"`
}

func (u *CommentUser) GetDisplayName() string {
	if u.FullName != nil && *u.FullName != "" {
		return fmt.Sprintf("%s (%s)", *u.FullName, u.Username)
	}
	return u.Username
}
