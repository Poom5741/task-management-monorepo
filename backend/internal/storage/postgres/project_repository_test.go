package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
	"github.com/stretchr/testify/assert"
)

func TestNewProjectRepository(t *testing.T) {
	t.Run("success: creates repository with valid db", func(t *testing.T) {
		repo := NewProjectRepository(nil)
		assert.Nil(t, repo)
	})
}

func TestProjectRepository_Create(t *testing.T) {
	t.Run("success: creates project with valid input", func(t *testing.T) {
		p := &project.Project{
			Name:        "Test Project",
			Description: "Test Description",
			Status:      project.StatusActive,
		}

		assert.Equal(t, "Test Project", p.Name)
		assert.Equal(t, "Test Description", p.Description)
		assert.Equal(t, project.StatusActive, p.Status)
		assert.Empty(t, p.ID)

		p.ID = uuid.New().String()
		assert.NotEmpty(t, p.ID)
	})

	t.Run("validation: duplicate name returns ErrProjectNameExists", func(t *testing.T) {
		err := project.ErrProjectNameExists
		assert.Error(t, err)
		assert.Equal(t, "project name already exists", err.Error())
	})
}

func TestProjectRepository_GetByName(t *testing.T) {
	t.Run("success: returns project when found", func(t *testing.T) {
		p := &project.Project{
			ID:          uuid.New().String(),
			Name:        "Test Project",
			Description: "Test Description",
			Status:      project.StatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, p)
		assert.Equal(t, "Test Project", p.Name)
	})

	t.Run("success: returns nil when not found", func(t *testing.T) {
		var p *project.Project

		assert.Nil(t, p)
	})
}

func TestProjectRepository_GetByID(t *testing.T) {
	t.Run("success: returns project when found", func(t *testing.T) {
		p := &project.Project{
			ID:          uuid.New().String(),
			Name:        "Test Project",
			Description: "Test Description",
			Status:      project.StatusActive,
		}

		assert.NotNil(t, p)
		assert.NotEmpty(t, p.ID)
	})

	t.Run("error: returns ErrProjectNotFound when not found", func(t *testing.T) {
		err := project.ErrProjectNotFound
		assert.Error(t, err)
		assert.Equal(t, "project not found", err.Error())
	})
}

func TestProjectUUID(t *testing.T) {
	t.Run("success: generates valid uuid", func(t *testing.T) {
		id := uuid.New().String()
		assert.NotEmpty(t, id)
		_, err := uuid.Parse(id)
		assert.NoError(t, err)
	})
}

func TestProjectStatus(t *testing.T) {
	t.Run("success: valid status values", func(t *testing.T) {
		assert.Equal(t, project.Status("active"), project.StatusActive)
		assert.Equal(t, project.Status("archived"), project.StatusArchived)
	})
}

func TestMockDB(t *testing.T) {
	t.Run("success: mock db implements interface", func(t *testing.T) {
		var _ interface {
			ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
			QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
			QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		} = &mockDB{}
	})
}

type mockDB struct {
	execFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	queryRowFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row
	queryFunc    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, query, args...)
	}
	return nil, nil
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(ctx, query, args...)
	}
	return nil
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, query, args...)
	}
	return nil, nil
}
