package user

import (
	"match-me/api/middleware"
	"match-me/config"
	"match-me/internal/requests"
	"match-me/internal/usecases/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase       user.UserUsecase
	validationService *requests.ValidationService
	cfg               *config.Config
}

func NewUserHandler(userUsecase user.UserUsecase, cfg *config.Config, validationService *requests.ValidationService) *UserHandler {
	return &UserHandler{
		userUsecase:       userUsecase,
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
	usersGroup := r.Group("/users", middleware.VerifyUser(h.userUsecase, h.cfg))
	{
		usersGroup.GET("/:id", h.GetUserByID)
		usersGroup.GET("/:id/bio", h.GetUserBio)
		usersGroup.GET("/:id/profile", h.GetUserProfile)
	}

	// Convenience route for getting current user
	userMeGroup := r.Group("/me", middleware.VerifyUser(h.userUsecase, h.cfg))
	{
		userMeGroup.GET("/", h.GetCurrentUser)
		userMeGroup.PUT("/", h.UpdateCurrentUser)
		userMeGroup.DELETE("/", h.DeleteCurrentUser)
		userMeGroup.PUT("/password", h.UpdateCurrentUserPassword)
		userMeGroup.POST("/photos", h.UploadCurrentUserPhotos)
	}

	return r
}
