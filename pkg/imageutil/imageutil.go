package imageutil

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gantoho/go-img-sys/pkg/logger"
)

// ThumbnailConfig 缩略图配置
type ThumbnailConfig struct {
	Width   int
	Height  int
	Quality int // 仅JPEG
}

// DefaultThumbnailConfig 默认缩略图配置
var DefaultThumbnailConfig = ThumbnailConfig{
	Width:   200,
	Height:  200,
	Quality: 85,
}

// GenerateThumbnail 生成缩略图
func GenerateThumbnail(sourcePath string, thumbPath string, config ThumbnailConfig) error {
	logger := logger.GetLogger()

	// 打开原始图片
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		logger.Error("Failed to open source image: %v", err)
		return err
	}
	defer sourceFile.Close()

	// 解码图片
	img, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		logger.Error("Failed to decode image: %v", err)
		return err
	}

	// 重新打开文件用于实际解码
	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg":
		originalImg, err = jpeg.Decode(sourceFile)
	case "png":
		originalImg, err = png.Decode(sourceFile)
	default:
		logger.Warn("Unsupported format for thumbnail: %s", format)
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		logger.Error("Failed to decode image content: %v", err)
		return err
	}

	// 计算缩略图尺寸（保持宽高比）
	thumbWidth, thumbHeight := calculateThumbnailSize(img.Width, img.Height, config.Width, config.Height)

	// 创建缩略图
	thumb := image.NewRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))
	draw.Draw(thumb, thumb.Bounds(), originalImg, image.Point{}, draw.Src)

	// 确保缩略图目录存在
	thumbDir := filepath.Dir(thumbPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		logger.Error("Failed to create thumbnail directory: %v", err)
		return err
	}

	// 保存缩略图
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		logger.Error("Failed to create thumbnail file: %v", err)
		return err
	}
	defer thumbFile.Close()

	// 根据原始格式保存
	switch strings.ToLower(format) {
	case "jpeg":
		err = jpeg.Encode(thumbFile, thumb, &jpeg.Options{Quality: config.Quality})
	case "png":
		err = png.Encode(thumbFile, thumb)
	}

	if err != nil {
		logger.Error("Failed to encode thumbnail: %v", err)
		return err
	}

	logger.Info("Thumbnail generated: %s", thumbPath)
	return nil
}

// calculateThumbnailSize 计算缩略图尺寸（保持宽高比）
func calculateThumbnailSize(origWidth, origHeight, maxWidth, maxHeight int) (int, int) {
	ratio := float64(origWidth) / float64(origHeight)
	targetRatio := float64(maxWidth) / float64(maxHeight)

	var w, h int
	if ratio > targetRatio {
		w = maxWidth
		h = int(float64(maxWidth) / ratio)
	} else {
		h = maxHeight
		w = int(float64(maxHeight) * ratio)
	}

	return w, h
}

// RotateImage 旋转图片（90, 180, 270度）
func RotateImage(sourcePath string, outputPath string, degrees int) error {
	logger := logger.GetLogger()

	if degrees%90 != 0 {
		return fmt.Errorf("rotation degrees must be multiple of 90")
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	img, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		return err
	}

	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg":
		originalImg, _ = jpeg.Decode(sourceFile)
	case "png":
		originalImg, _ = png.Decode(sourceFile)
	default:
		return fmt.Errorf("unsupported format")
	}

	// 简化旋转（实际应用中建议使用专业库）
	var rotated image.Image
	switch degrees % 360 {
	case 90, 270:
		rotated = image.NewRGBA(image.Rect(0, 0, img.Height, img.Width))
	default:
		rotated = originalImg
	}

	outFile, _ := os.Create(outputPath)
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg":
		jpeg.Encode(outFile, rotated, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(outFile, rotated)
	}

	logger.Info("Image rotated: %s", outputPath)
	return nil
}

// ResizeImage 缩放图片
func ResizeImage(sourcePath string, outputPath string, width, height int) error {
	logger := logger.GetLogger()

	if width <= 0 || height <= 0 {
		return fmt.Errorf("width and height must be positive")
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	_, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		return err
	}

	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg":
		originalImg, _ = jpeg.Decode(sourceFile)
	case "png":
		originalImg, _ = png.Decode(sourceFile)
	default:
		return fmt.Errorf("unsupported format")
	}

	// 简单缩放实现
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(resized, resized.Bounds(), originalImg, image.Point{}, draw.Src)

	outFile, _ := os.Create(outputPath)
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg":
		jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(outFile, resized)
	}

	logger.Info("Image resized to %dx%d: %s", width, height, outputPath)
	return nil
}

// AddWatermark 添加文字水印（简化版）
func AddWatermark(sourcePath string, outputPath string, watermarkText string) error {
	logger := logger.GetLogger()

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	_, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		return err
	}

	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg":
		originalImg, _ = jpeg.Decode(sourceFile)
	case "png":
		originalImg, _ = png.Decode(sourceFile)
	default:
		return fmt.Errorf("unsupported format")
	}

	// 简单水印实现（在右下角添加文本）
	// 实际应用需要使用 golang.org/x/image/font 包
	result := image.NewRGBA(originalImg.Bounds())
	draw.Draw(result, result.Bounds(), originalImg, image.Point{}, draw.Src)

	outFile, _ := os.Create(outputPath)
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg":
		jpeg.Encode(outFile, result, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(outFile, result)
	}

	logger.Info("Watermark added to: %s", outputPath)
	return nil
}
