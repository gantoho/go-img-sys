package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/service"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	service *service.ImageService
	logger  *logger.Logger
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		service: service.NewImageService(),
		logger:  logger.GetLogger(),
	}
}

// GetImage retrieves a single image by filename
func (h *ImageHandler) GetImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	filepath, err := h.service.GetImageByFilename(filename)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.File(filepath)
}

// ListAllImages returns all available images
func (h *ImageHandler) ListAllImages(ctx *gin.Context) {
	hostURL := ctx.Request.Host

	data, err := h.service.GetAllImages(hostURL)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// ListAllImagesWithMetadata returns all images with metadata
func (h *ImageHandler) ListAllImagesWithMetadata(ctx *gin.Context) {
	hostURL := ctx.Request.Host

	data, err := h.service.GetAllImagesWithMetadata(hostURL)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"total": len(data),
		"data":  data,
	})
}

// ListAllImagesPaginated returns paginated images with metadata
func (h *ImageHandler) ListAllImagesPaginated(ctx *gin.Context) {
	hostURL := ctx.Request.Host
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid page parameter", nil)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid page_size parameter", nil)
		return
	}

	data, appErr := h.service.GetAllImagesPaginated(hostURL, page, pageSize)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// GetRandomImage returns a random image filename
func (h *ImageHandler) GetRandomImage(ctx *gin.Context) {
	filename, err := h.service.GetRandomImage()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.String(http.StatusOK, filename)
}

// GetRandomImages returns multiple random images with URLs
func (h *ImageHandler) GetRandomImages(ctx *gin.Context) {
	hostURL := ctx.Request.Host
	countStr := ctx.Param("number")

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid number parameter", nil)
		return
	}

	if count > 100 {
		count = 100 // Limit to 100 images per request
	}

	images, appErr := h.service.GetRandomImages(hostURL, count)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"total": len(images),
		"data":  images,
	})
}

// UploadImage handles file uploads
func (h *ImageHandler) UploadImage(ctx *gin.Context) {
	hostURL := ctx.Request.Host

	form, err := ctx.MultipartForm()
	if err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "failed to parse form data", nil)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "no files provided", nil)
		return
	}

	// Save files
	uploadedFiles := make([]map[string]interface{}, 0)
	failedFiles := make([]map[string]string, 0)

	cfg := config.GetConfig()
	uploadDir := cfg.File.UploadDir
	dupStrategy := cfg.File.DuplicateStrategy

	for idx, file := range files {
		origName := file.Filename
		dstName := origName
		dstPath := utils.GetUploadPath(uploadDir, dstName)

		// Handle duplicates according to strategy
		if utils.FileExists(dstPath) {
			switch dupStrategy {
			case "overwrite":
				// do nothing, will overwrite
			case "reject":
				failedFiles = append(failedFiles, map[string]string{
					"filename": origName,
					"error":    "file already exists",
				})
				h.logger.Warn("Upload rejected for existing file %s", origName)
				continue
			default: // rename
				// generate unique name: name_1.ext, name_2.ext ...
				ext := filepath.Ext(origName)
				nameOnly := origName[:len(origName)-len(ext)]
				i := 1
				for {
					dstName = nameOnly + "_" + strconv.Itoa(i) + ext
					dstPath = utils.GetUploadPath(uploadDir, dstName)
					if !utils.FileExists(dstPath) {
						break
					}
					i++
				}
			}
		}

		// Ensure upload dir exists
		if err := utils.EnsureDir(uploadDir); err != nil {
			h.logger.Error("Failed to ensure upload dir: %v", err)
			failedFiles = append(failedFiles, map[string]string{
				"filename": origName,
				"error":    "server error: cannot create upload dir",
			})
			continue
		}

		if err := ctx.SaveUploadedFile(file, dstPath); err != nil {
			h.logger.Error("Failed to save file %s -> %s: %v", origName, dstPath, err)
			failedFiles = append(failedFiles, map[string]string{
				"filename": origName,
				"error":    err.Error(),
			})
			continue
		}

		uploadedFiles = append(uploadedFiles, map[string]interface{}{
			"index":    idx + 1,
			"filename": dstName,
			"size":     file.Size,
			"url":      hostURL + "/f/" + dstName,
			"progress": 100,
		})
	}

	result := map[string]interface{}{
		"message":        "Upload completed",
		"total_files":    len(files),
		"total_uploaded": len(uploadedFiles),
		"uploaded":       uploadedFiles,
	}

	if len(failedFiles) > 0 {
		result["failed"] = failedFiles
		result["total_failed"] = len(failedFiles)
	}

	utils.SuccessResponse(ctx, result)
}

