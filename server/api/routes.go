package api

import (
	"match-me/api/websocket"
	"match-me/config"
	"match-me/ent"
	"match-me/internal/adapters/connection"
	"match-me/internal/adapters/user"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/repositories/connections"
	"match-me/internal/requests"
	wscore "match-me/internal/websocket"

	"log"

	"github.com/gin-gonic/gin"
)

func registerRoutes(client *ent.Client, r *gin.Engine, cfg *config.Config, cld cloudinary.Cloudinary) {

	connectionRepo := connections.NewConnectionRepository(client)

	chatHub := wscore.NewChatHub()
	typingHub := wscore.NewTypingHub()
	statusHub := wscore.NewStatusHub()

	go chatHub.Run()
	go typingHub.Run()
	go statusHub.Run()

	webSocketService := wscore.NewWebSocketService(chatHub, typingHub, statusHub)

	log.Println("🚀 Registering API routes...")
	userHandler := user.NewUserHandler(client, cfg, connectionRepo, requests.NewValidationService(), cld)
	userHandler.RegisterRoutes(r)

	webSocketHandler := websocket.NewWebSocketHandler(chatHub, typingHub, statusHub, connectionRepo, userHandler.UserUsecase, cfg)
	webSocketHandler.RegisterRoutes(r)

	connectionHandler := connection.NewConnectionHandler(
		client, cfg, requests.NewValidationService(), cld, webSocketService, userHandler.UserUsecase)

	connectionHandler.RegisterRoutes(r)

}
