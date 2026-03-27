package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/task"
)

type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	if db == nil {
		return nil
	}
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}

	query := `
		INSERT INTO tasks (id, project_id, title, description, status, priority, due_date, labels, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.ProjectID,
		t.Title,
		t.Description,
		string(t.Status),
		string(t.Priority),
		t.DueDate,
		pq.Array(t.Labels),
		t.CreatedAt,
		t.UpdatedAt,
	)

	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*task.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, due_date, labels, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	var t task.Task
	var status, priority string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.ProjectID, &t.Title, &t.Description,
		&status, &priority, &t.DueDate, pq.Array(&t.Labels),
		&t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, task.ErrTaskNotFound
		}
		return nil, err
	}

	t.Status = task.Status(status)
	t.Priority = task.Priority(priority)
	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context, projectID string, filter *task.TaskListFilter) ([]*task.Task, int, error) {
	return nil, 0, nil
}

func (r *TaskRepository) Update(ctx context.Context, id string, input *task.UpdateTaskInput) (*task.Task, error) {
	return nil, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TaskRepository) Search(ctx context.Context, query string, page, pageSize int) ([]*task.Task, int, error) {
	return nil, 0, nil
}

func (r *TaskRepository) ProjectExists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
