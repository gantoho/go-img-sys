package router

import (
	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/handler"
	"github.com/gantoho/go-img-sys/internal/middleware"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Apply global middleware
	router.Use(middleware.RequestTimingMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.CORSMiddleware())

	imageHandler := handler.NewImageHandler()

	// Initialize JWT
	cfg := config.GetConfig()
	auth.InitJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpire)
	jwtManager := auth.GetJWTManager()

	// API routes
	api := router.Group("/api")

	// Auth routes - no protection required
	auth := api.Group("/auth")
	{
		// Login endpoint
		auth.POST("/login", imageHandler.Login)

		// Refresh token endpoint (optional JWT for convenience)
		auth.POST("/refresh", middleware.OptionalJWTMiddleware(jwtManager), imageHandler.RefreshToken)
	}

	// v1 routes - public read operations
	v1 := api.Group("/v1")
	{
		// Health check
		v1.GET("/health", imageHandler.HealthCheck)

		// Image read operations (public)
		v1.GET("/images", imageHandler.ListAllImages)
		v1.GET("/images/metadata", imageHandler.ListAllImagesWithMetadata)
		v1.GET("/images/paginated", imageHandler.ListAllImagesPaginated)
		v1.GET("/images/search", imageHandler.SearchImages)
		v1.GET("/images/random", imageHandler.GetRandomImage)
		v1.GET("/images/random/:number", imageHandler.GetRandomImages)
	}

	// v1 protected routes - write operations require JWT
	v1Protected := api.Group("/v1")
	v1Protected.Use(middleware.AuthMiddleware())
	{
		v1Protected.POST("/images/upload", imageHandler.UploadImage)
		v1Protected.DELETE("/images/:filename", imageHandler.DeleteImage)
		v1Protected.POST("/images/delete", imageHandler.DeleteImages)
	}

	// v1 admin routes - requires JWT with admin role
	v1Admin := api.Group("/v1/admin")
	v1Admin.Use(middleware.AuthMiddleware())
	{
		v1Admin.GET("/api-keys", imageHandler.ListAPIKeys)
		v1Admin.DELETE("/api-keys", imageHandler.RevokeAPIKey)
	}

	// Public endpoint for creating API keys (no authentication)
	api.POST("/v1/admin/api-keys", imageHandler.CreateAPIKey)
	api.POST("/v1/admin/api-keys/validate", imageHandler.ValidateAPIKey)

	// v1 utility routes - statistics, export, cleanup (public read, protected write)
	v1Util := api.Group("/v1/util")
	{
		// Public read operations
		v1Util.GET("/statistics", imageHandler.GetStatistics)
		v1Util.GET("/disk-usage", imageHandler.GetDiskUsage)
	}

	// v1 utility protected routes
	v1UtilProtected := api.Group("/v1/util")
	v1UtilProtected.Use(middleware.AuthMiddleware())
	{
		v1UtilProtected.POST("/export", imageHandler.ExportFiles)
		v1UtilProtected.POST("/export-all", imageHandler.ExportAllFiles)
		v1UtilProtected.POST("/cleanup", imageHandler.Cleanup)
		v1UtilProtected.POST("/generate-thumbnails", imageHandler.StartThumbnailGeneration)
	}

	// Direct file access
	router.GET("/f/:filename", imageHandler.GetImage)

	router.GET("/bgimg", imageHandler.GetRandomImage)

	// Legacy support - public read operations
	legacyV1 := router.Group("/v1")
	{
		legacyV1.GET("/", imageHandler.HealthCheck)
		legacyV1.GET("/all", imageHandler.ListAllImages)
		legacyV1.GET("/get/:number", imageHandler.GetRandomImages)
		legacyV1.GET("/bgimg", imageHandler.GetRandomImage)
	}

	// Legacy protected routes - also support JWT now
	legacyV1Protected := router.Group("/v1")
	legacyV1Protected.Use(middleware.AuthMiddleware())
	{
		legacyV1Protected.POST("/upload", imageHandler.UploadImage)
	}
}
