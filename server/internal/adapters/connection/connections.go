package connection

import (
	"net/http"

	"match-me/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserConnections handles GET /connections
func (h *ConnectionHandler) GetUserConnections(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user connections
	connections, err := h.ConnectionUsecase.GetUserConnections(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get connections",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connections": connections,
		"count":       len(connections),
	})
}

// DeleteConnection handles DELETE /connections/:connectionId
func (h *ConnectionHandler) DeleteConnection(c *gin.Context) {
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

	// Delete connection
	err = h.ConnectionUsecase.DeleteConnection(c.Request.Context(), user.ID, connectionID)
	if err != nil {
		if err.Error() == "connection not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Connection not found",
				"details": "The specified connection does not exist",
			})
			return
		}
		if err.Error() == "unauthorized: user is not part of this connection" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"details": "You are not authorized to delete this connection",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete connection",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Connection deleted successfully",
	})
}