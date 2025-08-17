package user

import (
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

func NewUserHandler(client *ent.Client, cfg *config.Config, validationService *requests.ValidationService) *UserHandler {

	userRepo := userRepo.NewUserRepository(client)
	userUsecase := userUsecase.NewUserUsecase(userRepo, cfg.JWTSecret, cloudinary.NewCloudinary())
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
	}

	// Convenience route for getting current user
	userMeGroup := r.Group("/me", middleware.VerifyUser(h.UserUsecase, h.cfg.JWTSecret))
	{
		userMeGroup.GET("/", h.GetCurrentUser)
		userMeGroup.PUT("/", h.UpdateCurrentUser)
		userMeGroup.DELETE("/", h.DeleteCurrentUser)
		userMeGroup.PUT("/password", h.UpdateCurrentUserPassword)
		userMeGroup.POST("/photos", h.UploadCurrentUserPhotos)
	}

	return r
}
