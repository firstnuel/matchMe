package api

import (
	"match-me/config"
	"match-me/ent"
	"match-me/internal/adapters/user"
	"match-me/internal/pkg/cloudinary"
	userRepo "match-me/internal/repositories/user"
	"match-me/internal/requests"
	userUsecase "match-me/internal/usecases/user"

	"log"

	"github.com/gin-gonic/gin"
)

func registerRoutes(client *ent.Client, r *gin.Engine, cfg *config.Config) {

	log.Println("🚀 Registering API routes...")

	userRepo := userRepo.NewUserRepository(client)
	userUsecase := userUsecase.NewUserUsecase(userRepo, cfg.JWTSecret, cloudinary.NewCloudinary())
	userHandler := user.NewUserHandler(userUsecase, cfg, requests.NewValidationService())

	userHandler.RegisterRoutes(r)

}
