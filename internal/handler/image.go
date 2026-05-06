package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/internal/service"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/imageutil"
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
// @Summary      Get image file
// @Description  直接通过文件名获取图片文件
// @Tags         Files
// @Produce      image/*
// @Param        filename path string true "图片文件名"
// @Success      200 {file} file "图片文件"
// @Failure      404 {object} utils.Response "文件不存在"
// @Router       /f/{filename} [get]
func (h *ImageHandler) GetImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	filepath, err := h.service.GetImageByFilename(ctx.Request.Context(), filename)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.File(filepath)
}

// ListAllImages returns all available images
// @Summary      List all images
// @Description  获取所有图片的 URL 列表
// @Tags         Images
// @Produce      json
// @Success      200 {object} utils.Response{data=service.ImageData} "图片列表"
// @Router       /api/v1/images [get]
func (h *ImageHandler) ListAllImages(ctx *gin.Context) {
	hostURL := ctx.Request.Host

	data, err := h.service.GetAllImages(ctx.Request.Context(), hostURL)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// ListAllImagesWithMetadata returns all images with metadata
// @Summary      List all images with metadata
// @Description  获取所有图片及元数据（大小、MIME类型、修改时间等）
// @Tags         Images
// @Produce      json
// @Success      200 {object} utils.Response{data=object} "含元数据的图片列表"
// @Router       /api/v1/images/metadata [get]
func (h *ImageHandler) ListAllImagesWithMetadata(ctx *gin.Context) {
	hostURL := ctx.Request.Host

	data, err := h.service.GetAllImagesWithMetadata(ctx.Request.Context(), hostURL)
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
// @Summary      List images with pagination
// @Description  分页获取图片列表，支持自定义页码和每页数量
// @Tags         Images
// @Produce      json
// @Param        page query int false "页码，从1开始" default(1)
// @Param        page_size query int false "每页数量，最大100" default(20)
// @Success      200 {object} utils.Response{data=service.PaginatedImageData} "分页结果"
// @Router       /api/v1/images/paginated [get]
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

	data, appErr := h.service.GetAllImagesPaginated(ctx.Request.Context(), hostURL, page, pageSize)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// GetRandomImage returns a random image filename
// @Summary      Get random image
// @Description  获取随机一张图片的文件名
// @Tags         Images
// @Produce      plain
// @Success      200 {string} string "随机图片文件名"
// @Router       /api/v1/images/random [get]
func (h *ImageHandler) GetRandomImage(ctx *gin.Context) {
	filename, err := h.service.GetRandomImage(ctx.Request.Context())
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	ctx.String(http.StatusOK, filename)
}

// GetRandomImages returns multiple random images with URLs
// @Summary      Get N random images
// @Description  获取指定数量的随机图片 URL 列表，最多100张
// @Tags         Images
// @Produce      json
// @Param        number path int true "随机图片数量" maximum(100)
// @Success      200 {object} utils.Response{data=object} "随机图片列表"
// @Router       /api/v1/images/random/{number} [get]
func (h *ImageHandler) GetRandomImages(ctx *gin.Context) {
	hostURL := ctx.Request.Host
	countStr := ctx.Param("number")

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid number parameter", nil)
		return
	}

	if count > 100 {
		count = 100
	}

	images, appErr := h.service.GetRandomImages(ctx.Request.Context(), hostURL, count)
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
// @Summary      Upload images
// @Description  上传图片文件，支持多文件上传。使用 API Key 认证。
// @Tags         Images
// @Accept       multipart/form-data
// @Produce      json
// @Param        files formData file true "图片文件，支持多文件" 
// @Success      200 {object} utils.Response{data=object} "上传结果"
// @Failure      400 {object} utils.Response "请求错误"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/upload [post]
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

	uploadedFiles := make([]map[string]interface{}, 0)
	failedFiles := make([]map[string]string, 0)

	cfg := config.GetConfig()
	uploadDir := cfg.File.UploadDir
	dupStrategy := cfg.File.DuplicateStrategy

	for idx, file := range files {
		origName := file.Filename
		dstName := origName
		dstPath := utils.GetUploadPath(uploadDir, dstName)

		if utils.FileExists(dstPath) {
			switch dupStrategy {
			case "overwrite":
			case "reject":
				failedFiles = append(failedFiles, map[string]string{
					"filename": origName,
					"error":    "file already exists",
				})
				h.logger.Warn("Upload rejected for existing file %s", origName)
				continue
			default:
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

	if len(uploadedFiles) > 0 {
		h.service.InvalidateCache()
	}

	utils.SuccessResponse(ctx, result)
}

// SearchImages searches and filters images
// @Summary      Search images
// @Description  按文件名、大小范围、文件类型搜索过滤图片，支持分页
// @Tags         Images
// @Produce      json
// @Param        filename query string false "按文件名搜索（模糊匹配）"
// @Param        min_size query int false "文件大小下限（字节）"
// @Param        max_size query int false "文件大小上限（字节）"
// @Param        type query string false "文件类型后缀（如 jpg、png）"
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(20)
// @Success      200 {object} utils.Response{data=service.PaginatedImageData} "搜索结果"
// @Router       /api/v1/images/search [get]
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

	data, appErr := h.service.SearchImages(ctx.Request.Context(), hostURL, filename, minSize, maxSize, fileType, page, pageSize)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}

	utils.SuccessResponse(ctx, data)
}

// DeleteImage deletes a single image by filename
// @Summary      Delete image
// @Description  根据文件名删除单个图片
// @Tags         Images
// @Produce      json
// @Param        filename path string true "要删除的文件名"
// @Success      200 {object} utils.Response "删除成功"
// @Failure      404 {object} utils.Response "文件不存在"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/{filename} [delete]
func (h *ImageHandler) DeleteImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if err := h.service.DeleteImage(ctx.Request.Context(), filename); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, map[string]string{
		"message": "image deleted successfully",
	})
}

