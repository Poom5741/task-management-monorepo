package task

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, t *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context, projectID string, filter *TaskListFilter) ([]*Task, int, error)
	Update(ctx context.Context, id string, input *UpdateTaskInput) (*Task, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, page, pageSize int) ([]*Task, int, error)
	ProjectExists(ctx context.Context, id string) (bool, error)
}

type Usecase interface {
	CreateTask(ctx context.Context, input *CreateTaskInput) (*Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	ListTasks(ctx context.Context, projectID string, filter *TaskListFilter) ([]*Task, int, error)
	UpdateTask(ctx context.Context, id string, input *UpdateTaskInput) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
	SearchTasks(ctx context.Context, query string, page, pageSize int) ([]*Task, int, error)
	AddDependency(ctx context.Context, taskID string, input *CreateDependencyInput) (*TaskDependency, error)
	RemoveDependency(ctx context.Context, taskID, dependencyID string) error
}
