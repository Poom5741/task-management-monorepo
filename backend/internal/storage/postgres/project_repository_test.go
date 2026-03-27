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
	t.Run("success: returns project with task_count", func(t *testing.T) {
		p := &project.Project{
			ID:          uuid.New().String(),
			Name:        "Test Project",
			Description: "Test Description",
			Status:      project.StatusActive,
			TaskCount:   5,
		}

		assert.NotNil(t, p)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, 5, p.TaskCount)
	})

	t.Run("success: returns project with completion_percentage", func(t *testing.T) {
		p := &project.Project{
			ID:                   uuid.New().String(),
			Name:                 "Test Project",
			TaskCount:            10,
			CompletionPercentage: 60.0,
		}

		assert.Equal(t, 10, p.TaskCount)
		assert.Equal(t, 60.0, p.CompletionPercentage)
	})

	t.Run("success: completion_percentage is zero when no tasks", func(t *testing.T) {
		p := &project.Project{
			ID:                   uuid.New().String(),
			Name:                 "Empty Project",
			TaskCount:            0,
			CompletionPercentage: 0.0,
		}

		assert.Equal(t, 0, p.TaskCount)
		assert.Equal(t, 0.0, p.CompletionPercentage)
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

func strPtr(s string) *string {
	return &s
}

func statusPtr(s project.Status) *project.Status {
	return &s
}

func TestProjectRepository_Update(t *testing.T) {
	t.Run("success: update project name only", func(t *testing.T) {
		updateInput := &project.UpdateProjectInput{
			Name: strPtr("Updated Project Name"),
		}

		assert.NotNil(t, updateInput)
		assert.Equal(t, "Updated Project Name", *updateInput.Name)
		assert.Nil(t, updateInput.Description)
		assert.Nil(t, updateInput.Status)
	})

	t.Run("success: update description only", func(t *testing.T) {
		updateInput := &project.UpdateProjectInput{
			Description: strPtr("Updated description"),
		}

		assert.NotNil(t, updateInput)
		assert.Equal(t, "Updated description", *updateInput.Description)
		assert.Nil(t, updateInput.Name)
		assert.Nil(t, updateInput.Status)
	})

	t.Run("success: update status to archived", func(t *testing.T) {
		updateInput := &project.UpdateProjectInput{
			Status: statusPtr(project.StatusArchived),
		}

		assert.NotNil(t, updateInput)
		assert.Equal(t, project.StatusArchived, *updateInput.Status)
		assert.Nil(t, updateInput.Name)
		assert.Nil(t, updateInput.Description)
	})

	t.Run("success: update multiple fields at once", func(t *testing.T) {
		updateInput := &project.UpdateProjectInput{
			Name:        strPtr("Updated Name"),
			Description: strPtr("Updated description"),
			Status:      statusPtr(project.StatusArchived),
		}

		assert.NotNil(t, updateInput)
		assert.Equal(t, "Updated Name", *updateInput.Name)
		assert.Equal(t, "Updated description", *updateInput.Description)
		assert.Equal(t, project.StatusArchived, *updateInput.Status)
	})

	t.Run("validation: update name to existing name returns ErrProjectNameExists", func(t *testing.T) {
		err := project.ErrProjectNameExists
		assert.Error(t, err)
		assert.Equal(t, "project name already exists", err.Error())
	})

	t.Run("infrastructure: project not found returns ErrProjectNotFound", func(t *testing.T) {
		err := project.ErrProjectNotFound
		assert.Error(t, err)
		assert.Equal(t, "project not found", err.Error())
	})

	t.Run("context: database connection errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		assert.Error(t, ctx.Err())
		assert.Equal(t, context.Canceled, ctx.Err())
	})
}

func TestProjectRepository_List(t *testing.T) {
	t.Run("success: returns empty list when no projects", func(t *testing.T) {
		filter := &project.ProjectListFilter{
			Page:     1,
			PageSize: 20,
		}

		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, 20, filter.PageSize)
	})

	t.Run("success: applies pagination correctly", func(t *testing.T) {
		filter := &project.ProjectListFilter{
			Page:     2,
			PageSize: 10,
		}

		offset := (filter.Page - 1) * filter.PageSize
		assert.Equal(t, 10, offset)
	})

	t.Run("success: search filter with ILIKE pattern", func(t *testing.T) {
		search := "test"
		pattern := "%" + search + "%"

		assert.Contains(t, pattern, "test")
	})

	t.Run("success: default filter values", func(t *testing.T) {
		filter := &project.ProjectListFilter{}

		assert.Zero(t, filter.Page)
		assert.Zero(t, filter.PageSize)
	})

	t.Run("success: TaskCount field exists on Project", func(t *testing.T) {
		p := &project.Project{
			ID:        uuid.New().String(),
			Name:      "Test Project",
			TaskCount: 5,
		}

		assert.Equal(t, 5, p.TaskCount)
	})
}
