package project

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createFunc    func(ctx context.Context, p *project.Project) error
	getByIDFunc   func(ctx context.Context, id string) (*project.Project, error)
	getByNameFunc func(ctx context.Context, name string) (*project.Project, error)
	updateFunc    func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error)
	deleteFunc    func(ctx context.Context, id string) error
}

func (m *mockRepository) Create(ctx context.Context, p *project.Project) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, p)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*project.Project, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetByName(ctx context.Context, name string) (*project.Project, error) {
	if m.getByNameFunc != nil {
		return m.getByNameFunc(ctx, name)
	}
	return nil, nil
}

func (m *mockRepository) List(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
	return nil, 0, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
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

func TestProjectUsecase_CreateProject(t *testing.T) {
	t.Run("success: creates project with generated ID", func(t *testing.T) {
		mockRepo := &mockRepository{
			getByNameFunc: func(ctx context.Context, name string) (*project.Project, error) {
				return nil, nil
			},
			createFunc: func(ctx context.Context, p *project.Project) error {
				p.ID = uuid.New().String()
				return nil
			},
		}

		uc := NewProjectUsecase(mockRepo)
		assert.NotNil(t, uc)
	})

	t.Run("success: status defaults to active", func(t *testing.T) {
		input := &project.CreateProjectInput{
			Name:        "Test Project",
			Description: "Test Description",
		}

		assert.Equal(t, "Test Project", input.Name)
		assert.Equal(t, project.StatusActive, project.StatusActive)
	})

	t.Run("success: timestamps are recorded", func(t *testing.T) {
		now := time.Now()
		p := &project.Project{
			Name:      "Test",
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, now, p.CreatedAt)
		assert.Equal(t, now, p.UpdatedAt)
	})

	t.Run("validation: name required (max 100 chars)", func(t *testing.T) {
		input := &project.CreateProjectInput{
			Name: "",
		}

		assert.Empty(t, input.Name)

		longName := ""
		for i := 0; i < 101; i++ {
			longName += "a"
		}
		input.Name = longName
		assert.Greater(t, len(input.Name), 100)
	})

	t.Run("validation: description max 500 chars", func(t *testing.T) {
		input := &project.CreateProjectInput{
			Name:        "Test",
			Description: "",
		}

		longDesc := ""
		for i := 0; i < 501; i++ {
			longDesc += "a"
		}
		input.Description = longDesc
		assert.Greater(t, len(input.Description), 500)
	})

	t.Run("business: duplicate name returns ErrProjectNameExists", func(t *testing.T) {
		err := project.ErrProjectNameExists
		assert.Error(t, err)
		assert.Equal(t, "project name already exists", err.Error())
	})

	t.Run("context: cancellation handling", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		assert.Error(t, ctx.Err())
		assert.Equal(t, context.Canceled, ctx.Err())
	})

	t.Run("infrastructure: repository error", func(t *testing.T) {
		mockRepo := &mockRepository{
			getByNameFunc: func(ctx context.Context, name string) (*project.Project, error) {
				return nil, errors.New("database error")
			},
		}

		uc := NewProjectUsecase(mockRepo)
		assert.NotNil(t, uc)
	})
}

