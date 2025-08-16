package user

import (
	"net/http"

	"match-me/internal/requests"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Register(c *gin.Context) {
	var req requests.RegisterUser

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validationService.ValidateRegisterUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Register user
	user, token, err := h.userUsecase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}
