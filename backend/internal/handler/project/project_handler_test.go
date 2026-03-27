package project

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
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
	"github.com/stretchr/testify/assert"
)

type mockUsecase struct {
	createProjectFunc func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error)
	getProjectFunc    func(ctx context.Context, id string) (*project.Project, error)
	listProjectsFunc  func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error)
	updateProjectFunc func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error)
	deleteProjectFunc func(ctx context.Context, id string) error
	statisticsFunc    func(ctx context.Context, id string) (*project.ProjectStatistics, error)
}

func (m *mockUsecase) CreateProject(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
	if m.createProjectFunc != nil {
		return m.createProjectFunc(ctx, input)
	}
	return nil, nil
}

func (m *mockUsecase) GetProject(ctx context.Context, id string) (*project.Project, error) {
	if m.getProjectFunc != nil {
		return m.getProjectFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUsecase) ListProjects(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
	if m.listProjectsFunc != nil {
		return m.listProjectsFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockUsecase) UpdateProject(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
	if m.updateProjectFunc != nil {
		return m.updateProjectFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *mockUsecase) DeleteProject(ctx context.Context, id string) error {
	if m.deleteProjectFunc != nil {
		return m.deleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *mockUsecase) GetProjectStatistics(ctx context.Context, id string) (*project.ProjectStatistics, error) {
	if m.statisticsFunc != nil {
		return m.statisticsFunc(ctx, id)
	}
	return nil, nil
}

func TestProjectHandler_CreateProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success: POST /api/v1/projects returns 201 with project", func(t *testing.T) {
		mockUC := &mockUsecase{
			createProjectFunc: func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
				return &project.Project{
					ID:          uuid.New().String(),
					Name:        input.Name,
					Description: input.Description,
					Status:      project.StatusActive,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		assert.NotNil(t, handler)

		router := gin.New()
		router.POST("/projects", handler.CreateProject)

		body := map[string]string{
			"name":        "Test Project",
			"description": "Test Description",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("validation: missing name returns 400", func(t *testing.T) {
		mockUC := &mockUsecase{
			createProjectFunc: func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
				return nil, &project.ValidationError{
					Field:   "name",
					Message: "name is required",
				}
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.POST("/projects", handler.CreateProject)

		body := map[string]string{
			"description": "Test Description",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation: name > 100 chars returns 400", func(t *testing.T) {
		mockUC := &mockUsecase{
			createProjectFunc: func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
				return nil, &project.ValidationError{
					Field:   "name",
					Message: "name must be at most 100 characters",
				}
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.POST("/projects", handler.CreateProject)

		longName := ""
		for i := 0; i < 101; i++ {
			longName += "a"
		}
		body := map[string]string{
			"name": longName,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("business: duplicate name returns 409", func(t *testing.T) {
		mockUC := &mockUsecase{
			createProjectFunc: func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
				return nil, project.ErrProjectNameExists
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.POST("/projects", handler.CreateProject)

		body := map[string]string{
			"name": "Existing Project",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("infrastructure: internal error returns 500", func(t *testing.T) {
		mockUC := &mockUsecase{
			createProjectFunc: func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.POST("/projects", handler.CreateProject)

		body := map[string]string{
			"name": "Test Project",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNewProjectHandler(t *testing.T) {
	t.Run("success: creates handler with valid usecase", func(t *testing.T) {
		mockUC := &mockUsecase{}
		handler := NewProjectHandler(mockUC)
		assert.NotNil(t, handler)
	})

	t.Run("error: nil usecase returns nil", func(t *testing.T) {
		handler := NewProjectHandler(nil)
		assert.Nil(t, handler)
	})
}

func TestProjectHandler_GetProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success: GET /api/v1/projects/:id returns 200 with task statistics", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			getProjectFunc: func(ctx context.Context, id string) (*project.Project, error) {
				assert.Equal(t, projectID, id)
				return &project.Project{
					ID:                   projectID,
					Name:                 "Test Project",
					Description:          "Test Description",
					Status:               project.StatusActive,
					TaskCount:            10,
					CompletionPercentage: 60.0,
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects/:id", handler.GetProject)

		req := httptest.NewRequest(http.MethodGet, "/projects/"+projectID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(10), response["task_count"])
		assert.Equal(t, float64(60), response["completion_percentage"])
	})

	t.Run("validation: empty id param returns 400", func(t *testing.T) {
		mockUC := &mockUsecase{}
		handler := NewProjectHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: ""}}
		c.Request = httptest.NewRequest(http.MethodGet, "/projects/", nil)

		handler.GetProject(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("business: not found returns 404", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			getProjectFunc: func(ctx context.Context, id string) (*project.Project, error) {
				return nil, project.ErrProjectNotFound
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects/:id", handler.GetProject)

		req := httptest.NewRequest(http.MethodGet, "/projects/"+projectID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("infrastructure: internal error returns 500", func(t *testing.T) {
		projectID := uuid.New().String()
		mockUC := &mockUsecase{
			getProjectFunc: func(ctx context.Context, id string) (*project.Project, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects/:id", handler.GetProject)

		req := httptest.NewRequest(http.MethodGet, "/projects/"+projectID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProjectHandler_ListProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success: GET /api/v1/projects returns list with default pagination", func(t *testing.T) {
		mockUC := &mockUsecase{
			listProjectsFunc: func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
				assert.Equal(t, 1, filter.Page)
				assert.Equal(t, 20, filter.PageSize)
				return []*project.Project{
					{
						ID:          uuid.New().String(),
						Name:        "Project 1",
						Description: "Description 1",
						Status:      project.StatusActive,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}, 1, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects", handler.ListProjects)

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success: custom page_size from query param", func(t *testing.T) {
		mockUC := &mockUsecase{
			listProjectsFunc: func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
				assert.Equal(t, 1, filter.Page)
				assert.Equal(t, 50, filter.PageSize)
				return []*project.Project{}, 0, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects", handler.ListProjects)

		req := httptest.NewRequest(http.MethodGet, "/projects?page_size=50", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success: page param parsed correctly", func(t *testing.T) {
		mockUC := &mockUsecase{
			listProjectsFunc: func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
				assert.Equal(t, 2, filter.Page)
				return []*project.Project{}, 0, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects", handler.ListProjects)

		req := httptest.NewRequest(http.MethodGet, "/projects?page=2", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success: search filter applied", func(t *testing.T) {
		mockUC := &mockUsecase{
			listProjectsFunc: func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
				assert.Equal(t, "test", filter.Search)
				return []*project.Project{}, 0, nil
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects", handler.ListProjects)

		req := httptest.NewRequest(http.MethodGet, "/projects?search=test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("infrastructure: internal error returns 500", func(t *testing.T) {
		mockUC := &mockUsecase{
			listProjectsFunc: func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
				return nil, 0, errors.New("database error")
			},
		}

		handler := NewProjectHandler(mockUC)
		router := gin.New()
		router.GET("/projects", handler.ListProjects)

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