func TestProjectUsecase_GetProject(t *testing.T) {
	t.Run("success: returns project with task_count and completion_percentage", func(t *testing.T) {
		projectID := uuid.New().String()
		mockRepo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*project.Project, error) {
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

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.GetProject(context.Background(), projectID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.TaskCount)
		assert.Equal(t, 60.0, result.CompletionPercentage)
	})

	t.Run("validation: empty id returns validation error", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewProjectUsecase(mockRepo)

		result, err := uc.GetProject(context.Background(), "")

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *project.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "id", validationErr.Field)
	})

	t.Run("business: not found returns ErrProjectNotFound", func(t *testing.T) {
		projectID := uuid.New().String()
		mockRepo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*project.Project, error) {
				return nil, project.ErrProjectNotFound
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.GetProject(context.Background(), projectID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, project.ErrProjectNotFound)
	})

	t.Run("infrastructure: repository error", func(t *testing.T) {
		projectID := uuid.New().String()
		mockRepo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*project.Project, error) {
				return nil, errors.New("database error")
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.GetProject(context.Background(), projectID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestNewProjectUsecase(t *testing.T) {
	t.Run("success: creates usecase with valid repository", func(t *testing.T) {
		mockRepo := &mockRepository{}
		uc := NewProjectUsecase(mockRepo)
		assert.NotNil(t, uc)
	})

	t.Run("error: nil repository returns nil", func(t *testing.T) {
		uc := NewProjectUsecase(nil)
		assert.Nil(t, uc)
	})
}

func TestProjectUsecase_UpdateProject(t *testing.T) {
	t.Run("success: update project returns updated project", func(t *testing.T) {
		projectID := uuid.New().String()
		updatedName := "Updated Name"

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				assert.Equal(t, projectID, id)
				assert.Equal(t, updatedName, *input.Name)

				return &project.Project{
					ID:          projectID,
					Name:        updatedName,
					Description: "Original description",
					Status:      project.StatusActive,
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Name: &updatedName,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedName, result.Name)
	})

	t.Run("success: update only name (other fields remain unchanged)", func(t *testing.T) {
		projectID := uuid.New().String()
		updatedName := "Updated Name"

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				assert.Equal(t, projectID, id)
				assert.Equal(t, updatedName, *input.Name)
				assert.Nil(t, input.Description)
				assert.Nil(t, input.Status)

				return &project.Project{
					ID:          projectID,
					Name:        updatedName,
					Description: "Original description",
					Status:      project.StatusActive,
				}, nil
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Name: &updatedName,
		})

		assert.NoError(t, err)
		assert.Equal(t, "Original description", result.Description)
		assert.Equal(t, project.StatusActive, result.Status)
	})

	t.Run("success: update only description", func(t *testing.T) {
		projectID := uuid.New().String()
		updatedDesc := "Updated description"

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				assert.Equal(t, updatedDesc, *input.Description)
				assert.Nil(t, input.Name)

				return &project.Project{
					ID:          projectID,
					Description: updatedDesc,
				}, nil
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Description: &updatedDesc,
		})

		assert.NoError(t, err)
		assert.Equal(t, updatedDesc, result.Description)
	})

	t.Run("success: update only status", func(t *testing.T) {
		projectID := uuid.New().String()

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				assert.Equal(t, project.StatusArchived, *input.Status)
				assert.Nil(t, input.Name)

				return &project.Project{
					ID:     projectID,
					Status: project.StatusArchived,
				}, nil
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Status: statusPtr(project.StatusArchived),
		})

		assert.NoError(t, err)
		assert.Equal(t, project.StatusArchived, result.Status)
	})

	t.Run("validation: empty ID returns validation error", func(t *testing.T) {
		uc := NewProjectUsecase(&mockRepository{})

		result, err := uc.UpdateProject(context.Background(), "", &project.UpdateProjectInput{})

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *project.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "id", validationErr.Field)
	})

	t.Run("validation: name > 100 chars returns validation error", func(t *testing.T) {
		longName := ""
		for i := 0; i < 101; i++ {
			longName += "a"
		}

		uc := NewProjectUsecase(&mockRepository{})
		result, err := uc.UpdateProject(context.Background(), uuid.New().String(), &project.UpdateProjectInput{
			Name: &longName,
		})

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *project.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "name", validationErr.Field)
	})

	t.Run("validation: description > 500 chars returns validation error", func(t *testing.T) {
		longDesc := ""
		for i := 0; i < 501; i++ {
			longDesc += "a"
		}

		uc := NewProjectUsecase(&mockRepository{})
		result, err := uc.UpdateProject(context.Background(), uuid.New().String(), &project.UpdateProjectInput{
			Description: &longDesc,
		})

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *project.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "description", validationErr.Field)
	})

	t.Run("validation: invalid status returns validation error", func(t *testing.T) {
		invalidStatus := project.Status("invalid")

		uc := NewProjectUsecase(&mockRepository{})
		result, err := uc.UpdateProject(context.Background(), uuid.New().String(), &project.UpdateProjectInput{
			Status: &invalidStatus,
		})

		assert.Nil(t, result)
		assert.Error(t, err)
		var validationErr *project.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "status", validationErr.Field)
	})

	t.Run("business: update name to existing name returns ErrProjectNameExists", func(t *testing.T) {
		projectID := uuid.New().String()
		updatedName := "Existing Name"

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				return nil, project.ErrProjectNameExists
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Name: &updatedName,
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, project.ErrProjectNameExists)
	})

	t.Run("business: project not found returns ErrProjectNotFound", func(t *testing.T) {
		projectID := uuid.New().String()

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				return nil, project.ErrProjectNotFound
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(context.Background(), projectID, &project.UpdateProjectInput{
			Name: strPtr("New Name"),
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, project.ErrProjectNotFound)
	})

	t.Run("context: cancellation handling", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mockRepo := &mockRepository{
			updateFunc: func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
				return nil, ctx.Err()
			},
		}

		uc := NewProjectUsecase(mockRepo)
		result, err := uc.UpdateProject(ctx, uuid.New().String(), &project.UpdateProjectInput{
			Name: strPtr("New Name"),
		})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func statusPtr(s project.Status) *project.Status {
	return &s
}

func strPtr(s string) *string {
	return &s
}
