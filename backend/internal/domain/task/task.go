package task

import (
	"time"
)

type Status string
type Priority string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"

	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type Task struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskInput struct {
	ProjectID   string     `json:"project_id" validate:"required"`
	Title       string     `json:"title" validate:"required,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	Priority    Priority   `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date" validate:"omitempty"`
	Labels      []string   `json:"labels" validate:"max=10"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title" validate:"omitempty,max=200"`
	Description *string    `json:"description" validate:"omitempty,max=2000"`
	Status      *Status    `json:"status" validate:"omitempty,oneof=todo in-progress done cancelled"`
	Priority    *Priority  `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date" validate:"omitempty"`
	Labels      []string   `json:"labels" validate:"omitempty,max=10"`
}

type TaskListFilter struct {
	Page      int        `form:"page" validate:"min=1"`
	PageSize  int        `form:"page_size" validate:"min=1,max=100"`
	Status    []Status   `form:"status" validate:"omitempty"`
	Priority  []Priority `form:"priority" validate:"omitempty"`
	Labels    []string   `form:"labels" validate:"omitempty"`
	Overdue   *bool      `form:"overdue" validate:"omitempty"`
	SortBy    string     `form:"sort_by" validate:"omitempty,oneof=created due priority"`
	SortOrder string     `form:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type TaskDependency struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	DependsOnID string `json:"depends_on_id"`
}

type CreateDependencyInput struct {
	DependsOnID string `json:"depends_on_id" validate:"required"`
}
