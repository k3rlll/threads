package auth

import (
	"context"
	"errors"
	"time"

	"main/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthRepo defines the interface for authentication-related database operations.
type AuthRepo interface {
	// CreateUser creates a new user in the database with the provided details and returns the user ID.
	CreateUser(ctx context.Context, userID uuid.UUID, email, username, passwordHash string) (uuid.UUID, error)

	// GetUserByLogin retrieves the user ID and password hash based on the provided login (username or email).
	GetUserByLogin(ctx context.Context, login string) (userID uuid.UUID, passwordHash string, err error)

	// StoreSession saves the session associated with a user in the database, allowing for session management and token revocation.
	StoreSession(ctx context.Context, userID uuid.UUID, session entity.Session) error

	// DeleteSession removes a specific session for a user, effectively logging them out from that ONE SPECIFIC SESSION.
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error

	// DeleteAllSessions removes all sessions associated with a user, effectively logging them out from !ALL! devices.
	DeleteAllSessions(ctx context.Context, userID uuid.UUID) error
}

// JWTManager defines the interface for JWT token management.
type JWTManager interface {
	NewAccessToken(userID uuid.UUID) (string, error)
	VerifyAccessToken(token string) (string, error)
}

type AuthUsecase struct {
	authRepo   AuthRepo
	JWTManager JWTManager
}

func NewAuthUsecase(authRepo AuthRepo, JWTManager JWTManager) *AuthUsecase {
	return &AuthUsecase{
		authRepo:   authRepo,
		JWTManager: JWTManager}
}

// RegisterUser validates the input, hashes the password, and creates a new user in the database.
// It returns the user ID as a string or an error if the registration fails.
func (uc *AuthUsecase) RegisterUser(ctx context.Context, username, email, password string) (userID uuid.UUID, err error) {

	if !validateUsername(username) {
		return uuid.Nil, errors.New("username must be between 3 and 30 characters")
	}

	if !validateEmail(email) {
		return uuid.Nil, errors.New("invalid email format")
	}

	if !validatePassword(password) {
		return uuid.Nil, errors.New("password does not meet the requirements")
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return uuid.Nil, err
	}
	userID, err = uuid.NewUUID()
	if err != nil {
		return uuid.Nil, err
	}

	return uc.authRepo.CreateUser(ctx, userID, email, username, passwordHash)

}

// LoginUser authenticates the user by verifying the provided credentials.
// If successful, it generates an access token and a refresh token, stores the session in the database, and returns the access token.
// If authentication fails, it returns an error.
func (uc *AuthUsecase) LoginUser(ctx context.Context, login, password, userAgent string, ip string) (accessToken string, err error) {
	userID, passwordHash, err := uc.authRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if !verifyPassword(password, passwordHash) {
		return "", errors.New("invalid credentials")
	}

	accessToken, err = uc.JWTManager.NewAccessToken(userID)
	if err != nil {
		return "", err
	}

	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	session := entity.Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now().Unix(),
		ExpiresAt:    time.Now().Add(15 * 24 * time.Hour).Unix(),
		UserAgent:    userAgent,
		ClientIP:     ip,
	}

	err = uc.authRepo.StoreSession(ctx, userID, session)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// LogoutSession logs out the user from a specific session by deleting that session from the database.
func (uc *AuthUsecase) LogoutSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	err := uc.authRepo.DeleteSession(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	return nil
}

// LogoutAllSessions logs out the user from all sessions by deleting all sessions associated with the user from the database.
func (uc *AuthUsecase) LogoutAllSessions(ctx context.Context, userID uuid.UUID) error {
	err := uc.authRepo.DeleteAllSessions(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// hashPassword hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

func verifyPassword(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

// ValidatePassword checks if the password meets certain criteria
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !containsUppercase(password) {
		return false
	}
	return true
}

// Simple function to check if a string contains at least one uppercase letter
func containsUppercase(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

// Simple function to check if a string contains '@' symbol
func containsAtSymbol(s string) bool {
	for _, char := range s {
		if char == '@' {
			return true
		}
	}
	return false
}

func validateEmail(email string) bool {
	// Simple email validation
	if len(email) < 5 || len(email) > 50 {
		return false
	}
	if !containsAtSymbol(email) {
		return false
	}
	return true
}

func validateUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	return true
}
