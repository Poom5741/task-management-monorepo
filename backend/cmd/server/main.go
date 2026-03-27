package main

import (
	"log"

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
}
