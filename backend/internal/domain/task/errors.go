package task

import "errors"

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrInvalidTaskStatus   = errors.New("invalid task status")
	ErrInvalidTaskPriority = errors.New("invalid task priority")
	ErrInvalidDueDate      = errors.New("due date must be in the future")
	ErrTooManyLabels       = errors.New("task cannot have more than 10 labels")
	ErrProjectNotFound     = errors.New("project not found")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