// DeleteImages deletes multiple images
// @Summary      Batch delete images
// @Description  批量删除图片文件
// @Tags         Images
// @Accept       json
// @Produce      json
// @Param        body body object true "要删除的文件名列表" SchemaExample({"filenames":["image1.jpg","image2.png"]})
// @Success      200 {object} utils.Response{data=object} "批量删除结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/delete [post]
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

	result := h.service.DeleteImages(ctx.Request.Context(), req.Filenames)
	utils.SuccessResponse(ctx, result)
}

// HealthCheck returns server health status
// @Summary      Health check
// @Description  服务健康检查端点，包含磁盘空间检测
// @Tags         System
// @Produce      json
// @Success      200 {object} utils.Response{data=object} "服务状态"
// @Router       /api/v1/health [get]
func (h *ImageHandler) HealthCheck(ctx *gin.Context) {
	diskOK := true
	var freeSpace int64

	statService := service.NewStatisticsService()
	usage := statService.GetDiskUsage(ctx.Request.Context())

	info := map[string]interface{}{
		"status":       "ok",
		"version":      utils.Version,
		"environment":  config.GetConfig().Server.Env,
		"total_files":  usage.UsedSpace,
		"disk_usage":   usage.Percentage,
	}

	if usage.Percentage > 90 {
		diskOK = false
	}

	if freeSpace > 0 {
		info["free_space"] = freeSpace
	}

	if !diskOK {
		info["status"] = "degraded"
		info["message"] = "disk usage exceeds 90%"
	}

	_ = freeSpace

	utils.SuccessResponse(ctx, info)
}

// CreateAPIKey creates a new API key
// @Summary      Create API key
// @Description  创建新的 API Key，可指定过期天数（1-365天）。此接口无需认证。
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        expire_days body int false "过期天数（1-365），默认7天" SchemaExample(30)
// @Success      200 {object} utils.Response{data=object} "创建成功，返回明文 Key"
// @Router       /api/v1/admin/api-keys [post]
func (h *ImageHandler) CreateAPIKey(ctx *gin.Context) {
	expireDays := 7

	var req struct {
		ExpireDays int `json:"expire_days"`
	}
	_ = ctx.ShouldBindJSON(&req)
	if req.ExpireDays != 0 {
		expireDays = req.ExpireDays
	} else {
		if v := ctx.Query("expire_days"); v != "" {
			if d, err := strconv.Atoi(v); err == nil {
				expireDays = d
			}
		} else {
			if v := ctx.PostForm("expire_days"); v != "" {
				if d, err := strconv.Atoi(v); err == nil {
					expireDays = d
				}
			}
		}
	}

	if expireDays < 1 || expireDays > 365 {
		utils.CustomResponse(ctx, http.StatusBadRequest, "expire_days must be between 1 and 365", nil)
		return
	}

	keyManager := auth.GetManager()
	plainKey := keyManager.CreateRandomKey(expireDays)
	keyHash := auth.GenerateKey(plainKey)
	if err := keyManager.SaveToFile(auth.DefaultStorePath()); err != nil {
		h.logger.Error("Failed to save API keys: %v", err)
	}

	h.logger.Info("New API key created, expires in %d days", expireDays)

	utils.SuccessResponse(ctx, map[string]interface{}{
		"api_key":     plainKey,
		"key_hash":    keyHash,
		"expire_days": expireDays,
		"message":     "API key created successfully. Please save it safely!",
	})
}

