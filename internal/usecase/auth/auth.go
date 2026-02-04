package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, userID uuid.UUID, email, username, passwordHash string) (string, error)
}

type AuthUsecase struct {
	authRepo AuthRepo
}

func NewAuthUsecase(authRepo AuthRepo) *AuthUsecase {
	return &AuthUsecase{authRepo: authRepo}
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, email, username, password string) (string, error) {

	if !validateUsername(username) {
		return "", errors.New("username must be between 3 and 30 characters")
	}

	if !validateEmail(email) {
		return "", errors.New("invalid email format")
	}

	if !validatePassword(password) {
		return "",errors.New("password does not meet the requirements")
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	userID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return uc.authRepo.CreateUser(ctx, userID, email, username, passwordHash)

}


func (uc *AuthUsecase) LoginUser(ctx context.Context, login, password string) (string, error) {
	
}

// hashPassword hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
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
