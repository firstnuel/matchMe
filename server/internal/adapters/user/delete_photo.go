package user

import (
	"match-me/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) DeleteUserPhoto(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Get photo ID from URL parameter
	photoIDStr := c.Param("photoId")
	if photoIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing photo ID",
			"details": "Photo ID is required in the URL path",
		})
		return
	}

	// Parse photo ID
	photoID, err := uuid.Parse(photoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid photo ID",
			"details": "Photo ID must be a valid UUID",
		})
		return
	}

	// Delete photo via usecase
	err = h.UserUsecase.DeleteUserPhoto(c.Request.Context(), user.ID, photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete photo",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":  "Photo deleted successfully",
		"photo_id": photoID.String(),
	})
}