package task

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/task"
	"github.com/stretchr/testify/assert"
)

type mockUsecase struct {
	createTaskFunc    func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error)
	getTaskFunc       func(ctx context.Context, id string) (*task.Task, error)
	listTasksFunc     func(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error)
	updateTaskFunc    func(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error)
	deleteTaskFunc    func(ctx context.Context, id string) error
	searchTasksFunc   func(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error)
	addDependencyFunc func(ctx context.Context, taskID string, input *task.CreateDependencyInput) (*task.TaskDependency, error)
	removeDepFunc     func(ctx context.Context, taskID, dependencyID string) error
}

func (m *mockUsecase) CreateTask(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, input)
	}
	return nil, nil
}

func (m *mockUsecase) GetTask(ctx context.Context, id string) (*task.Task, error) {
	if m.getTaskFunc != nil {
		return m.getTaskFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUsecase) ListTasks(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error) {
	if m.listTasksFunc != nil {
		return m.listTasksFunc(ctx, projectID, filter)
	}
	return nil, 0, nil
}

func (m *mockUsecase) UpdateTask(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error) {
	if m.updateTaskFunc != nil {
		return m.updateTaskFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *mockUsecase) DeleteTask(ctx context.Context, id string) error {
	if m.deleteTaskFunc != nil {
		return m.deleteTaskFunc(ctx, id)
	}
	return nil
}

func (m *mockUsecase) SearchTasks(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error) {
	if m.searchTasksFunc != nil {
		return m.searchTasksFunc(ctx, query, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockUsecase) AddDependency(ctx context.Context, taskID string, input *task.CreateDependencyInput) (*task.TaskDependency, error) {
	if m.addDependencyFunc != nil {
		return m.addDependencyFunc(ctx, taskID, input)
	}
	return nil, nil
}

func (m *mockUsecase) RemoveDependency(ctx context.Context, taskID, dependencyID string) error {
	if m.removeDepFunc != nil {
		return m.removeDepFunc(ctx, taskID, dependencyID)
	}
	return nil
}

func TestNewTaskHandler(t *testing.T) {
	t.Run("success: creates handler with valid usecase", func(t *testing.T) {
		mockUC := &mockUsecase{}
		handler := NewTaskHandler(mockUC)
		assert.NotNil(t, handler)
	})

	t.Run("error: nil usecase returns nil", func(t *testing.T) {
		handler := NewTaskHandler(nil)
		assert.Nil(t, handler)
	})
}

func TestTaskHandler_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success: POST /api/v1/projects/:projectId/tasks returns 201 with task", func(t *testing.T) {
		projectID := uuid.New().String()
		now := time.Now()

		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				assert.Equal(t, projectID, input.ProjectID)
				assert.Equal(t, "Test Task", input.Title)
				return &task.Task{
					ID:        uuid.New().String(),
					ProjectID: projectID,
					Title:     "Test Task",
					Status:    task.StatusTodo,
					Priority:  task.PriorityMedium,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		body := map[string]interface{}{
			"title": "Test Task",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Test Task", response["title"])
	})

	t.Run("success: creates task with all fields", func(t *testing.T) {
		projectID := uuid.New().String()
		dueDate := time.Now().Add(24 * time.Hour)

		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				assert.Equal(t, "Full Task", input.Title)
				assert.Equal(t, "Description", input.Description)
				assert.Equal(t, task.PriorityHigh, input.Priority)
				assert.NotNil(t, input.DueDate)
				assert.Equal(t, []string{"backend"}, input.Labels)
				return &task.Task{
					ID:          uuid.New().String(),
					ProjectID:   projectID,
					Title:       input.Title,
					Description: input.Description,
					Status:      task.StatusTodo,
					Priority:    input.Priority,
					DueDate:     input.DueDate,
					Labels:      input.Labels,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		body := map[string]interface{}{
			"title":       "Full Task",
			"description": "Description",
			"priority":    "high",
			"due_date":    dueDate.Format(time.RFC3339),
			"labels":      []string{"backend"},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("validation: empty projectId returns 400", func(t *testing.T) {
		mockUC := &mockUsecase{}
		handler := NewTaskHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectId", Value: ""}}
		c.Request = httptest.NewRequest(http.MethodPost, "/projects//tasks", nil)

		handler.CreateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: missing title returns 400", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, &task.ValidationError{
					Field:   "title",
					Message: "title is required",
				}
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		body := map[string]interface{}{
			"description": "Test Description",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: title > 200 chars returns 400", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, &task.ValidationError{
					Field:   "title",
					Message: "title must be at most 200 characters",
				}
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		longTitle := ""
		for i := 0; i < 201; i++ {
			longTitle += "a"
		}
		body := map[string]string{
			"title": longTitle,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: description > 2000 chars returns 400", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, &task.ValidationError{
					Field:   "description",
					Message: "description must be at most 2000 characters",
				}
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		longDesc := ""
		for i := 0; i < 2001; i++ {
			longDesc += "a"
		}
		body := map[string]string{
			"title":       "Test",
			"description": longDesc,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: past due_date returns 400", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, &task.ValidationError{
					Field:   "due_date",
					Message: "due date must be in the future",
				}
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		pastDate := time.Now().Add(-24 * time.Hour)
		body := map[string]interface{}{
			"title":    "Test",
			"due_date": pastDate.Format(time.RFC3339),
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: > 10 labels returns 400", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, &task.ValidationError{
					Field:   "labels",
					Message: "task cannot have more than 10 labels",
				}
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		labels := make([]string, 11)
		for i := 0; i < 11; i++ {
			labels[i] = "label"
		}
		body := map[string]interface{}{
			"title":  "Test",
			"labels": labels,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("business: project not found returns 404", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, task.ErrProjectNotFound
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		body := map[string]string{
			"title": "Test Task",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("infrastructure: internal error returns 500", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			createTaskFunc: func(ctx context.Context, input *task.CreateTaskInput) (*task.Task, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewTaskHandler(mockUC)
		router := gin.New()
		router.POST("/projects/:projectId/tasks", handler.CreateTask)

		body := map[string]string{
			"title": "Test Task",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID+"/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
