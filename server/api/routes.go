package api

import (
	"match-me/config"
	"match-me/ent"
	"match-me/internal/adapters/user"
	"match-me/internal/pkg/cloudinary"

	"match-me/internal/requests"

	"log"

	"github.com/gin-gonic/gin"
)

func registerRoutes(client *ent.Client, r *gin.Engine, cfg *config.Config, cld cloudinary.Cloudinary) {

	log.Println("🚀 Registering API routes...")
	userHandler := user.NewUserHandler(client, cfg, requests.NewValidationService(), cld)
	userHandler.RegisterRoutes(r)

}
