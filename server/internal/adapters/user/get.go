package user

import (
	"match-me/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": "User ID must be a valid UUID",
		})
		return
	}

	// Get user from usecase
	user, err := h.userUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Current user retrieved successfully",
		"user":    user,
	})
}
