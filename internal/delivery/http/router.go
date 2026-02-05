package http

import (
	"log/slog"
	handler "main/internal/delivery/grpc/auth"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func MapRoutes(e *echo.Echo, authHandler *handler.AuthHandler, logger slog.Logger) {
	// Middlewares
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Auth routes
	v1 := e.Group("/api/v1")

	
	authGroup := v1.Group("/auth")
	authGroup.POST("/register", authHandler.RegisterUser)
	authGroup.POST("/login", authHandler.LoginUser)
	authGroup.POST("/refresh", authHandler.RefreshToken)
}
