package models

type User struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Role      string       `json:"role"`
	IsActive  bool         `json:"is_active"`
	CreatedAt FlexibleTime `json:"created_at"`
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
	AssigneeID         *int          `json:"assignee_id,omitempty"`
	AssigneeUsername   *string       `json:"assignee_username,omitempty"`
	CompanyName        string        `json:"company_name"`
	PausedAt           *FlexibleTime `json:"paused_at,omitempty"`
	TotalWorkedSeconds int           `json:"total_worked_seconds"`
	TotalHours         float64       `json:"total_hours"`
	Sessions           []TaskSession `json:"sessions"`
	CanEdit            bool          `json:"can_edit"`
	CanDelete          bool          `json:"can_delete"`
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