// ListAPIKeys lists all API keys
// @Summary      List API keys
// @Description  列出所有 API Key 的摘要信息（不返回明文 Key）
// @Tags         Admin
// @Produce      json
// @Success      200 {object} utils.Response{data=object} "Key 列表"
// @Security     ApiKeyAuth
// @Router       /api/v1/admin/api-keys [get]
func (h *ImageHandler) ListAPIKeys(ctx *gin.Context) {
	keyManager := auth.GetManager()
	keys := keyManager.ListKeys()

	utils.SuccessResponse(ctx, map[string]interface{}{
		"total": len(keys),
		"keys":  keys,
	})
}

// RevokeAPIKey revokes an API key
// @Summary      Revoke API key
// @Description  撤销指定的 API Key，使其立即失效
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        body body object true "要撤销的 API Key" SchemaExample({"api_key":"demo-key-12345"})
// @Success      200 {object} utils.Response "撤销成功"
// @Failure      404 {object} utils.Response "Key 未找到"
// @Security     ApiKeyAuth
// @Router       /api/v1/admin/api-keys [delete]
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
	if err := keyManager.SaveToFile(auth.DefaultStorePath()); err != nil {
		h.logger.Error("Failed to save API keys: %v", err)
	}

	h.logger.Info("API key revoked")
	utils.SuccessResponse(ctx, map[string]string{
		"message": "API key revoked successfully",
	})
}

// ValidateAPIKey validates an API key
// @Summary      Validate API key
// @Description  校验 API Key 的有效性、过期状态等信息
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        api_key query string false "要校验的明文 API Key"
// @Param        key_hash query string false "要校验的 Key 哈希值"
// @Success      200 {object} utils.Response{data=object} "校验结果"
// @Router       /api/v1/admin/api-keys/validate [post]
func (h *ImageHandler) ValidateAPIKey(ctx *gin.Context) {
	var req struct {
		APIKey  string `json:"api_key"`
		KeyHash string `json:"key_hash"`
	}
	_ = ctx.ShouldBindJSON(&req)
	if req.APIKey == "" {
		req.APIKey = ctx.Query("api_key")
	}
	if req.KeyHash == "" {
		req.KeyHash = ctx.Query("key_hash")
	}
	keyManager := auth.GetManager()
	var valid bool
	var info *auth.APIKey
	if req.APIKey != "" {
		valid = keyManager.ValidateKey(req.APIKey)
		info = keyManager.GetKeyInfo(req.APIKey)
	} else if req.KeyHash != "" {
		valid = keyManager.ValidateKeyHash(req.KeyHash)
		info = keyManager.GetKeyInfoByHash(req.KeyHash)
	} else {
		utils.CustomResponse(ctx, http.StatusBadRequest, "api_key or key_hash required", nil)
		return
	}
	if info == nil {
		utils.CustomResponse(ctx, http.StatusNotFound, "API key not found", nil)
		return
	}
	utils.SuccessResponse(ctx, map[string]interface{}{
		"valid":      valid,
		"active":     info.Active,
		"is_expired": time.Now().After(info.ExpiresAt),
		"expires_at": info.ExpiresAt.Unix(),
		"created_at": info.CreatedAt.Unix(),
		"key_hash":   info.Key,
	})
}

// GetStatistics returns file statistics
// @Summary      Get file statistics
// @Description  获取文件统计信息：总数、总大小、格式分布、最大文件等
// @Tags         Utility
// @Produce      json
// @Success      200 {object} utils.Response{data=service.FileStats} "统计信息"
// @Router       /api/v1/util/statistics [get]
func (h *ImageHandler) GetStatistics(ctx *gin.Context) {
	statService := service.NewStatisticsService()
	stats := statService.GetStatistics(ctx.Request.Context())

	utils.SuccessResponse(ctx, stats)
}

