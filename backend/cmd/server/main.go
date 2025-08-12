package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds server configuration
type Config struct {
	Debug bool
	Port  string
}

// Server represents the main server instance
type Server struct {
	config *Config
	logger *zap.SugaredLogger
	router *gin.Engine
}

// NewServer creates a new server instance
func NewServer(config *Config) *Server {
	// Configure logger
	var logger *zap.Logger
	if config.Debug {
		// Development logger with colored output
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	} else {
		// Production logger
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	// Set Gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		config: config,
		logger: sugaredLogger,
		router: gin.Default(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Health check endpoint (always available)
	s.router.GET("/health", s.healthCheck)

	// API routes
	api := s.router.Group("/api")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.handleRegister)
			auth.POST("/login", s.handleLogin)
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/:id", s.getUser)
				users.GET("/:id/profile", s.getUserProfile)
				users.GET("/:id/bio", s.getUserBio)
			}

			// Current user routes
			me := protected.Group("/me")
			{
				me.GET("", s.getMe)
				me.GET("/profile", s.getMyProfile)
				me.GET("/bio", s.getMyBio)
			}

			// Recommendations and connections
			protected.GET("/recommendations", s.getRecommendations)
			protected.GET("/connections", s.getConnections)
		}
	}

	// Debug routes (only in debug mode)
	if s.config.Debug {
		s.setupDebugRoutes()
	}
}

// setupDebugRoutes adds debug-only routes
func (s *Server) setupDebugRoutes() {
	s.logger.Info("Setting up debug routes")

	// Debug endpoints
	s.router.GET("/debug/health", s.healthCheck)
	s.router.GET("/debug/status", s.debugStatus)

	// Development tools
	s.router.GET("/debug/users", s.listUsers) // Only in debug mode
}

// authMiddleware provides JWT authentication middleware
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT token validation
		// For now, just pass through
		c.Next()
	}
}

// Handler stubs (to be implemented)
func (s *Server) handleRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - to be implemented"})
}

func (s *Server) handleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - to be implemented"})
}

func (s *Server) getUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint - to be implemented"})
}

func (s *Server) getUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user profile endpoint - to be implemented"})
}

func (s *Server) getUserBio(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user bio endpoint - to be implemented"})
}

func (s *Server) getMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get me endpoint - to be implemented"})
}

func (s *Server) getMyProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get my profile endpoint - to be implemented"})
}

func (s *Server) getMyBio(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get my bio endpoint - to be implemented"})
}

func (s *Server) getRecommendations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get recommendations endpoint - to be implemented"})
}

func (s *Server) getConnections(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get connections endpoint - to be implemented"})
}

func (s *Server) listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint - to be implemented"})
}

// Health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "match-me",
		"timestamp": "2024-01-01T00:00:00Z",
	})
}

// Debug status endpoint
func (s *Server) debugStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "debug",
		"service":    "match-me",
		"debug_mode": s.config.Debug,
		"timestamp":  "2024-01-01T00:00:00Z",
	})
}

// Start initializes and starts the server
func (s *Server) Start() error {
	s.logger.Infof("Starting Match-Me server on port %s (debug: %v)", s.config.Port, s.config.Debug)

	// Start server
	return s.router.Run(":" + s.config.Port)
}

func main() {
	// Parse command line flags
	debug := flag.Bool("d", false, "Enable debug mode")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Load configuration
	config := &Config{
		Debug: *debug,
		Port:  *port,
	}

	// Create and start server
	server := NewServer(config)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