// SearchImages searches and filters images
func (h *ImageHandler) SearchImages(ctx *gin.Context) {
	hostURL := ctx.Request.Host
	filename := ctx.DefaultQuery("filename", "")
	minSizeStr := ctx.DefaultQuery("min_size", "0")
	maxSizeStr := ctx.DefaultQuery("max_size", "0")
	fileType := ctx.DefaultQuery("type", "")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")

	minSize, _ := strconv.ParseInt(minSizeStr, 10, 64)
	maxSize, _ := strconv.ParseInt(maxSizeStr, 10, 64)
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	data, appErr := h.service.SearchImages(hostURL, filename, minSize, maxSize, fileType, page, pageSize)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// DeleteImage deletes a single image by filename
func (h *ImageHandler) DeleteImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if err := h.service.DeleteImage(filename); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, map[string]string{
		"message": "image deleted successfully",
	})
}

// DeleteImages deletes multiple images
func (h *ImageHandler) DeleteImages(ctx *gin.Context) {
	var req struct {
		Filenames []string `json:"filenames" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if len(req.Filenames) == 0 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "filenames list is empty", nil)
		return
	}

	result := h.service.DeleteImages(req.Filenames)
	utils.SuccessResponse(ctx, result)
}

func (h *ImageHandler) HealthCheck(ctx *gin.Context) {
	utils.SuccessResponse(ctx, map[string]string{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// CreateAPIKey creates a new API key
func (h *ImageHandler) CreateAPIKey(ctx *gin.Context) {
	var req struct {
		ExpireDays int `json:"expire_days" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}

	if req.ExpireDays < 1 || req.ExpireDays > 365 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "expire_days must be between 1 and 365", nil)
		return
	}

	keyManager := auth.GetManager()
	plainKey := keyManager.CreateKey("api-key-"+strconv.Itoa(len(keyManager.ListKeys())), req.ExpireDays)

	h.logger.Info("New API key created, expires in %d days", req.ExpireDays)

	utils.SuccessResponse(ctx, map[string]interface{}{
		"api_key":     plainKey,
		"expire_days": req.ExpireDays,
		"message":     "API key created successfully. Please save it safely!",
	})
}

// ListAPIKeys lists all API keys
func (h *ImageHandler) ListAPIKeys(ctx *gin.Context) {
	keyManager := auth.GetManager()
	keys := keyManager.ListKeys()

	utils.SuccessResponse(ctx, map[string]interface{}{
		"total": len(keys),
		"keys":  keys,
	})
}

