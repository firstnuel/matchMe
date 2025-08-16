package user

import (
	"match-me/internal/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(c *gin.Context) {
	var req requests.LoginRequest

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

	// Login user
	user, token, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}
