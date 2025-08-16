package user

import (
	"match-me/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) DeleteUser(c *gin.Context) {
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

	// Delete user from usecase
	err = h.userUsecase.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Set the user ID parameter for the existing DeleteUser handler
	c.Params = append(c.Params, gin.Param{Key: "id", Value: user.ID.String()})
	h.DeleteUser(c)
}
