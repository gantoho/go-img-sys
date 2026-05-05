package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
)

type StatisticsService struct {
	config *config.Config
	logger *logger.Logger
}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
	}
}

type FileStats struct {
	TotalFiles      int                   `json:"total_files"`
	TotalSize       int64                 `json:"total_size"`
	TotalSizeStr    string                `json:"total_size_str"`
	AverageFileSize int64                 `json:"average_file_size"`
	FormatStats     map[string]FormatStat `json:"format_stats"`
	LargestFile     string                `json:"largest_file"`
	LargestFileSize int64                 `json:"largest_file_size"`
}

type FormatStat struct {
	Count      int     `json:"count"`
	Size       int64   `json:"size"`
	SizeStr    string  `json:"size_str"`
	Percentage float64 `json:"percentage"`
}

func (s *StatisticsService) GetStatistics(ctx context.Context) *FileStats {
	stats := &FileStats{
		FormatStats: make(map[string]FormatStat),
	}

	uploadDir := s.config.File.UploadDir
	var largestSize int64

	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "thumbs") {
			return nil
		}

		size := info.Size()
		stats.TotalFiles++
		stats.TotalSize += size

		if size > largestSize {
			largestSize = size
			stats.LargestFile = filepath.Base(path)
			stats.LargestFileSize = size
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != "" {
			formatStat := stats.FormatStats[ext]
			formatStat.Count++
			formatStat.Size += size
			stats.FormatStats[ext] = formatStat
		}
		return nil
	})

	if stats.TotalFiles > 0 {
		stats.AverageFileSize = stats.TotalSize / int64(stats.TotalFiles)
	}

	stats.TotalSizeStr = utils.GetFileSizeFormatted(stats.TotalSize)

	for fileFmt, stat := range stats.FormatStats {
		stat.SizeStr = utils.GetFileSizeFormatted(stat.Size)
		if stats.TotalSize > 0 {
			stat.Percentage = float64(stat.Size) / float64(stats.TotalSize) * 100
		}
		stats.FormatStats[fileFmt] = stat
	}

	s.logger.Info("Statistics computed: %d files, %.2f MB total", stats.TotalFiles, float64(stats.TotalSize)/1024/1024)

	return stats
}

type DiskUsage struct {
	UsedSpace    int64   `json:"used_space"`
	UsedSpaceStr string  `json:"used_space_str"`
	Limit        int64   `json:"limit"`
	LimitStr     string  `json:"limit_str"`
	Percentage   float64 `json:"percentage"`
}

func (s *StatisticsService) GetDiskUsage(ctx context.Context) *DiskUsage {
	stats := s.GetStatistics(ctx)
	maxSize := s.config.File.MaxSize * 1024 * 1024

	usage := &DiskUsage{
		UsedSpace:    stats.TotalSize,
		UsedSpaceStr: stats.TotalSizeStr,
		Limit:        maxSize,
		LimitStr:     utils.GetFileSizeFormatted(maxSize),
	}

	if maxSize > 0 {
		usage.Percentage = float64(stats.TotalSize) / float64(maxSize) * 100
	}

	return usage
}
