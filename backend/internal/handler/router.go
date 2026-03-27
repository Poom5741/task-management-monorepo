package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	projectHandler "github.com/poom5741/task-management-monorepo/backend/internal/handler/project"
	taskHandler "github.com/poom5741/task-management-monorepo/backend/internal/handler/task"
)

type Handlers struct {
	Project *projectHandler.ProjectHandler
	Task    *taskHandler.TaskHandler
}

func SetupRouter(projectH *projectHandler.ProjectHandler, taskH *taskHandler.TaskHandler) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		projects := v1.Group("/projects")
		{
			projects.GET("", projectH.ListProjects)
			projects.POST("", projectH.CreateProject)
			projects.GET("/:id", projectH.GetProject)
			projects.PATCH("/:id", projectH.UpdateProject)
			projects.DELETE("/:id", projectH.DeleteProject)

			projects.GET("/:id/tasks", taskH.ListTasks)
			projects.POST("/:id/tasks", taskH.CreateTask)
		}

		tasks := v1.Group("/tasks")
		{
			tasks.GET("/search", taskH.SearchTasks)
			tasks.GET("/:id", taskH.GetTask)
			tasks.PATCH("/:id", taskH.UpdateTask)
			tasks.DELETE("/:id", taskH.DeleteTask)

			tasks.POST("/:id/dependencies", taskH.AddDependency)
			tasks.DELETE("/:id/dependencies/:dependencyId", taskH.RemoveDependency)
		}

		labels := v1.Group("/labels")
		{
			labels.GET("", listLabels)
			labels.POST("", createLabel)
			labels.GET("/:id", getLabel)
			labels.PATCH("/:id", updateLabel)
			labels.DELETE("/:id", deleteLabel)
		}

		v1.GET("/dashboard", getDashboard)
	}

	return r
}

func listLabels(c *gin.Context) {
	c.JSON(200, gin.H{"message": "list labels"})
}

func createLabel(c *gin.Context) {
	c.JSON(201, gin.H{"message": "create label"})
}

func getLabel(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get label"})
}

func updateLabel(c *gin.Context) {
	c.JSON(200, gin.H{"message": "update label"})
}

func deleteLabel(c *gin.Context) {
	c.JSON(204, nil)
}

func getDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get dashboard"})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