// GetDiskUsage returns disk usage info
// @Summary      Get disk usage
// @Description  获取磁盘使用情况：已用空间、限制、使用率
// @Tags         Utility
// @Produce      json
// @Success      200 {object} utils.Response{data=service.DiskUsage} "磁盘使用信息"
// @Router       /api/v1/util/disk-usage [get]
func (h *ImageHandler) GetDiskUsage(ctx *gin.Context) {
	statService := service.NewStatisticsService()
	usage := statService.GetDiskUsage(ctx.Request.Context())

	utils.SuccessResponse(ctx, usage)
}

// ExportFiles exports specified files as ZIP
// @Summary      Export files as ZIP
// @Description  将指定文件列表导出为 ZIP 压缩包
// @Tags         Utility
// @Accept       json
// @Produce      json
// @Param        body body object true "要导出的文件名列表" SchemaExample({"filenames":["image1.jpg","image2.png"]})
// @Success      200 {object} utils.Response{data=service.ExportResult} "导出结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/util/export [post]
func (h *ImageHandler) ExportFiles(ctx *gin.Context) {
	var req struct {
		Filenames []string `json:"filenames" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	exportService := service.NewExportService()
	result, err := exportService.ExportMultipleFiles(ctx.Request.Context(), req.Filenames, "./files")
	if err != nil {
		utils.CustomResponse(ctx, http.StatusInternalServerError, "export failed", nil)
		return
	}

	utils.SuccessResponse(ctx, result)
}

// ExportAllFiles exports all files as ZIP
// @Summary      Export all files as ZIP
// @Description  将所有图片文件导出为 ZIP 压缩包
// @Tags         Utility
// @Produce      json
// @Success      200 {object} utils.Response{data=service.ExportResult} "导出结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/util/export-all [post]
func (h *ImageHandler) ExportAllFiles(ctx *gin.Context) {
	exportService := service.NewExportService()
	result, err := exportService.ExportAllFiles(ctx.Request.Context(), "./files")
	if err != nil {
		utils.CustomResponse(ctx, http.StatusInternalServerError, "export failed", nil)
		return
	}

	utils.SuccessResponse(ctx, result)
}

// Cleanup performs disk cleanup
// @Summary      Cleanup files
// @Description  执行磁盘清理操作：删除孤立缩略图、过期文件、空目录
// @Tags         Utility
// @Accept       json
// @Produce      json
// @Param        body body service.CleanupConfig true "清理配置"
// @Success      200 {object} utils.Response{data=service.CleanupResult} "清理结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/util/cleanup [post]
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
		maxAge = 24 * time.Hour * 30
	}

	result := maintService.Cleanup(ctx.Request.Context(), service.CleanupConfig{
		RemoveOrphanThumbnails: req.RemoveOrphanThumbnails,
		RemoveOldFiles:         req.RemoveOldFiles,
		MaxFileAge:             maxAge,
		RemoveEmptyDirs:        req.RemoveEmptyDirs,
	})

	utils.SuccessResponse(ctx, result)
}

// StartThumbnailGeneration triggers thumbnail generation for existing images
// @Summary      Generate thumbnails
// @Description  为已有图片生成缩略图（后台异步任务）
// @Tags         Utility
// @Produce      json
// @Param        filenames query string true "要生成缩略图的文件名（逗号分隔）"
// @Success      200 {object} utils.Response "任务已触发"
// @Security     ApiKeyAuth
// @Router       /api/v1/util/generate-thumbnails [post]
func (h *ImageHandler) StartThumbnailGeneration(ctx *gin.Context) {
	filenames := ctx.Query("filenames")

	if filenames == "" {
		utils.CustomResponse(ctx, http.StatusBadRequest, "filenames parameter required", nil)
		return
	}

	fileList := strings.Split(filenames, ",")
	total := len(fileList)
	uploadDir := config.GetConfig().File.UploadDir

	go func() {
		success := 0
		failed := 0
		for _, name := range fileList {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			srcPath := filepath.Join(uploadDir, name)
			thumbDir := filepath.Join(uploadDir, "thumbs")
			dstPath := filepath.Join(thumbDir, name)

			if err := imageutil.GenerateThumbnail(srcPath, dstPath, imageutil.DefaultThumbnailConfig); err != nil {
				h.logger.Error("Thumbnail failed for %s: %v", name, err)
				failed++
			} else {
				success++
			}
		}
		h.logger.Info("Thumbnail generation complete: %d/%d success, %d failed", success, total, failed)
	}()

	utils.SuccessResponse(ctx, map[string]interface{}{
		"message":  "Thumbnail generation started (background task)",
		"total":    total,
	})
}

// InitChunkUpload initializes a chunk upload session
// @Summary      Init chunk upload
// @Description  初始化分片上传会话，返回 upload_id
// @Tags         Upload
// @Accept       json
// @Produce      json
// @Param        body body object true "分片上传初始化参数" SchemaExample({"filename":"large.jpg","file_size":10485760})
// @Success      200 {object} utils.Response{data=object} "upload_id"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/upload/chunk/init [post]
func (h *ImageHandler) InitChunkUpload(ctx *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		FileSize int64  `json:"file_size"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}
	uploadID, err := h.service.InitChunkUpload(ctx.Request.Context(), req.Filename, req.FileSize)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, map[string]interface{}{
		"upload_id": uploadID,
		"filename":  req.Filename,
	})
}

