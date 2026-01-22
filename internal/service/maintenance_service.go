package service

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
)

// MaintenanceService 维护服务
type MaintenanceService struct {
	config *config.Config
	logger *logger.Logger
}

// NewMaintenanceService 创建维护服务
func NewMaintenanceService() *MaintenanceService {
	return &MaintenanceService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
	}
}

// CleanupConfig 清理配置
type CleanupConfig struct {
	RemoveOrphanThumbnails bool          // 删除孤立缩略图
	RemoveOldFiles         bool          // 删除旧文件
	MaxFileAge             time.Duration // 最大文件年龄
	RemoveEmptyDirs        bool          // 删除空目录
}

// CleanupResult 清理结果
type CleanupResult struct {
	FilesRemoved      int
	ThumbnailsRemoved int
	DirsRemoved       int
	SizeFreed         int64
	Errors            []string
}

// Cleanup 执行清理操作
func (m *MaintenanceService) Cleanup(cfg CleanupConfig) *CleanupResult {
	result := &CleanupResult{
		Errors: make([]string, 0),
	}

	uploadDir := m.config.File.UploadDir

	if cfg.RemoveOrphanThumbnails {
		m.cleanupOrphanThumbnails(uploadDir, result)
	}

	if cfg.RemoveOldFiles {
		m.cleanupOldFiles(uploadDir, cfg.MaxFileAge, result)
	}

	if cfg.RemoveEmptyDirs {
		m.cleanupEmptyDirs(uploadDir, result)
	}

	m.logger.Info("Cleanup completed: %d files removed, %d thumbnails removed, %d dirs removed, %.2f MB freed",
		result.FilesRemoved, result.ThumbnailsRemoved, result.DirsRemoved, float64(result.SizeFreed)/1024/1024)

	return result
}

// cleanupOrphanThumbnails 清理孤立缩略图
func (m *MaintenanceService) cleanupOrphanThumbnails(uploadDir string, result *CleanupResult) {
	thumbDir := filepath.Join(uploadDir, "thumbs")
	if _, err := os.Stat(thumbDir); os.IsNotExist(err) {
		return
	}

	filepath.Walk(thumbDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 检查对应的原始文件是否存在
		relPath, _ := filepath.Rel(thumbDir, path)
		originalPath := filepath.Join(uploadDir, relPath)

		if !utils.FileExists(originalPath) {
			size := info.Size()
			os.Remove(path)
			result.ThumbnailsRemoved++
			result.SizeFreed += size
			m.logger.Info("Orphan thumbnail removed: %s", path)
		}

		return nil
	})
}

// cleanupOldFiles 清理旧文件
func (m *MaintenanceService) cleanupOldFiles(uploadDir string, maxAge time.Duration, result *CleanupResult) {
	cutoffTime := time.Now().Add(-maxAge)

	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			size := info.Size()
			os.Remove(path)
			result.FilesRemoved++
			result.SizeFreed += size
			m.logger.Info("Old file removed: %s", path)
		}

		return nil
	})
}

// cleanupEmptyDirs 清理空目录
func (m *MaintenanceService) cleanupEmptyDirs(uploadDir string, result *CleanupResult) {
	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}

		if path == uploadDir {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err == nil && len(entries) == 0 {
			os.Remove(path)
			result.DirsRemoved++
			m.logger.Info("Empty directory removed: %s", path)
		}

		return nil
	})
}

// StartAutoCleanup 启动自动清理（后台定时任务）
func (m *MaintenanceService) StartAutoCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			cfg := CleanupConfig{
				RemoveOrphanThumbnails: true,
				RemoveOldFiles:         true,
				MaxFileAge:             24 * time.Hour * 30, // 30天
				RemoveEmptyDirs:        true,
			}
			m.Cleanup(cfg)
		}
	}()

	m.logger.Info("Auto cleanup started with interval: %v", interval)
}
