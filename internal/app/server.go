package app

import (
	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/router"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	engine *gin.Engine
}

func New() *Server {
	cfg := config.Init()
	log := logger.Init()

	return &Server{
		config: cfg,
		logger: log,
		engine: gin.Default(),
	}
}

// Start initializes and starts the server
func (s *Server) Start() {
	// Ensure upload directory exists
	if err := utils.EnsureDir(s.config.File.UploadDir); err != nil {
		s.logger.Fatal("Failed to create upload directory: %v", err)
	}

	// Initialize API key manager with default keys
	keyManager := auth.GetManager()
	keyManager.InitDefaultKeys()
	go keyManager.CleanupExpiredKeys() // Start background cleanup

	s.logger.Info("API Key Manager initialized with default keys")
	s.logger.Info("Default API Keys: demo-key-12345 (30 days), test-key-67890 (7 days)")

	// Register routes
	router.RegisterRoutes(s.engine)

	// Print startup info
	s.logger.Info("Starting Image Server on %s", s.config.Server.Port)
	s.logger.Info("Upload directory: %s", s.config.File.UploadDir)

	// Start server
	addr := s.config.Server.Port
	if err := s.engine.Run(addr); err != nil {
		s.logger.Fatal("Server failed to start: %v", err)
	}
}

// Close gracefully closes the server
func (s *Server) Close() {
	s.logger.Info("Shutting down server...")
	defer s.logger.Close()
}
