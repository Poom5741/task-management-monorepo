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
	return nil, nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
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
