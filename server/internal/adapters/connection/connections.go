package connection

import (
	"net/http"

	"match-me/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserConnections handles GET /connections/details
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

// GetUserConnectionsIds handles GET /connections
func (h *ConnectionHandler) GetUserConnectionsIds(c *gin.Context) {
	// Get authenticated user
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user connection IDs
	connectionIDs, err := h.ConnectionUsecase.GetUserConnectionsIds(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get connection IDs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connection_ids": connectionIDs,
		"count":          len(connectionIDs),
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

func (h *ConnectionHandler) SkipConnection(c *gin.Context) {
	user, _ := middleware.GetUserFromGinContext(c)

	var req struct {
		TargetUserID string `json:"target_userId" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	targetUUID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	if err := h.InteractionUsecase.RecordSkippedProfile(c.Request.Context(), user.ID, targetUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to record skipped profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile skipped successfully"})
}
