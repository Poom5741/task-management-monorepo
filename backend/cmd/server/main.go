package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/poom5741/task-management-monorepo/backend/internal/handler"
	projectHandler "github.com/poom5741/task-management-monorepo/backend/internal/handler/project"
	taskHandler "github.com/poom5741/task-management-monorepo/backend/internal/handler/task"
	"github.com/poom5741/task-management-monorepo/backend/internal/storage/postgres"
	projectUsecase "github.com/poom5741/task-management-monorepo/backend/internal/usecase/project"
	taskUsecase "github.com/poom5741/task-management-monorepo/backend/internal/usecase/task"
	"github.com/poom5741/task-management-monorepo/backend/pkg/config"
	"github.com/poom5741/task-management-monorepo/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.LogLevel)

	logger.Info("Starting Task Management API server...")
	logger.Info("Server will run on port: " + cfg.ServerPort)

	db, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	projectRepo := postgres.NewProjectRepository(db)
	projectUC := projectUsecase.NewProjectUsecase(projectRepo)
	projectH := projectHandler.NewProjectHandler(projectUC)

	taskRepo := postgres.NewTaskRepository(db)
	taskUC := taskUsecase.NewTaskUsecase(taskRepo)
	taskH := taskHandler.NewTaskHandler(taskUC)

	router := handler.SetupRouter(projectH, taskH)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		logger.Info("Server listening on port " + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exited")
}
