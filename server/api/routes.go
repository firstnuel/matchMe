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

	wshub := wscore.NewHub()
	webSocketService := wscore.NewWebSocketService(wshub)

	log.Println("🚀 Registering API routes...")
	userHandler := user.NewUserHandler(client, cfg, requests.NewValidationService(), cld)
	userHandler.RegisterRoutes(r)

	connectionRepo := connections.NewConnectionRepository(client)

	webSocketHandler := websocket.NewWebSocketHandler(wshub, connectionRepo, userHandler.UserUsecase, cfg)
	webSocketHandler.RegisterRoutes(r)

	connectionHandler := connection.NewConnectionHandler(
		client, cfg, requests.NewValidationService(), cld, webSocketService, userHandler.UserUsecase)

	connectionHandler.RegisterRoutes(r)

}
