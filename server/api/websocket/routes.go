package websocket

import (
	"log"
	"match-me/api/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterWebSocketRoutes registers WebSocket routes
func (ws *WebSocketHandler) RegisterRoutes(r *gin.Engine) *gin.Engine {

	// WebSocket routes group
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.VerifyUser(ws.UserUsecase, ws.cfg.JWTSecret))
	{
		// Chat WebSocket - for real-time messaging in a specific connection
		wsGroup.GET("/chat/:connectionId", ws.HandleChatConnection)

		// Status WebSocket - for general user online/offline status
		wsGroup.GET("/status", ws.HandleStatusConnection)

		// Typing WebSocket - for typing indicators in a specific connection
		wsGroup.GET("/typing/:connectionId", ws.HandleTypingConnection)

		// Status endpoints (HTTP endpoints for debugging/admin)
		wsGroup.GET("/online-users", ws.GetOnlineUsers)
		wsGroup.GET("/connection/:connectionId/status", ws.GetConnectionStatus)
	}
	log.Println("💫 All websocket routes registered")
	return r
}
