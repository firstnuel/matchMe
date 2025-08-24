package user

import (
	"log"
	"match-me/api/middleware"
	"match-me/config"
	"match-me/ent"
	"match-me/internal/pkg/cloudinary"
	userRepo "match-me/internal/repositories/user"
	"match-me/internal/requests"
	userUsecase "match-me/internal/usecases/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase       userUsecase.UserUsecase
	validationService *requests.ValidationService
	cfg               *config.Config
}

func NewUserHandler(client *ent.Client, cfg *config.Config,
	validationService *requests.ValidationService,
	cld cloudinary.Cloudinary) *UserHandler {

	userRepo := userRepo.NewUserRepository(client)
	userUsecase := userUsecase.NewUserUsecase(userRepo, cfg.JWTSecret, cld)
	return &UserHandler{
		UserUsecase:       userUsecase,
		validationService: validationService,
		cfg:               cfg,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) *gin.Engine {
	// Public routes (no authentication required)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}

	// Protected routes (authentication required)
	usersGroup := r.Group("/users", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		usersGroup.GET("/:id", h.GetUserByID)
		usersGroup.GET("/:id/bio", h.GetUserBio)
		usersGroup.GET("/:id/profile", h.GetUserProfile)
		usersGroup.GET("/:id/distance", h.GetDistanceBetweenUsers)
	}

	// Convenience route for getting current user
	userMeGroup := r.Group("api", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		userMeGroup.GET("/me", h.GetCurrentUser)
		userMeGroup.PUT("/me", h.UpdateUser)
		userMeGroup.DELETE("/me", h.DeleteCurrentUser)
		userMeGroup.PUT("/password", h.UpdatePassword)
		userMeGroup.POST("/me/photos", h.UploadUserPhotos)
		userMeGroup.DELETE("/me/photos/:photoId", h.DeleteUserPhoto)
		userMeGroup.GET("/me/recommendations", h.GetRecommendations)
	}

	log.Println("💫 All user routes registered")
	return r
}
