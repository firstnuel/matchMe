package user

import (
	"net/http"

	"match-me/api/middleware"
	"match-me/internal/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	var req requests.UpdateUser

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validationService.ValidateUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Update user
	user, err := h.UserUsecase.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Update failed",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
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

	var req requests.UpdatePasswordRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validationService.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Update password
	err = h.UserUsecase.UpdatePassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Password update failed",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Set the user ID parameter for the existing UpdateUser handler
	c.Params = append(c.Params, gin.Param{Key: "id", Value: user.ID.String()})
	h.UpdateUser(c)
}

func (h *UserHandler) UpdateCurrentUserPassword(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Set the user ID parameter for the existing UpdatePassword handler
	c.Params = append(c.Params, gin.Param{Key: "id", Value: user.ID.String()})
	h.UpdatePassword(c)
}
