package project

import "errors"

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectNameExists    = errors.New("project name already exists")
	ErrInvalidProjectStatus = errors.New("invalid project status")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
