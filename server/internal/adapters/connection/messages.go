package connection

import (
	"net/http"
	"strconv"

	"match-me/api/middleware"
	"match-me/internal/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SendTextMessage handles POST /messages/text
func (h *ConnectionHandler) SendTextMessage(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req requests.SendTextMessageBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Send text message
	message, err := h.MessageUsecase.SendTextMessage(c.Request.Context(), user.ID, req.ConnectionID, req.Content)
	if err != nil {
		if err.Error() == "connection not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Connection not found",
				"details": "The specified connection does not exist",
			})
			return
		}
		if err.Error() == "connection is not active" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Connection inactive",
				"details": "Cannot send messages to an inactive connection",
			})
			return
		}
		if err.Error() == "unauthorized: user is not part of this connection" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to send messages in this connection",
			})
			return
		}
		if err.Error() == "message content cannot be empty" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid message",
				"details": "Message content cannot be empty",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send message",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Text message sent successfully",
		"data":    message,
	})
}

// SendMediaMessage handles POST /messages/media
func (h *ConnectionHandler) SendMediaMessage(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse connection ID from form data
	connectionIDStr := c.PostForm("connection_id")
	if connectionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing connection_id",
			"details": "Connection ID is required",
		})
		return
	}

	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid connection ID",
			"details": "Connection ID must be a valid UUID",
		})
		return
	}

	// Get media file from form
	file, err := c.FormFile("media")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing media file",
			"details": "Media file is required",
		})
		return
	}

	// Send media message
	message, err := h.MessageUsecase.SendMediaMessage(c.Request.Context(), user.ID, connectionID, file)
	if err != nil {
		if err.Error() == "connection not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Connection not found",
				"details": "The specified connection does not exist",
			})
			return
		}
		if err.Error() == "connection is not active" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Connection inactive",
				"details": "Cannot send messages to an inactive connection",
			})
			return
		}
		if err.Error() == "unauthorized: user is not part of this connection" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to send messages in this connection",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send media message",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Media message sent successfully",
		"data":    message,
	})
}

// GetConnectionMessages handles GET /messages/connection/:connectionId
func (h *ConnectionHandler) GetConnectionMessages(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse connection ID
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid connection ID",
			"details": "Connection ID must be a valid UUID",
		})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit",
			"details": "Limit must be a number between 0 and 100",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid offset",
			"details": "Offset must be a non-negative number",
		})
		return
	}

	// Get connection messages
	messages, err := h.MessageUsecase.GetConnectionMessages(c.Request.Context(), user.ID, connectionID, limit, offset)
	if err != nil {
		if err.Error() == "connection not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Connection not found",
				"details": "The specified connection does not exist",
			})
			return
		}
		if err.Error() == "connection is not active" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Connection inactive",
				"details": "Cannot access messages from an inactive connection",
			})
			return
		}
		if err.Error() == "unauthorized: user is not part of this connection" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to access messages from this connection",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get messages",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
		"limit":    limit,
		"offset":   offset,
	})
}

// MarkMessagesAsRead handles PUT /messages/connection/:connectionId/read
func (h *ConnectionHandler) MarkMessagesAsRead(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse connection ID
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid connection ID",
			"details": "Connection ID must be a valid UUID",
		})
		return
	}

	// Mark messages as read
	err = h.MessageUsecase.MarkMessagesAsRead(c.Request.Context(), user.ID, connectionID)
	if err != nil {
		if err.Error() == "connection not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Connection not found",
				"details": "The specified connection does not exist",
			})
			return
		}
		if err.Error() == "connection is not active" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Connection inactive",
				"details": "Cannot mark messages as read for an inactive connection",
			})
			return
		}
		if err.Error() == "unauthorized: user is not part of this connection" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to mark messages as read in this connection",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to mark messages as read",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Messages marked as read successfully",
	})
}

// GetUnreadCount handles GET /messages/unread-count
func (h *ConnectionHandler) GetUnreadCount(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get unread count
	count, err := h.MessageUsecase.GetUnreadCount(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get unread count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": count,
	})
}
