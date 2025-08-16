package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func VerifyJwtToken(ctx context.Context, tokenStr string, purpose, jwtSecret string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	purposeStr, ok := claims["purpose"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid or missing purpose in token claims")
	}

	if purposeStr != purpose {
		return uuid.Nil, fmt.Errorf("token purpose mismatch")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return uuid.Nil, fmt.Errorf("missing token expiration")
	}

	expirationTime := time.Unix(int64(expFloat), 0)
	if expirationTime.Before(time.Now()) {
		return uuid.Nil, fmt.Errorf("token has expired")
	}

	userIdStr, ok := claims["userId"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	return userID, nil
}
