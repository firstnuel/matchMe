package connection

import (
	"log"
	"match-me/api/middleware"
	"match-me/config"
	"match-me/ent"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/repositories/connections"
	"match-me/internal/requests"
	connectionUsecases "match-me/internal/usecases/connections"
	userUsecase "match-me/internal/usecases/user"
	"match-me/internal/websocket"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	ConnectionUsecase        connectionUsecases.ConnectionUsecase
	ConnectionRequestUsecase connectionUsecases.ConnectionRequestUsecase
	MessageUsecase           connectionUsecases.MessageUsecase
	UserUsecase              userUsecase.UserUsecase
	validationService        *requests.ValidationService
	cfg                      *config.Config
}

func NewConnectionHandler(
	client *ent.Client,
	cfg *config.Config,
	validationService *requests.ValidationService,
	cld cloudinary.Cloudinary,
	wsService *websocket.WebSocketService,
	userUsecase userUsecase.UserUsecase,
) *ConnectionHandler {

	// Create repositories
	connectionRepo := connections.NewConnectionRepository(client)
	requestRepo := connections.NewConnectionRequestRepository(client)
	messageRepo := connections.NewMessageRepository(client)

	// Create usecases
	connectionUsecase := connectionUsecases.NewConnectionUsecase(connectionRepo)
	connectionRequestUsecase := connectionUsecases.NewConnectionRequestUsecase(requestRepo, connectionRepo, wsService)
	messageUsecase := connectionUsecases.NewMessageUsecase(messageRepo, connectionRepo, cld, wsService)

	return &ConnectionHandler{
		ConnectionUsecase:        connectionUsecase,
		ConnectionRequestUsecase: connectionRequestUsecase,
		MessageUsecase:           messageUsecase,
		UserUsecase:              userUsecase,
		validationService:        validationService,
		cfg:                      cfg,
	}
}

func (h *ConnectionHandler) RegisterRoutes(r *gin.Engine) *gin.Engine {
	// All connection routes require authentication
	connectionGroup := r.Group("/connections", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		// Connection management
		connectionGroup.GET("/", h.GetUserConnections)
		connectionGroup.DELETE("/:connectionId", h.DeleteConnection)
	}

	// Connection request routes
	requestGroup := r.Group("/connection-requests", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		requestGroup.POST("/", h.SendConnectionRequest)
		requestGroup.GET("/", h.GetPendingRequests)
		requestGroup.PUT("/:requestId/accept", h.AcceptRequest)
		requestGroup.PUT("/:requestId/decline", h.DeclineRequest)
	}

	// Message routes
	messageGroup := r.Group("/messages", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		messageGroup.POST("/text", h.SendTextMessage)
		messageGroup.POST("/media", h.SendMediaMessage)
		messageGroup.GET("/connection/:connectionId", h.GetConnectionMessages)
		messageGroup.PUT("/connection/:connectionId/read", h.MarkMessagesAsRead)
		messageGroup.GET("/unread-count", h.GetUnreadCount)
		messageGroup.GET("/chat-list", h.GetChatList)
	}

	log.Println("💫 All connection routes registered")
	return r
}