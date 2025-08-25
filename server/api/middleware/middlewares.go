package middleware

import (
	"context"
	"log"
	"match-me/internal/models"
	"match-me/internal/pkg/jwt"
	"match-me/internal/usecases/user"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Ping returns a middleware that responds with "pong" to GET requests on "/ping".
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && c.Request.URL.Path == "/ping" {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:        []string{"http://localhost:5173", "http://192.168.100.44:5173", "http://127.0.0.1:5500"}, // frontend origin
		AllowMethods:        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:        []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:       []string{"Content-Length"},
		AllowCredentials:    true,
		MaxAge:              12 * time.Hour,
		AllowPrivateNetwork: true,
	})
}

// UserContextKey is the key used to store user in context
const UserContextKey = "user"

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}

// GetUserFromGinContext extracts the user from Gin context
func GetUserFromGinContext(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	userModel, ok := user.(*models.User)
	return userModel, ok
}

// VerifyUser is middleware that extracts a user from a valid access token
func VerifyUser(userUC user.UserUsecase, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		tokenStr := ""

		// Extract token from Bearer header
		if authHeader != "" {
			if token, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
				tokenStr = token
			}
		} else if token := c.Request.URL.Query().Get("token"); token != "" {
			tokenStr = token
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Missing access token. Please login.",
			})
			c.Abort()
			return
		}

		// Verify access token
		userID, err := jwt.VerifyJwtToken(c.Request.Context(), tokenStr, jwt.PurposeLogin, jwtSecret)
		if err != nil {
			log.Printf("[middleware]: VerifyUser Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "Invalid or expired access token. Please login.",
			})
			c.Abort()
			return
		}

		// Load user from DB
		user, err := userUC.GetUserByID(c.Request.Context(), userID, models.AccessLevelFull)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"details": "User not found",
			})
			c.Abort()
			return
		}

		// Store user in context
		c.Set(UserContextKey, user)

		// Proceed to next middleware/handler
		c.Next()
	}
}
