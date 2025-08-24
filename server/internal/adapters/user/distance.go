package user

import (
	"match-me/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) GetDistanceBetweenUsers(c *gin.Context) {
	// Get current user from context
	currentUser, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Get target user ID from URL parameter
	targetUserIDStr := c.Param("id")
	if targetUserIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing user ID",
			"details": "Target user ID is required in the URL path",
		})
		return
	}

	// Parse target user ID
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": "User ID must be a valid UUID",
		})
		return
	}

	// Get distance between users via usecase
	distance, err := h.UserUsecase.GetDistanceBetweenUsers(c.Request.Context(), currentUser.ID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get distance",
			"details": err.Error(),
		})
		return
	}

	// Return distance response
	c.JSON(http.StatusOK, gin.H{
		"distance":     distance,
		"unit":         "km",
		"current_user": currentUser.ID.String(),
		"target_user":  targetUserID.String(),
	})
}
