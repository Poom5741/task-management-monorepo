package label

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, l *Label) error
	GetByID(ctx context.Context, id string) (*Label, error)
	List(ctx context.Context, filter *LabelListFilter) ([]*Label, int, error)
	Update(ctx context.Context, id string, input *UpdateLabelInput) (*Label, error)
	Delete(ctx context.Context, id string) error
}

type Usecase interface {
	CreateLabel(ctx context.Context, input *CreateLabelInput) (*Label, error)
	GetLabel(ctx context.Context, id string) (*Label, error)
	ListLabels(ctx context.Context, filter *LabelListFilter) ([]*Label, int, error)
	UpdateLabel(ctx context.Context, id string, input *UpdateLabelInput) (*Label, error)
	DeleteLabel(ctx context.Context, id string) error
}
