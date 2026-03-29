package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userID uint, role string) (string, string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, errAccess := accessTokenRaw.SignedString(secretKey)

	if errAccess != nil {
		return "", "", errAccess
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, errRefresh := refreshTokenRaw.SignedString(secretKey)

	if errRefresh != nil {
		return "", "", errRefresh
	}

	return accessTokenString, refreshTokenString, nil
}
