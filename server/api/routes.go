package api

import (
	"match-me/config"
	"match-me/ent"
	"match-me/internal/adapters/user"

	"match-me/internal/requests"

	"log"

	"github.com/gin-gonic/gin"
)

func registerRoutes(client *ent.Client, r *gin.Engine, cfg *config.Config) {

	log.Println("🚀 Registering API routes...")
	userHandler := user.NewUserHandler(client, cfg, requests.NewValidationService())
	userHandler.RegisterRoutes(r)

}
