package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey string
	tokenTTL  time.Duration
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

func (manager *JWTManager) NewAcccessToken(userID string) (string, error) {
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(manager.tokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	})
	tokenString, err := jwtClaims.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (manager *JWTManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

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
