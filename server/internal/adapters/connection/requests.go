package connection

import (
	"net/http"

	"match-me/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendConnectionRequestBody represents the request body for sending a connection request
type SendConnectionRequestBody struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Message    string    `json:"message,omitempty"`
}

// SendConnectionRequest handles POST /connection-requests
func (h *ConnectionHandler) SendConnectionRequest(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req SendConnectionRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Send connection request
	request, err := h.ConnectionRequestUsecase.SendRequest(c.Request.Context(), user.ID, req.ReceiverID, req.Message)
	if err != nil {
		if err.Error() == "cannot send connection request to yourself" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": "Cannot send connection request to yourself",
			})
			return
		}
		if err.Error() == "connection already exists between users" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Connection exists",
				"details": "You are already connected to this user",
			})
			return
		}
		if err.Error() == "connection request already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Request exists",
				"details": "A connection request already exists between these users",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send connection request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Connection request sent successfully",
		"request": request,
	})
}

// GetPendingRequests handles GET /connection-requests
func (h *ConnectionHandler) GetPendingRequests(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get pending requests
	requests, err := h.ConnectionRequestUsecase.GetPendingRequests(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get pending requests",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"count":    len(requests),
	})
}

// AcceptRequest handles PUT /connection-requests/:requestId/accept
func (h *ConnectionHandler) AcceptRequest(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request ID
	requestIDStr := c.Param("requestId")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request ID",
			"details": "Request ID must be a valid UUID",
		})
		return
	}

	// Accept request
	connection, err := h.ConnectionRequestUsecase.AcceptRequest(c.Request.Context(), user.ID, requestID)
	if err != nil {
		if err.Error() == "connection request not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Request not found",
				"details": "The specified connection request does not exist",
			})
			return
		}
		if err.Error() == "unauthorized: user is not the receiver of this request" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to accept this request",
			})
			return
		}
		if err.Error() == "request is no longer pending" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Request not pending",
				"details": "This request has already been responded to",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to accept request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Connection request accepted successfully",
		"connection": connection,
	})
}

// DeclineRequest handles PUT /connection-requests/:requestId/decline
func (h *ConnectionHandler) DeclineRequest(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request ID
	requestIDStr := c.Param("requestId")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request ID",
			"details": "Request ID must be a valid UUID",
		})
		return
	}

	// Decline request
	err = h.ConnectionRequestUsecase.DeclineRequest(c.Request.Context(), user.ID, requestID)
	if err != nil {
		if err.Error() == "connection request not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Request not found",
				"details": "The specified connection request does not exist",
			})
			return
		}
		if err.Error() == "unauthorized: user is not the receiver of this request" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to decline this request",
			})
			return
		}
		if err.Error() == "request is no longer pending" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Request not pending",
				"details": "This request has already been responded to",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to decline request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Connection request declined successfully",
	})
}