// UploadChunk uploads a chunk of a file
// @Summary      Upload chunk
// @Description  上传文件分片（二进制）
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        upload_id formData string true "上传会话 ID"
// @Param        chunk_index formData int true "分片序号（从0开始）"
// @Param        file formData file true "分片数据"
// @Success      200 {object} utils.Response{data=object} "分片上传结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/upload/chunk [post]
func (h *ImageHandler) UploadChunk(ctx *gin.Context) {
	uploadID := ctx.PostForm("upload_id")
	chunkIdxStr := ctx.PostForm("chunk_index")
	chunkIdx, err := strconv.Atoi(chunkIdxStr)
	if err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid chunk_index", nil)
		return
	}
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "file chunk required", nil)
		return
	}
	defer file.Close()

	if err := h.service.SaveChunk(ctx.Request.Context(), uploadID, chunkIdx, file); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.SuccessResponse(ctx, map[string]interface{}{
		"upload_id":   uploadID,
		"chunk_index": chunkIdx,
	})
}

// CompleteChunkUpload merges all chunks into the final file
// @Summary      Complete chunk upload
// @Description  合并所有分片为最终文件
// @Tags         Upload
// @Accept       json
// @Produce      json
// @Param        body body object true "完成分片上传" SchemaExample({"upload_id":"xxx"})
// @Success      200 {object} utils.Response{data=object} "合并结果"
// @Security     ApiKeyAuth
// @Router       /api/v1/images/upload/chunk/complete [post]
func (h *ImageHandler) CompleteChunkUpload(ctx *gin.Context) {
	var req struct {
		UploadID string `json:"upload_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}
	result, appErr := h.service.CompleteChunkUpload(ctx.Request.Context(), req.UploadID)
	if appErr != nil {
		utils.ErrorResponse(ctx, appErr)
		return
	}
	h.service.InvalidateCache()
	utils.SuccessResponse(ctx, result)
}

// ReloadConfig reloads server configuration from environment
// @Summary      Reload configuration
// @Description  重新加载服务器配置（从环境变量）
// @Tags         System
// @Produce      json
// @Success      200 {object} utils.Response "配置已重载"
// @Router       /api/v1/system/reload [post]
func (h *ImageHandler) ReloadConfig(ctx *gin.Context) {
	config.GetConfig().Reload()
	h.logger.Info("Configuration reloaded")
	utils.SuccessResponse(ctx, map[string]string{
		"message": "configuration reloaded",
	})
}

// Login handles user login and returns JWT token
// @Summary      User login
// @Description  用户登录，返回 JWT Token（开发账号：admin/admin123, user/user123）
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object true "登录凭证" SchemaExample({"username":"admin","password":"admin123"})
// @Success      200 {object} utils.Response{data=object} "登录成功，返回 Token"
// @Failure      401 {object} utils.Response "用户名或密码错误"
// @Router       /api/auth/login [post]
func (h *ImageHandler) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.CustomResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

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
// @Summary      Refresh JWT token
// @Description  刷新 JWT Token 的有效期
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=object} "刷新成功，返回新 Token"
// @Failure      401 {object} utils.Response "Token 无效"
// @Router       /api/auth/refresh [post]
func (h *ImageHandler) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		utils.CustomResponse(ctx, http.StatusUnauthorized, "missing authorization header", nil)
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

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
