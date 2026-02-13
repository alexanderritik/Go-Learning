package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator interface {
	GenerateToken(userID string, duration time.Duration) (string, error)
	ValidateToken(tokenStr string) (string, error)
}

type TokenManager struct {
	secretKey []byte
	issuer    string
}

func NewTokenManager(secret string, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

func (m *TokenManager) GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                          // Subject (User ID)
		"iss": m.issuer,                        // Issuer
		"exp": time.Now().Add(duration).Unix(), // Expiration
		"iat": time.Now().Unix(),               // Issued At
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	return token.SignedString(m.secretKey)
}

func (m *TokenManager) ValidateToken(tokenStr string) (string, error) {

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}

	return "", errors.New("invalid token")
}
