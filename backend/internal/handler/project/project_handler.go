package project

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/poom5741/task-management-monorepo/backend/internal/domain/project"
)

type ProjectHandler struct {
	usecase project.Usecase
}

func NewProjectHandler(usecase project.Usecase) *ProjectHandler {
	if usecase == nil {
		return nil
	}
	return &ProjectHandler{usecase: usecase}
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	input := &project.CreateProjectInput{
		Name:        req.Name,
		Description: req.Description,
	}

	p, err := h.usecase.CreateProject(c.Request.Context(), input)
	if err != nil {
		var validationErr *project.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation_error",
				Message: validationErr.Message,
			})
			return
		}

		if errors.Is(err, project.ErrProjectNameExists) {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "conflict",
				Message: "project name already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "id is required",
		})
		return
	}

	p, err := h.usecase.GetProject(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
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

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	filter := &project.ProjectListFilter{
		Page:     1,
		PageSize: 20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			filter.PageSize = ps
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	projects, total, err := h.usecase.ListProjects(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  projects,
		"total": total,
	})
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "id is required",
		})
		return
	}

	var req project.UpdateProjectInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	p, err := h.usecase.UpdateProject(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
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

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "id is required",
		})
		return
	}

	if err := h.usecase.DeleteProject(c.Request.Context(), id); err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
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

	c.Status(http.StatusNoContent)
}
