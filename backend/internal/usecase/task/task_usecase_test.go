package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createFunc        func(ctx context.Context, tsk *task.Task) error
	getByIDFunc       func(ctx context.Context, id string) (*task.Task, error)
	listFunc          func(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error)
	updateFunc        func(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error)
	deleteFunc        func(ctx context.Context, id string) error
	searchFunc        func(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error)
	projectExistsFunc func(ctx context.Context, id string) (bool, error)
}

func (m *mockRepository) Create(ctx context.Context, tsk *task.Task) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, tsk)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*task.Task, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) List(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, projectID, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepository) Search(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockRepository) ProjectExists(ctx context.Context, id string) (bool, error) {
	if m.projectExistsFunc != nil {
		return m.projectExistsFunc(ctx, id)
	}
	return true, nil
}

func TestNewTaskUsecase(t *testing.T) {
	t.Run("success: creates usecase with valid repository", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)
		assert.NotNil(t, uc)
	})

	t.Run("error: nil repository returns nil", func(t *testing.T) {
		uc := NewTaskUsecase(nil)
		assert.Nil(t, uc)
	})
}

func TestTaskUsecase_CreateTask(t *testing.T) {
	t.Run("success: creates task with generated ID and defaults", func(t *testing.T) {
		projectID := uuid.New().String()
		now := time.Now()

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				assert.Equal(t, projectID, id)
				return true, nil
			},
			createFunc: func(ctx context.Context, tsk *task.Task) error {
				assert.NotEmpty(t, tsk.ID)
				assert.Equal(t, projectID, tsk.ProjectID)
				assert.Equal(t, "Test Task", tsk.Title)
				assert.Equal(t, task.StatusTodo, tsk.Status)
				assert.Equal(t, task.PriorityMedium, tsk.Priority)
				assert.NotZero(t, tsk.CreatedAt)
				assert.NotZero(t, tsk.UpdatedAt)
				return nil
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: projectID,
			Title:     "Test Task",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, task.StatusTodo, result.Status)
		assert.Equal(t, task.PriorityMedium, result.Priority)
		assert.True(t, result.CreatedAt.After(now.Add(-time.Second)))
	})

	t.Run("success: creates task with all fields", func(t *testing.T) {
		projectID := uuid.New().String()
		dueDate := time.Now().Add(24 * time.Hour)
		labels := []string{"backend", "urgent"}

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return true, nil
			},
			createFunc: func(ctx context.Context, tsk *task.Task) error {
				assert.Equal(t, "Full Task", tsk.Title)
				assert.Equal(t, "Description", tsk.Description)
				assert.Equal(t, task.PriorityHigh, tsk.Priority)
				assert.Equal(t, dueDate.Unix(), tsk.DueDate.Unix())
				assert.Equal(t, labels, tsk.Labels)
				return nil
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID:   projectID,
			Title:       "Full Task",
			Description: "Description",
			Priority:    task.PriorityHigh,
			DueDate:     &dueDate,
			Labels:      labels,
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("success: priority defaults to medium when empty", func(t *testing.T) {
		projectID := uuid.New().String()

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return true, nil
			},
			createFunc: func(ctx context.Context, tsk *task.Task) error {
				assert.Equal(t, task.PriorityMedium, tsk.Priority)
				return nil
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: projectID,
			Title:     "Test",
			Priority:  "",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, task.PriorityMedium, result.Priority)
	})

	t.Run("validation: nil input returns error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)

		result, err := uc.CreateTask(context.Background(), nil)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "input is required", err.Error())
	})

	t.Run("validation: empty title returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *task.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "title", validationErr.Field)
		assert.Equal(t, "title is required", validationErr.Message)
	})

	t.Run("validation: title > 200 chars returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)

		longTitle := ""
		for i := 0; i < 201; i++ {
			longTitle += "a"
		}
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     longTitle,
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *task.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "title", validationErr.Field)
	})

	t.Run("validation: description > 2000 chars returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)

		longDesc := ""
		for i := 0; i < 2001; i++ {
			longDesc += "a"
		}
		input := &task.CreateTaskInput{
			ProjectID:   uuid.New().String(),
			Title:       "Test",
			Description: longDesc,
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *task.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "description", validationErr.Field)
	})

	t.Run("validation: past due_date returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)

		pastDate := time.Now().Add(-24 * time.Hour)
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "Test",
			DueDate:   &pastDate,
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *task.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "due_date", validationErr.Field)
	})

	t.Run("validation: > 10 labels returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewTaskUsecase(mockRepo)

		labels := make([]string, 11)
		for i := 0; i < 11; i++ {
			labels[i] = "label" + string(rune('0'+i))
		}
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "Test",
			Labels:    labels,
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *task.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "labels", validationErr.Field)
	})

	t.Run("business: project not found returns ErrProjectNotFound", func(t *testing.T) {
		projectID := uuid.New().String()

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return false, nil
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: projectID,
			Title:     "Test",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, task.ErrProjectNotFound)
	})

	t.Run("business: status defaults to todo", func(t *testing.T) {
		projectID := uuid.New().String()

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return true, nil
			},
			createFunc: func(ctx context.Context, tsk *task.Task) error {
				assert.Equal(t, task.StatusTodo, tsk.Status)
				return nil
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: projectID,
			Title:     "Test",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, task.StatusTodo, result.Status)
	})

	t.Run("context: cancellation handling", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return false, ctx.Err()
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "Test",
		}

		result, err := uc.CreateTask(ctx, input)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("infrastructure: repository error on ProjectExists", func(t *testing.T) {
		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return false, errors.New("database error")
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "Test",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})

	t.Run("infrastructure: repository error on Create", func(t *testing.T) {
		mockRepo := &mockRepository{
			projectExistsFunc: func(ctx context.Context, id string) (bool, error) {
				return true, nil
			},
			createFunc: func(ctx context.Context, t *task.Task) error {
				return errors.New("insert error")
			},
		}

		uc := NewTaskUsecase(mockRepo)
		input := &task.CreateTaskInput{
			ProjectID: uuid.New().String(),
			Title:     "Test",
		}

		result, err := uc.CreateTask(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "insert error", err.Error())
	})
}
