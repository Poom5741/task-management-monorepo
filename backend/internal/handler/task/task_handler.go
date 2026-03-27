package task

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/task"
)

type TaskHandler struct {
	usecase task.Usecase
}

func NewTaskHandler(usecase task.Usecase) *TaskHandler {
	if usecase == nil {
		return nil
	}
	return &TaskHandler{usecase: usecase}
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Description string     `json:"description" binding:"max=2000"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
	Labels      []string   `json:"labels" binding:"max=10"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "projectId is required",
		})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	input := &task.CreateTaskInput{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    task.Priority(req.Priority),
		DueDate:     req.DueDate,
		Labels:      req.Labels,
	}

	t, err := h.usecase.CreateTask(c.Request.Context(), input)
	if err != nil {
		var validationErr *task.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: validationErr.Message,
			})
			return
		}

		if errors.Is(err, task.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "project not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get task"})
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list tasks"})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update task"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) SearchTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "search tasks"})
}

func (h *TaskHandler) AddDependency(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "add dependency"})
}

func (h *TaskHandler) RemoveDependency(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
