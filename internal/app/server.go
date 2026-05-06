package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/router"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     *config.Config
	logger     *logger.Logger
	engine     *gin.Engine
	httpServer *http.Server
	stopCh     chan struct{}
}

func New(cfg *config.Config) *Server {
	log := logger.Init()

	if cfg.Server.Env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &Server{
		config: cfg,
		logger: log,
		engine: gin.New(),
		stopCh: make(chan struct{}),
	}

	router.RegisterRoutes(srv.engine)
	return srv
}

func (s *Server) Start() {
	if err := utils.EnsureDir(s.config.File.UploadDir); err != nil {
		s.logger.Fatal("Failed to create upload directory: %v", err)
	}

	keyManager := auth.GetManager()
	if err := keyManager.LoadFromFile(auth.DefaultStorePath()); err != nil {
		keyManager.InitDefaultKeys()
		if err := keyManager.SaveToFile(auth.DefaultStorePath()); err != nil {
			s.logger.Error("Failed to save default API keys: %v", err)
		}
	}
	go keyManager.CleanupExpiredKeys()

	if s.config.Server.Env != "release" {
		s.logger.Info("API Key Manager initialized with default keys")
		s.logger.Info("Default API Keys: demo-key-12345 (30 days), test-key-67890 (7 days)")
	}

	s.startFileWatcher()

	addr := s.config.Server.Port
	s.logger.Info("Starting Image Server on %s", addr)
	s.logger.Info("Upload directory: %s", s.config.File.UploadDir)
	s.logger.Info("Environment: %s", s.config.Server.Env)
	s.logger.Info("Metrics: http://localhost%s/metrics", addr)
	s.logger.Info("Swagger UI: http://localhost%s/swagger/index.html", addr)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Shutting down server...")

	close(s.stopCh)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown: %v", err)
	}
	s.Close()
}

func (s *Server) startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Warn("File watcher not available: %v", err)
		return
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-s.stopCh:
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					s.logger.Debug("File change detected: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				s.logger.Error("File watcher error: %v", err)
			}
		}
	}()

	if err := watcher.Add(s.config.File.UploadDir); err != nil {
		s.logger.Warn("Cannot watch upload dir: %v", err)
	}

	if err := watcher.Add("."); err != nil {
		s.logger.Warn("Cannot watch config dir: %v", err)
	} else {
		s.logger.Info("Config file watcher started")
	}
}

func (s *Server) Close() {
	s.logger.Info("Server stopped")
	s.logger.Close()
}
