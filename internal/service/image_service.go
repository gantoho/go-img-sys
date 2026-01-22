package service

import (
	"math/rand"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/cache"
	"github.com/gantoho/go-img-sys/pkg/errors"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
)

type ImageData struct {
	Total int      `json:"total"`
	Data  []string `json:"data"`
}

type ImageMetaData struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	SizeStr  string `json:"size_str"`
	MimeType string `json:"mime_type"`
	ModTime  int64  `json:"mod_time"`
}

type PaginatedImageData struct {
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Pages    int             `json:"pages"`
	Data     []ImageMetaData `json:"data"`
}

type ImageService struct {
	config *config.Config
	logger *logger.Logger
	cache  *cache.Cache
}

func NewImageService() *ImageService {
	return &ImageService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
		cache:  cache.NewCache(),
	}
}

// GetImageByFilename retrieves a single image by filename
func (s *ImageService) GetImageByFilename(filename string) (string, *errors.AppError) {
	filepath := utils.GetUploadPath(s.config.File.UploadDir, filename)

	if !utils.FileExists(filepath) {
		s.logger.Warn("File not found: %s", filename)
		return "", errors.ErrFileNotFound
	}

	return filepath, nil
}

// GetAllImages returns all image files with their URLs
func (s *ImageService) GetAllImages(hostURL string) (*ImageData, *errors.AppError) {
	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to list images", err)
	}

	if len(fileInfos) == 0 {
		return &ImageData{Total: 0, Data: []string{}}, nil
	}

	data := &ImageData{
		Total: len(fileInfos),
		Data:  make([]string, 0, len(fileInfos)),
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			data.Data = append(data.Data, hostURL+"/f/"+fileInfo.Name())
		}
	}

	return data, nil
}

// GetAllImagesWithMetadata returns all image files with detailed metadata
func (s *ImageService) GetAllImagesWithMetadata(hostURL string) ([]ImageMetaData, *errors.AppError) {
	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to list images", err)
	}

	result := make([]ImageMetaData, 0)

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			filePath := utils.GetUploadPath(s.config.File.UploadDir, fileInfo.Name())
			size, _ := utils.GetFileInfo(filePath)
			metadata := ImageMetaData{
				Filename: fileInfo.Name(),
				URL:      hostURL + "/f/" + fileInfo.Name(),
				Size:     size,
				SizeStr:  utils.GetFileSizeFormatted(size),
				MimeType: utils.GetMimeType(fileInfo.Name()),
				ModTime:  fileInfo.ModTime().Unix(),
			}
			result = append(result, metadata)
		}
	}

	return result, nil
}

// GetAllImagesPaginated returns paginated image files with metadata
func (s *ImageService) GetAllImagesPaginated(hostURL string, page, pageSize int) (*PaginatedImageData, *errors.AppError) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Try to get from cache
	cacheKey := "images_list"
	if cached, ok := s.cache.Get(cacheKey); ok {
		fileInfos := cached.([]os.FileInfo)
		return s.paginateFiles(fileInfos, hostURL, page, pageSize), nil
	}

	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to list images", err)
	}

	// Cache for 5 minutes
	s.cache.Set(cacheKey, fileInfos, 5*time.Minute)

	return s.paginateFiles(fileInfos, hostURL, page, pageSize), nil
}

// Helper function to paginate file list
func (s *ImageService) paginateFiles(fileInfos []os.FileInfo, hostURL string, page, pageSize int) *PaginatedImageData {
	validFiles := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			validFiles = append(validFiles, fileInfo)
		}
	}

	total := len(validFiles)
	pages := (total + pageSize - 1) / pageSize

	if page > pages && total > 0 {
		page = pages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	result := &PaginatedImageData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
		Data:     make([]ImageMetaData, 0, pageSize),
	}

	if start < total {
		for i := start; i < end; i++ {
			file := validFiles[i]
			filePath := utils.GetUploadPath(s.config.File.UploadDir, file.Name())
			size, _ := utils.GetFileInfo(filePath)
			metadata := ImageMetaData{
				Filename: file.Name(),
				URL:      hostURL + "/f/" + file.Name(),
				Size:     size,
				SizeStr:  utils.GetFileSizeFormatted(size),
				MimeType: utils.GetMimeType(file.Name()),
				ModTime:  file.ModTime().Unix(),
			}
			result.Data = append(result.Data, metadata)
		}
	}

	return result
}

// GetRandomImage returns a random image filename
func (s *ImageService) GetRandomImage() (string, *errors.AppError) {
	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return "", errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to get random image", err)
	}

	if len(fileInfos) == 0 {
		return "", errors.ErrNoFiles
	}

	validFiles := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			validFiles = append(validFiles, fileInfo)
		}
	}

	if len(validFiles) == 0 {
		return "", errors.ErrNoFiles
	}

	randomIndex := rand.Intn(len(validFiles))
	return validFiles[randomIndex].Name(), nil
}

// GetRandomImages returns multiple random images
func (s *ImageService) GetRandomImages(hostURL string, count int) ([]string, *errors.AppError) {
	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to get random images", err)
	}

	validFiles := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			validFiles = append(validFiles, fileInfo)
		}
	}

	if len(validFiles) == 0 {
		return nil, errors.ErrNoFiles
	}

	if count > len(validFiles) {
		count = len(validFiles)
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		randomIndex := rand.Intn(len(validFiles))
		result = append(result, hostURL+"/f/"+validFiles[randomIndex].Name())
	}

	return result, nil
}

