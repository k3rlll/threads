package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey      string
	accessTokenTTL time.Duration
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:      secretKey,
		accessTokenTTL: tokenTTL,
	}
}

// NewAccessToken generates a new JWT access token for the given user ID.
func (manager *JWTManager) NewAccessToken(userID uuid.UUID) (string, error) {
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(manager.accessTokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	})
	tokenString, err := jwtClaims.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyAccessToken verifies the access token and returns the user ID if the token is valid.
func (manager *JWTManager) VerifyAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return manager.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil || sub == "" {
		return "", jwt.ErrTokenMalformed
	}

	return sub, nil
}
