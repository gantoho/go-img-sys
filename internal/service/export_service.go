package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/logger"
)

// ExportService 导出服务
type ExportService struct {
	config *config.Config
	logger *logger.Logger
}

// NewExportService 创建导出服务
func NewExportService() *ExportService {
	return &ExportService{
		config: config.GetConfig(),
		logger: logger.GetLogger(),
	}
}

// ExportResult 导出结果
type ExportResult struct {
	ZipPath    string `json:"zip_path"`
	FileCount  int    `json:"file_count"`
	TotalSize  int64  `json:"total_size"`
	SizeStr    string `json:"size_str"`
	Compressed bool   `json:"compressed"`
}

// ExportMultipleFiles 导出多个文件为ZIP
func (e *ExportService) ExportMultipleFiles(filenames []string, outputDir string) (*ExportResult, error) {
	if len(filenames) == 0 {
		e.logger.Warn("No files to export")
		return nil, fmt.Errorf("no files provided")
	}

	// 生成ZIP文件名
	zipName := "export_" + getCurrentTimestamp() + ".zip"
	zipPath := filepath.Join(outputDir, zipName)

	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		e.logger.Error("Failed to create zip file: %v", err)
		return nil, err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	result := &ExportResult{
		ZipPath: zipPath,
	}

	uploadDir := e.config.File.UploadDir

	for _, filename := range filenames {
		filePath := filepath.Join(uploadDir, filename)

		// 安全检查：防止路径遍历
		absPath, _ := filepath.Abs(filePath)
		absUploadDir, _ := filepath.Abs(uploadDir)
		if !strings.HasPrefix(absPath, absUploadDir) {
			e.logger.Warn("Invalid file path: %s", filePath)
			continue
		}

		// 检查文件是否存在
		fileInfo, err := os.Stat(filePath)
		if err != nil || fileInfo.IsDir() {
			e.logger.Warn("File not found or is directory: %s", filePath)
			continue
		}

		// 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			e.logger.Error("Failed to open file %s: %v", filePath, err)
			continue
		}

		// 添加到ZIP
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			file.Close()
			continue
		}

		header.Name = filename
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			file.Close()
			continue
		}

		_, err = io.Copy(writer, file)
		file.Close()

		if err == nil {
			result.FileCount++
			result.TotalSize += fileInfo.Size()
			e.logger.Info("File added to zip: %s", filename)
		}
	}

	result.SizeStr = getFileSizeStr(result.TotalSize)
	result.Compressed = true

	e.logger.Info("Export completed: %s (files: %d, size: %s)", zipPath, result.FileCount, result.SizeStr)

	return result, nil
}

// ExportAllFiles 导出所有文件
func (e *ExportService) ExportAllFiles(outputDir string) (*ExportResult, error) {
	uploadDir := e.config.File.UploadDir

	// 获取所有文件
	var filenames []string
	filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// 忽略缩略图
		if strings.Contains(path, "thumbs") {
			return nil
		}

		rel, _ := filepath.Rel(uploadDir, path)
		filenames = append(filenames, rel)
		return nil
	})

	return e.ExportMultipleFiles(filenames, outputDir)
}

// getCurrentTimestamp 获取当前时间戳格式字符串
func getCurrentTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// getFileSizeStr 获取格式化的文件大小
func getFileSizeStr(size int64) string {
	const (
		Byte     = 1
		KiloByte = Byte * 1024
		MegaByte = KiloByte * 1024
		GigaByte = MegaByte * 1024
	)

	switch {
	case size >= GigaByte:
		return formatSize(float64(size), "GB", GigaByte)
	case size >= MegaByte:
		return formatSize(float64(size), "MB", MegaByte)
	case size >= KiloByte:
		return formatSize(float64(size), "KB", KiloByte)
	default:
		return formatSize(float64(size), "B", Byte)
	}
}

func formatSize(size float64, unit string, divisor int64) string {
	return fmt.Sprintf("%.2f %s", size/float64(divisor), unit)
}
