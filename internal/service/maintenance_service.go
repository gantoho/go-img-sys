package service

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
)

type MaintenanceService struct {
	config *config.Config
	logger *logger.Logger
}

func NewMaintenanceService() *MaintenanceService {
	return &MaintenanceService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
	}
}

type CleanupConfig struct {
	RemoveOrphanThumbnails bool
	RemoveOldFiles         bool
	MaxFileAge             time.Duration
	RemoveEmptyDirs        bool
}

type CleanupResult struct {
	FilesRemoved      int
	ThumbnailsRemoved int
	DirsRemoved       int
	SizeFreed         int64
	Errors            []string
}

func (m *MaintenanceService) Cleanup(ctx context.Context, cfg CleanupConfig) *CleanupResult {
	result := &CleanupResult{
		Errors: make([]string, 0),
	}

	uploadDir := m.config.File.UploadDir

	if cfg.RemoveOrphanThumbnails {
		if err := ctx.Err(); err != nil {
			return result
		}
		m.cleanupOrphanThumbnails(uploadDir, result)
	}

	if cfg.RemoveOldFiles {
		if err := ctx.Err(); err != nil {
			return result
		}
		m.cleanupOldFiles(uploadDir, cfg.MaxFileAge, result)
	}

	if cfg.RemoveEmptyDirs {
		if err := ctx.Err(); err != nil {
			return result
		}
		m.cleanupEmptyDirs(uploadDir, result)
	}

	m.logger.Info("Cleanup completed: %d files removed, %d thumbnails removed, %d dirs removed, %.2f MB freed",
		result.FilesRemoved, result.ThumbnailsRemoved, result.DirsRemoved, float64(result.SizeFreed)/1024/1024)

	return result
}

func (m *MaintenanceService) cleanupOrphanThumbnails(uploadDir string, result *CleanupResult) {
	thumbDir := filepath.Join(uploadDir, "thumbs")
	if _, err := os.Stat(thumbDir); os.IsNotExist(err) {
		return
	}

	filepath.Walk(thumbDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

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

func (m *MaintenanceService) StartAutoCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				m.logger.Info("Auto cleanup stopped")
				return
			case <-ticker.C:
				cfg := CleanupConfig{
					RemoveOrphanThumbnails: true,
					RemoveOldFiles:         true,
					MaxFileAge:             24 * time.Hour * 30,
					RemoveEmptyDirs:        true,
				}
				m.Cleanup(ctx, cfg)
			}
		}
	}()

	m.logger.Info("Auto cleanup started with interval: %v", interval)
}
