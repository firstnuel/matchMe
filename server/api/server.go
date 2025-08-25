package api

import (
	"match-me/api/middleware"
	"match-me/config"
	"match-me/ent"
	"match-me/internal/pkg/cloudinary"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(client *ent.Client, cfg *config.Config, cld cloudinary.Cloudinary) *http.Server {

	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middlewares
	router.Use(middleware.CORSMiddleware(cfg.ClientAddr))
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.Ping())

	// Register routes
	registerRoutes(client, router, cfg, cld)

	// HTTP server setup
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	return srv
}
