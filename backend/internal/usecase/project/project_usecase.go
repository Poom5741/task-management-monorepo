package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
)

type ProjectUsecase struct {
	repo project.Repository
}

func NewProjectUsecase(repo project.Repository) *ProjectUsecase {
	if repo == nil {
		return nil
	}
	return &ProjectUsecase{repo: repo}
}

func (uc *ProjectUsecase) CreateProject(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}

	if input.Name == "" {
		return nil, &project.ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	if len(input.Name) > 100 {
		return nil, &project.ValidationError{
			Field:   "name",
			Message: "name must be at most 100 characters",
		}
	}

	if len(input.Description) > 500 {
		return nil, &project.ValidationError{
			Field:   "description",
			Message: "description must be at most 500 characters",
		}
	}

	existing, err := uc.repo.GetByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, project.ErrProjectNameExists
	}

	now := time.Now()
	p := &project.Project{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Status:      project.StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *ProjectUsecase) GetProject(ctx context.Context, id string) (*project.Project, error) {
	if id == "" {
		return nil, &project.ValidationError{
			Field:   "id",
			Message: "id is required",
		}
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *ProjectUsecase) ListProjects(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *ProjectUsecase) UpdateProject(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
	if id == "" {
		return nil, &project.ValidationError{
			Field:   "id",
			Message: "id is required",
		}
	}

	if input == nil {
		return nil, errors.New("input is required")
	}

	if input.Name != nil {
		if len(*input.Name) > 100 {
			return nil, &project.ValidationError{
				Field:   "name",
				Message: "name must be at most 100 characters",
			}
		}
	}

	if input.Description != nil {
		if len(*input.Description) > 500 {
			return nil, &project.ValidationError{
				Field:   "description",
				Message: "description must be at most 500 characters",
			}
		}
	}

	if input.Status != nil {
		if *input.Status != project.StatusActive && *input.Status != project.StatusArchived {
			return nil, &project.ValidationError{
				Field:   "status",
				Message: "status must be either active or archived",
			}
		}
	}

	return uc.repo.Update(ctx, id, input)
}

func (uc *ProjectUsecase) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return &project.ValidationError{
			Field:   "id",
			Message: "id is required",
		}
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *ProjectUsecase) GetProjectStatistics(ctx context.Context, id string) (*project.ProjectStatistics, error) {
	return nil, nil
}
