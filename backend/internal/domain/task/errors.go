package task

import "errors"

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrInvalidTaskStatus   = errors.New("invalid task status")
	ErrInvalidTaskPriority = errors.New("invalid task priority")
	ErrInvalidDueDate      = errors.New("due date must be in the future")
	ErrTooManyLabels       = errors.New("task cannot have more than 10 labels")
)
