package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/logger"
	"github.com/gantoho/go-img-sys/pkg/utils"
)

// StatisticsService 统计服务
type StatisticsService struct {
	config *config.Config
	logger *logger.Logger
}

// NewStatisticsService 创建统计服务
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
	}
}

// FileStats 文件统计信息
type FileStats struct {
	TotalFiles      int                   `json:"total_files"`
	TotalSize       int64                 `json:"total_size"`
	TotalSizeStr    string                `json:"total_size_str"`
	AverageFileSize int64                 `json:"average_file_size"`
	FormatStats     map[string]FormatStat `json:"format_stats"`
	LargestFile     string                `json:"largest_file"`
	LargestFileSize int64                 `json:"largest_file_size"`
}

// FormatStat 格式统计
type FormatStat struct {
	Count      int     `json:"count"`
	Size       int64   `json:"size"`
	SizeStr    string  `json:"size_str"`
	Percentage float64 `json:"percentage"`
}

// GetStatistics 获取统计信息
func (s *StatisticsService) GetStatistics() *FileStats {
	stats := &FileStats{
		FormatStats: make(map[string]FormatStat),
	}

	uploadDir := s.config.File.UploadDir

	var largestSize int64

	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 忽略缩略图
		if strings.Contains(path, "thumbs") {
			return nil
		}

		size := info.Size()
		stats.TotalFiles++
		stats.TotalSize += size

		// 跟踪最大文件
		if size > largestSize {
			largestSize = size
			stats.LargestFile = filepath.Base(path)
			stats.LargestFileSize = size
		}

		// 统计格式
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != "" {
			formatStat := stats.FormatStats[ext]
			formatStat.Count++
			formatStat.Size += size
			stats.FormatStats[ext] = formatStat
		}

		return nil
	})

	// 计算平均大小
	if stats.TotalFiles > 0 {
		stats.AverageFileSize = stats.TotalSize / int64(stats.TotalFiles)
	}

	// 格式化大小字符串
	stats.TotalSizeStr = utils.GetFileSizeFormatted(stats.TotalSize)

	// 计算百分比
	for fmt, stat := range stats.FormatStats {
		stat.SizeStr = utils.GetFileSizeFormatted(stat.Size)
		if stats.TotalSize > 0 {
			stat.Percentage = float64(stat.Size) / float64(stats.TotalSize) * 100
		}
		stats.FormatStats[fmt] = stat
	}

	s.logger.Info("Statistics computed: %d files, %.2f MB total", stats.TotalFiles, float64(stats.TotalSize)/1024/1024)

	return stats
}

// DiskUsage 磁盘使用情况
type DiskUsage struct {
	UsedSpace    int64   `json:"used_space"`
	UsedSpaceStr string  `json:"used_space_str"`
	Limit        int64   `json:"limit"` // 设置的限制
	LimitStr     string  `json:"limit_str"`
	Percentage   float64 `json:"percentage"`
}

// GetDiskUsage 获取磁盘使用情况
func (s *StatisticsService) GetDiskUsage() *DiskUsage {
	stats := s.GetStatistics()
	maxSize := s.config.File.MaxSize * 1024 * 1024 // 转换为字节

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
