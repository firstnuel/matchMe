package api

import (
	"match-me/api/websocket"
	"match-me/config"
	"match-me/ent"
	"match-me/internal/adapters/connection"
	"match-me/internal/adapters/user"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/repositories/connections"
	"match-me/internal/repositories/interactions"
	"match-me/internal/requests"
	inUc "match-me/internal/usecases/interactions"
	wscore "match-me/internal/websocket"

	"log"

	"github.com/gin-gonic/gin"
)

func registerRoutes(client *ent.Client, r *gin.Engine, cfg *config.Config, cld cloudinary.Cloudinary) {

	connectionRepo := connections.NewConnectionRepository(client)
	connectionReqRepo := connections.NewConnectionRequestRepository(client)
	interactionRepo := interactions.NewUserInteractionRepository(client)

	chatHub := wscore.NewChatHub()
	typingHub := wscore.NewTypingHub()
	statusHub := wscore.NewStatusHub()

	go chatHub.Run()
	go typingHub.Run()
	go statusHub.Run()

	webSocketService := wscore.NewWebSocketService(chatHub, typingHub, statusHub)
	interactionService := inUc.NewUserInteractionUsecase(interactionRepo)
	validationService := requests.NewValidationService()

	log.Println("🚀 Registering API routes...")
	userHandler := user.NewUserHandler(
		client,
		cfg,
		connectionRepo,
		connectionReqRepo,
		interactionService,
		validationService,
		cld,
	)
	userHandler.RegisterRoutes(r)

	webSocketHandler := websocket.NewWebSocketHandler(
		chatHub,
		typingHub,
		statusHub,
		connectionRepo,
		userHandler.UserUsecase,
		cfg,
	)
	webSocketHandler.RegisterRoutes(r)

	connectionHandler := connection.NewConnectionHandler(
		client,
		cfg,
		validationService,
		cld,
		webSocketService,
		userHandler.UserUsecase,
		interactionService,
	)
	connectionHandler.RegisterRoutes(r)

}
