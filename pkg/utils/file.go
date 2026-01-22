package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Supported image formats
var SupportedImageFormats = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".bmp":  "image/bmp",
	".ico":  "image/x-icon",
	".svg":  "image/svg+xml",
}

// EnsureDir creates directory if not exists
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// ListFiles returns all files in a directory
func ListFiles(dir string) ([]os.FileInfo, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file.Readdir(-1)
}

// FileExists checks if file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetFileExt returns file extension
func GetFileExt(filename string) string {
	return filepath.Ext(filename)
}

// JoinPath safely joins paths
func JoinPath(dir, file string) string {
	return filepath.Join(dir, file)
}

// GetUploadPath returns the full upload path
func GetUploadPath(uploadDir, filename string) string {
	return filepath.Join(uploadDir, filename)
}

// IsValidImageFormat validates if file is a supported image format
func IsValidImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, exists := SupportedImageFormats[ext]
	return exists
}

// GetMimeType returns MIME type for image format
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mime, exists := SupportedImageFormats[ext]; exists {
		return mime
	}
	return "application/octet-stream"
}

// GetFileInfo returns file size in bytes and name
func GetFileInfo(filepath string) (int64, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileSizeFormatted returns formatted file size string
func GetFileSizeFormatted(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
