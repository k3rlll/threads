package grp

import (
	"context"
	"log/slog"
	authv1 "main/pkg/proto/gen/auth/v1"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	logger      slog.Logger
	AuthUsecase AuthUsecase
	e           *echo.Echo
}

type AuthUsecase interface {
	RegisterUser(ctx context.Context, username, email, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

func NewAuthHandler(logger slog.Logger, authUsecase AuthUsecase, e *echo.Echo) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		AuthUsecase: authUsecase,
		e:           e,
	}
}




func (h *AuthHandler) RegisterUser(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	userID, err := h.AuthUsecase.RegisterUser(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		return nil, err
	}
	return &authv1.RegisterResponse{
		UserId: userID}, nil

}

func (h *AuthHandler) LoginUser(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" {
		h.logger.Error("Login or password is empty")
		return nil, status.Error(codes.InvalidArgument, "login or password is empty")
	}
	userID, err := h.AuthUsecase.LoginUser(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		h.logger.Error("Failed to login user", "error", err)
		return nil, err
	}
	return &authv1.LoginResponse{
		UserId: userID}, nil

}
