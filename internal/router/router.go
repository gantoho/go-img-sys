package router

import (
	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/handler"
	"github.com/gantoho/go-img-sys/internal/middleware"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine) {
	imageHandler := handler.NewImageHandler()
	cfg := config.GetConfig()
	auth.InitJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpire)
	jwtManager := auth.GetJWTManager()
	registerRoutes(router, imageHandler, jwtManager)
}

func registerRoutes(router *gin.Engine, imageHandler *handler.ImageHandler, jwtManager *auth.JWTManager) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.RequestTimingMiddleware())
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.CORSMiddleware())

	router.GET("/metrics", middleware.PrometheusHandler())

	api := router.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", imageHandler.Login)
		authGroup.POST("/refresh", middleware.OptionalJWTMiddleware(jwtManager), imageHandler.RefreshToken)
	}

	v1 := api.Group("/v1")
	{
		v1.GET("/health", imageHandler.HealthCheck)
		v1.GET("/images", imageHandler.ListAllImages)
		v1.GET("/images/metadata", imageHandler.ListAllImagesWithMetadata)
		v1.GET("/images/paginated", imageHandler.ListAllImagesPaginated)
		v1.GET("/images/search", imageHandler.SearchImages)
		v1.GET("/images/random", imageHandler.GetRandomImage)
		v1.GET("/images/random/:number", imageHandler.GetRandomImages)
	}

	v1Protected := api.Group("/v1")
	v1Protected.Use(middleware.AuthMiddleware())
	{
		v1Protected.POST("/images/upload", imageHandler.UploadImage)
		v1Protected.DELETE("/images/:filename", imageHandler.DeleteImage)
		v1Protected.POST("/images/delete", imageHandler.DeleteImages)
		v1Protected.POST("/images/upload/chunk/init", imageHandler.InitChunkUpload)
		v1Protected.POST("/images/upload/chunk", imageHandler.UploadChunk)
		v1Protected.POST("/images/upload/chunk/complete", imageHandler.CompleteChunkUpload)
	}

	v1Admin := api.Group("/v1/admin")
	v1Admin.Use(middleware.AuthMiddleware())
	{
		v1Admin.GET("/api-keys", imageHandler.ListAPIKeys)
		v1Admin.DELETE("/api-keys", imageHandler.RevokeAPIKey)
	}

	api.POST("/v1/admin/api-keys", imageHandler.CreateAPIKey)
	api.POST("/v1/admin/api-keys/validate", imageHandler.ValidateAPIKey)

	v1Util := api.Group("/v1/util")
	{
		v1Util.GET("/statistics", imageHandler.GetStatistics)
		v1Util.GET("/disk-usage", imageHandler.GetDiskUsage)
	}

	v1UtilProtected := api.Group("/v1/util")
	v1UtilProtected.Use(middleware.AuthMiddleware())
	{
		v1UtilProtected.POST("/export", imageHandler.ExportFiles)
		v1UtilProtected.POST("/export-all", imageHandler.ExportAllFiles)
		v1UtilProtected.POST("/cleanup", imageHandler.Cleanup)
		v1UtilProtected.POST("/generate-thumbnails", imageHandler.StartThumbnailGeneration)
	}

	v1System := api.Group("/v1/system")
	v1System.Use(middleware.AuthMiddleware())
	{
		v1System.POST("/reload", imageHandler.ReloadConfig)
	}

	router.GET("/f/:filename", imageHandler.GetImage)
	router.GET("/bgimg", imageHandler.GetRandomImage)

	legacyV1 := router.Group("/v1")
	{
		legacyV1.GET("/", imageHandler.HealthCheck)
		legacyV1.GET("/all", imageHandler.ListAllImages)
		legacyV1.GET("/get/:number", imageHandler.GetRandomImages)
		legacyV1.GET("/bgimg", imageHandler.GetRandomImage)
	}

	legacyV1Protected := router.Group("/v1")
	legacyV1Protected.Use(middleware.AuthMiddleware())
	{
		legacyV1Protected.POST("/upload", imageHandler.UploadImage)
	}
}
