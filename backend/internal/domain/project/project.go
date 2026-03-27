package project

import (
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
)

type Project struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Status               Status    `json:"status"`
	TaskCount            int       `json:"task_count"`
	CompletionPercentage float64   `json:"completion_percentage"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateProjectInput struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type UpdateProjectInput struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Status      *Status `json:"status" validate:"omitempty,oneof=active archived"`
}

type ProjectListFilter struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Search   string `form:"search" validate:"max=100"`
}

type ProjectStatistics struct {
	TotalTasks      int     `json:"total_tasks"`
	TodoTasks       int     `json:"todo_tasks"`
	InProgressTasks int     `json:"in_progress_tasks"`
	DoneTasks       int     `json:"done_tasks"`
	CancelledTasks  int     `json:"cancelled_tasks"`
	OverdueTasks    int     `json:"overdue_tasks"`
	CompletionRate  float64 `json:"completion_rate"`
}