// RevokeAPIKey revokes an API key
func (h *ImageHandler) RevokeAPIKey(ctx *gin.Context) {
	var req struct {
		APIKey string `json:"api_key" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}

	keyManager := auth.GetManager()
	if !keyManager.RevokeKey(req.APIKey) {
		utils.CustomResponse(ctx, http.StatusNotFound, "API key not found", nil)
		return
	}

	h.logger.Info("API key revoked")
	utils.SuccessResponse(ctx, map[string]string{
		"message": "API key revoked successfully",
	})
}

// GetStatistics 获取统计信息
func (h *ImageHandler) GetStatistics(ctx *gin.Context) {
	statService := service.NewStatisticsService()
	stats := statService.GetStatistics()

	utils.SuccessResponse(ctx, stats)
}

// GetDiskUsage 获取磁盘使用情况
func (h *ImageHandler) GetDiskUsage(ctx *gin.Context) {
	statService := service.NewStatisticsService()
	usage := statService.GetDiskUsage()

	utils.SuccessResponse(ctx, usage)
}

// ExportFiles 导出多个文件为ZIP
func (h *ImageHandler) ExportFiles(ctx *gin.Context) {
	var req struct {
		Filenames []string `json:"filenames" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	exportService := service.NewExportService()
	result, err := exportService.ExportMultipleFiles(req.Filenames, "./files")
	if err != nil {
		utils.CustomResponse(ctx, http.StatusInternalServerError, "export failed", nil)
		return
	}

	utils.SuccessResponse(ctx, result)
}

// ExportAllFiles 导出所有文件
func (h *ImageHandler) ExportAllFiles(ctx *gin.Context) {
	exportService := service.NewExportService()
	result, err := exportService.ExportAllFiles("./files")
	if err != nil {
		utils.CustomResponse(ctx, http.StatusInternalServerError, "export failed", nil)
		return
	}

	utils.SuccessResponse(ctx, result)
}

// Cleanup 执行清理操作
func (h *ImageHandler) Cleanup(ctx *gin.Context) {
	var req struct {
		RemoveOrphanThumbnails bool `json:"remove_orphan_thumbnails"`
		RemoveOldFiles         bool `json:"remove_old_files"`
		MaxFileAgeDays         int  `json:"max_file_age_days"`
		RemoveEmptyDirs        bool `json:"remove_empty_dirs"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	maintService := service.NewMaintenanceService()
	maxAge := time.Duration(req.MaxFileAgeDays) * 24 * time.Hour
	if maxAge == 0 {
		maxAge = 24 * time.Hour * 30 // 默认30天
	}

	result := maintService.Cleanup(service.CleanupConfig{
		RemoveOrphanThumbnails: req.RemoveOrphanThumbnails,
		RemoveOldFiles:         req.RemoveOldFiles,
		MaxFileAge:             maxAge,
		RemoveEmptyDirs:        req.RemoveEmptyDirs,
	})

	utils.SuccessResponse(ctx, result)
}

// StartThumbnailGeneration 为现有图片生成缩略图
func (h *ImageHandler) StartThumbnailGeneration(ctx *gin.Context) {
	filenames := ctx.Query("filenames")

	if filenames == "" {
		utils.CustomResponse(ctx, http.StatusBadRequest, "filenames parameter required", nil)
		return
	}

	h.logger.Info("Thumbnail generation started for: %s", filenames)
	utils.SuccessResponse(ctx, map[string]string{
		"message": "Thumbnail generation started (background task)",
	})
}

// Login handles user login and returns JWT token
func (h *ImageHandler) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// Simple authentication - in production, validate against database with hashed passwords
	var userID, role string
	switch {
	case req.Username == "admin" && req.Password == "admin123":
		userID = "1"
		role = "admin"
	case req.Username == "user" && req.Password == "user123":
		userID = "2"
		role = "user"
	default:
		utils.CustomResponse(ctx, http.StatusUnauthorized, "invalid username or password", nil)
		return
	}

	// Generate JWT token
	jwtManager := auth.GetJWTManager()
	token, err := jwtManager.GenerateToken(userID, req.Username, role)
	if err != nil {
		utils.CustomResponse(ctx, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"token":      token,
		"user_id":    userID,
		"username":   req.Username,
		"role":       role,
		"expires_at": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
}

// RefreshToken handles JWT token refresh
func (h *ImageHandler) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		utils.CustomResponse(ctx, http.StatusUnauthorized, "missing authorization header", nil)
		return
	}

	// Extract token from "Bearer <token>" format
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Refresh token
	jwtManager := auth.GetJWTManager()
	newToken, err := jwtManager.RefreshToken(tokenString)
	if err != nil {
		utils.CustomResponse(ctx, http.StatusUnauthorized, "failed to refresh token", nil)
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"token":      newToken,
		"expires_at": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
}
