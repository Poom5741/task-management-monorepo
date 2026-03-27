package task

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/task"
)

type TaskUsecase struct {
	repo task.Repository
}

func NewTaskUsecase(repo task.Repository) *TaskUsecase {
	if repo == nil {
		return nil
	}
	return &TaskUsecase{repo: repo}
}

func (uc *TaskUsecase) CreateTask(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	if input.Title == "" {
		return nil, &task.ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}

	if len(input.Title) > 200 {
		return nil, &task.ValidationError{
			Field:   "title",
			Message: "title must be at most 200 characters",
		}
	}

	if len(input.Description) > 2000 {
		return nil, &task.ValidationError{
			Field:   "description",
			Message: "description must be at most 2000 characters",
		}
	}

	if input.DueDate != nil && input.DueDate.Before(time.Now()) {
		return nil, &task.ValidationError{
			Field:   "due_date",
			Message: "due date must be in the future",
		}
	}

	if len(input.Labels) > 10 {
		return nil, &task.ValidationError{
			Field:   "labels",
			Message: "task cannot have more than 10 labels",
		}
	}

	exists, err := uc.repo.ProjectExists(ctx, input.ProjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, task.ErrProjectNotFound
	}

	priority := input.Priority
	if priority == "" {
		priority = task.PriorityMedium
	}

	now := time.Now()
	t := &task.Task{
		ID:          uuid.New().String(),
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      task.StatusTodo,
		Priority:    priority,
		DueDate:     input.DueDate,
		Labels:      input.Labels,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (uc *TaskUsecase) GetTask(ctx context.Context, id string) (*task.Task, error) {
	return nil, nil
}

func (uc *TaskUsecase) ListTasks(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error) {
	return nil, 0, nil
}

func (uc *TaskUsecase) UpdateTask(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error) {
	return nil, nil
}

func (uc *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
	return nil
}

func (uc *TaskUsecase) SearchTasks(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error) {
	return nil, 0, nil
}

func (uc *TaskUsecase) AddDependency(ctx context.Context, taskID string, input *task.CreateDependencyInput) (*task.TaskDependency, error) {
	return nil, nil
}

func (uc *TaskUsecase) RemoveDependency(ctx context.Context, taskID, dependencyID string) error {
	return nil
}
