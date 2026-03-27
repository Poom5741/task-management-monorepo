package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
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
		SELECT 
			p.id, p.name, p.description, p.status, p.created_at, p.updated_at, p.deleted_at,
			COALESCE(t.task_count, 0) as task_count,
			CASE 
				WHEN COALESCE(t.task_count, 0) = 0 THEN 0
				ELSE ROUND((COALESCE(t.done_count, 0)::decimal / t.task_count::decimal) * 100, 2)
			END as completion_percentage
		FROM projects p
		LEFT JOIN (
			SELECT 
				project_id,
				COUNT(*) as task_count,
				COUNT(*) FILTER (WHERE status = 'done') as done_count
			FROM tasks
			WHERE deleted_at IS NULL
			GROUP BY project_id
		) t ON p.id = t.project_id
		WHERE p.id = $1 AND p.deleted_at IS NULL
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
		&p.TaskCount,
		&p.CompletionPercentage,
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
	if id == "" {
		return nil, project.ErrProjectNotFound
	}

	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	buildQuery := func() (string, []interface{}) {
		query := `UPDATE projects SET `
		args := []interface{}{}
		setClauses := []string{}
		argPos := 1

		if input.Name != nil {
			setClauses = append(setClauses, `name = $`+strconv.Itoa(argPos))
			args = append(args, *input.Name)
			argPos++
		}

		if input.Description != nil {
			setClauses = append(setClauses, `description = $`+strconv.Itoa(argPos))
			args = append(args, *input.Description)
			argPos++
		}

		if input.Status != nil {
			setClauses = append(setClauses, `status = $`+strconv.Itoa(argPos))
			args = append(args, string(*input.Status))
			argPos++
		}

		setClauses = append(setClauses, `updated_at = $`+strconv.Itoa(argPos))
		args = append(args, time.Now())
		argPos++

		query += strings.Join(setClauses, ", ")
		query += ` WHERE id = $` + strconv.Itoa(argPos) + ` AND deleted_at IS NULL`
		args = append(args, id)

		return query, args
	}

	query, args := buildQuery()
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, project.ErrProjectNameExists
		}
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, project.ErrProjectNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	taskQuery := `UPDATE tasks SET deleted_at = NOW() WHERE project_id = $1 AND deleted_at IS NULL`
	_, _ = r.db.ExecContext(ctx, taskQuery, id)

	query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return project.ErrProjectNotFound
	}

	return nil
}
