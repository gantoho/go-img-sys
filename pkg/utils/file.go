package utils

import (
	"fmt"
	"net/http"
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

var externalBaseURL string

// SetExternalBaseURL sets a global external base URL that overrides
// auto-detection in GetRequestBaseURL. Useful when running behind a reverse proxy.
func SetExternalBaseURL(url string) {
	externalBaseURL = url
}

// GetRequestBaseURL returns the full base URL (scheme + host) from an HTTP request.
// It checks X-Forwarded-Proto and X-Forwarded-Host headers for reverse proxy support.
// If an external base URL has been set via SetExternalBaseURL, it takes precedence.
func GetRequestBaseURL(r *http.Request) string {
	if externalBaseURL != "" {
		return externalBaseURL
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwd := r.Header.Get("X-Forwarded-Proto"); fwd == "https" || fwd == "http" {
		scheme = fwd
	}

	host := r.Host
	if fwdHost := r.Header.Get("X-Forwarded-Host"); fwdHost != "" {
		host = fwdHost
	}

	return scheme + "://" + host
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
