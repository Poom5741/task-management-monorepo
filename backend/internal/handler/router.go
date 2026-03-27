package handler

import (
	"github.com/gin-gonic/gin"
	projectHandler "github.com/poom5741/task-management-monorepo/backend/internal/handler/project"
)

type Handlers struct {
	Project *projectHandler.ProjectHandler
}

func SetupRouter(projectH *projectHandler.ProjectHandler) *gin.Engine {
	r := gin.Default()

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

			projects.GET("/:projectId/tasks", listTasks)
			projects.POST("/:projectId/tasks", createTask)
		}

		tasks := v1.Group("/tasks")
		{
			tasks.GET("/search", searchTasks)
			tasks.GET("/:id", getTask)
			tasks.PATCH("/:id", updateTask)
			tasks.DELETE("/:id", deleteTask)

			tasks.POST("/:taskId/dependencies", addDependency)
			tasks.DELETE("/:taskId/dependencies/:dependencyId", removeDependency)
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

func listTasks(c *gin.Context) {
	c.JSON(200, gin.H{"message": "list tasks"})
}

func createTask(c *gin.Context) {
	c.JSON(201, gin.H{"message": "create task"})
}

func getTask(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get task"})
}

func updateTask(c *gin.Context) {
	c.JSON(200, gin.H{"message": "update task"})
}

func deleteTask(c *gin.Context) {
	c.JSON(204, nil)
}

func searchTasks(c *gin.Context) {
	c.JSON(200, gin.H{"message": "search tasks"})
}

func addDependency(c *gin.Context) {
	c.JSON(201, gin.H{"message": "add dependency"})
}

func removeDependency(c *gin.Context) {
	c.JSON(204, nil)
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
