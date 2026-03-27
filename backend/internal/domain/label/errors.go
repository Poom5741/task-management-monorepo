package label

import "errors"

var (
	ErrLabelNotFound   = errors.New("label not found")
	ErrLabelNameExists = errors.New("label name already exists")
	ErrInvalidColor    = errors.New("invalid color format")
)
