package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	PurposePasswordReset string = "password_reset"
	PurposeLogin         string = "user_login"
)

func GenerateJWTToken(userID uuid.UUID, jwtSecret string, purpose string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	issuedAt := time.Now()

	claims := jwt.MapClaims{
		"userId":  userID.String(),
		"purpose": purpose,
		"exp":     expirationTime.Unix(),
		"iat":     issuedAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	return signedToken, err
}
