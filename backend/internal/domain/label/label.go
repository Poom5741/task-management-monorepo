package label

type Label struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	TaskCount int    `json:"task_count"`
}

type CreateLabelInput struct {
	Name  string `json:"name" validate:"required,max=50"`
	Color string `json:"color" validate:"required,hexcolor"`
}

type UpdateLabelInput struct {
	Name  *string `json:"name" validate:"omitempty,max=50"`
	Color *string `json:"color" validate:"omitempty,hexcolor"`
}

type LabelListFilter struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Search   string `form:"search" validate:"max=100"`
}
