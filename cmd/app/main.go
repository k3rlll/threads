package main

import (
	"log/slog"
	"main/internal/config"
	gRPCHandler "main/internal/delivery/grpc/auth"
	psql "main/internal/storage/postgres"
	authRepo "main/internal/storage/postgres/auth"
	authUs "main/internal/usecase/auth"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	config := config.LoadConfig()
	logger := setupLogger(config.Env)
	logger.Info("Application started", "env", config.Env)

	// Initialize Postgres connection
	DSN := config.PostgresConfig.DSN()
	pool, err := psql.NewPostgresConnection(DSN)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		return
	}
	defer pool.Close()
	logger.Info("Connected to the database successfully")
	e := echo.New()
	authRepo := authRepo.NewAuthRepo(pool)
	authUsecase := authUs.NewAuthUsecase(authRepo)
	authHadler := gRPCHandler.NewAuthHandler(*logger, authUsecase)
	
	
	_ = e.Start(":8080")
	logger.Info("gRPC server started on port 8080")
	
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "production":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "development", "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
