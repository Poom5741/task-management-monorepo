package project

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id string) (*Project, error)
	List(ctx context.Context, filter *ProjectListFilter) ([]*Project, int, error)
	Update(ctx context.Context, id string, input *UpdateProjectInput) (*Project, error)
	Delete(ctx context.Context, id string) error
}

type Usecase interface {
	CreateProject(ctx context.Context, input *CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	ListProjects(ctx context.Context, filter *ProjectListFilter) ([]*Project, int, error)
	UpdateProject(ctx context.Context, id string, input *UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, id string) error
	GetProjectStatistics(ctx context.Context, id string) (*ProjectStatistics, error)
}
