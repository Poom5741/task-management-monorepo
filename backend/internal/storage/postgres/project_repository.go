package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
)

type ProjectRepository struct {
	db *DB
}

func NewProjectRepository(db *DB) *ProjectRepository {
	if db == nil {
		return nil
	}
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *project.Project) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}

	query := `
		INSERT INTO projects (id, name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.Name,
		p.Description,
		string(p.Status),
		p.CreatedAt,
		p.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return project.ErrProjectNameExists
			}
		}
		return err
	}

	return nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*project.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL
	`

	var p project.Project
	var status string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, project.ErrProjectNotFound
		}
		return nil, err
	}

	p.Status = project.Status(status)
	return &p, nil
}

func (r *ProjectRepository) GetByName(ctx context.Context, name string) (*project.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at
		FROM projects
		WHERE name = $1 AND deleted_at IS NULL
	`

	var p project.Project
	var status string
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	p.Status = project.Status(status)
	return &p, nil
}

func (r *ProjectRepository) List(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
	return nil, 0, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
	return nil, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return nil
}