// UploadFile uploads a file to the server
func (s *ImageService) UploadFile(hostURL string, files []*multipart.FileHeader) (map[string]interface{}, *errors.AppError) {
	if len(files) == 0 {
		return nil, errors.NewError(400, "no files provided")
	}

	// Ensure upload directory exists
	if err := utils.EnsureDir(s.config.File.UploadDir); err != nil {
		s.logger.Error("Failed to create upload directory: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrInternalServer.Code, "Failed to create directory", err)
	}

	uploadedFiles := make([]string, 0)
	failedFiles := make([]map[string]string, 0)

	for _, file := range files {
		// Validate file
		if appErr := s.validateFile(file); appErr != nil {
			s.logger.Warn("File validation failed: %s, error: %s", file.Filename, appErr.Message)
			failedFiles = append(failedFiles, map[string]string{
				"filename": file.Filename,
				"error":    appErr.Message,
			})
			continue
		}

		// Save file
		// Note: In production, you should handle file name conflicts
		uploadedFiles = append(uploadedFiles, hostURL+"/f/"+file.Filename)

		// Save file using Gin's method would require gin.Context
		// For now, we'll prepare the path but need gin.Context in handler
	}

	result := map[string]interface{}{
		"total_uploaded": len(uploadedFiles),
		"uploaded":       uploadedFiles,
	}

	if len(failedFiles) > 0 {
		result["failed"] = failedFiles
		result["total_failed"] = len(failedFiles)
	}

	return result, nil
}

// SearchImages filters images by criteria
func (s *ImageService) SearchImages(hostURL string, filename string, minSize, maxSize int64, fileType string, page, pageSize int) (*PaginatedImageData, *errors.AppError) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	fileInfos, err := utils.ListFiles(s.config.File.UploadDir)
	if err != nil {
		s.logger.Error("Failed to list files: %v", err)
		return nil, errors.NewErrorWithCause(errors.ErrDirectoryFail.Code, "Failed to search images", err)
	}

	validFiles := make([]os.FileInfo, 0)

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		// Filter by filename
		if filename != "" && !strings.Contains(strings.ToLower(fileInfo.Name()), strings.ToLower(filename)) {
			continue
		}

		// Filter by file size
		if minSize > 0 && fileInfo.Size() < minSize {
			continue
		}
		if maxSize > 0 && fileInfo.Size() > maxSize {
			continue
		}

		// Filter by file type/extension
		if fileType != "" {
			ext := strings.ToLower(utils.GetFileExt(fileInfo.Name()))
			if !strings.EqualFold(ext, "."+fileType) && !strings.EqualFold(ext, fileType) {
				continue
			}
		}

		validFiles = append(validFiles, fileInfo)
	}

	total := len(validFiles)
	pages := (total + pageSize - 1) / pageSize

	if page > pages && total > 0 {
		page = pages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	result := &PaginatedImageData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
		Data:     make([]ImageMetaData, 0, pageSize),
	}

	if start < total {
		for i := start; i < end; i++ {
			file := validFiles[i]
			filePath := utils.GetUploadPath(s.config.File.UploadDir, file.Name())
			size, _ := utils.GetFileInfo(filePath)
			metadata := ImageMetaData{
				Filename: file.Name(),
				URL:      hostURL + "/f/" + file.Name(),
				Size:     size,
				SizeStr:  utils.GetFileSizeFormatted(size),
				MimeType: utils.GetMimeType(file.Name()),
				ModTime:  file.ModTime().Unix(),
			}
			result.Data = append(result.Data, metadata)
		}
	}

	return result, nil
}
func (s *ImageService) validateFile(file *multipart.FileHeader) *errors.AppError {
	// Check file size
	if file.Size > s.config.File.MaxSize*1024*1024 {
		return errors.ErrFileTooLarge
	}

	// Check if file format is supported
	if !utils.IsValidImageFormat(file.Filename) {
		return errors.NewError(400, "unsupported image format. Supported formats: jpg, jpeg, png, gif, webp, bmp, ico, svg")
	}

	return nil
}

// DeleteImage deletes a single image file
func (s *ImageService) DeleteImage(filename string) *errors.AppError {
	if !utils.IsValidImageFormat(filename) {
		return errors.NewError(400, "invalid filename format")
	}

	filepath := utils.GetUploadPath(s.config.File.UploadDir, filename)

	if !utils.FileExists(filepath) {
		return errors.ErrFileNotFound
	}

	if err := os.Remove(filepath); err != nil {
		s.logger.Error("Failed to delete file %s: %v", filename, err)
		return errors.NewErrorWithCause(500, "failed to delete file", err)
	}

	// Clear cache after deletion
	s.cache.Delete("images_list")
	s.logger.Info("File deleted: %s", filename)

	return nil
}

// DeleteImages deletes multiple image files
func (s *ImageService) DeleteImages(filenames []string) map[string]interface{} {
	deleted := make([]string, 0)
	failed := make([]map[string]string, 0)

	for _, filename := range filenames {
		if err := s.DeleteImage(filename); err != nil {
			failed = append(failed, map[string]string{
				"filename": filename,
				"error":    err.Message,
			})
		} else {
			deleted = append(deleted, filename)
		}
	}

	result := map[string]interface{}{
		"total_deleted": len(deleted),
		"deleted":       deleted,
	}

	if len(failed) > 0 {
		result["total_failed"] = len(failed)
		result["failed"] = failed
	}

	return result
}
