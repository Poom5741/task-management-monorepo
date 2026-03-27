package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
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
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	countQuery := `SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL`
	countArgs := []interface{}{}
	argPos := 1

	if filter.Search != "" {
		countQuery += ` AND name ILIKE $` + strconv.Itoa(argPos)
		countArgs = append(countArgs, "%"+filter.Search+"%")
		argPos++
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, description, status, created_at, updated_at FROM projects WHERE deleted_at IS NULL`
	queryArgs := []interface{}{}

	if filter.Search != "" {
		query += ` AND name ILIKE $1`
		queryArgs = append(queryArgs, "%"+filter.Search+"%")
		argPos = 2
	} else {
		argPos = 1
	}

	query += ` ORDER BY created_at DESC`
	query += ` LIMIT $` + strconv.Itoa(argPos) + ` OFFSET $` + strconv.Itoa(argPos+1)
	queryArgs = append(queryArgs, filter.PageSize, (filter.Page-1)*filter.PageSize)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		var p project.Project
		var status string
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		p.Status = project.Status(status)
		p.TaskCount = 0
		projects = append(projects, &p)
	}

	if projects == nil {
		projects = []*project.Project{}
	}

	return projects, total, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
	return nil, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return nil
}
