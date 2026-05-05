package imageutil

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gantoho/go-img-sys/pkg/logger"
	"golang.org/x/image/draw"
)

type ThumbnailConfig struct {
	Width   int
	Height  int
	Quality int
}

var DefaultThumbnailConfig = ThumbnailConfig{
	Width:   200,
	Height:  200,
	Quality: 85,
}

func GenerateThumbnail(sourcePath string, thumbPath string, config ThumbnailConfig) error {
	log := logger.GetLogger()

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Error("Failed to open source image: %v", err)
		return err
	}
	defer sourceFile.Close()

	imgCfg, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		log.Error("Failed to decode image config: %v", err)
		return err
	}

	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		originalImg, err = jpeg.Decode(sourceFile)
	case "png":
		originalImg, err = png.Decode(sourceFile)
	default:
		log.Warn("Unsupported format for thumbnail: %s", format)
		return fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		log.Error("Failed to decode image content: %v", err)
		return err
	}

	thumbWidth, thumbHeight := calculateThumbnailSize(imgCfg.Width, imgCfg.Height, config.Width, config.Height)

	thumb := image.NewRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))

	draw.CatmullRom.Scale(thumb, thumb.Bounds(), originalImg, originalImg.Bounds(), draw.Over, nil)

	thumbDir := filepath.Dir(thumbPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		log.Error("Failed to create thumbnail directory: %v", err)
		return err
	}

	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		log.Error("Failed to create thumbnail file: %v", err)
		return err
	}
	defer thumbFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(thumbFile, thumb, &jpeg.Options{Quality: config.Quality})
	case "png":
		err = png.Encode(thumbFile, thumb)
	}
	if err != nil {
		log.Error("Failed to encode thumbnail: %v", err)
		return err
	}

	log.Info("Thumbnail generated: %s", thumbPath)
	return nil
}

func calculateThumbnailSize(origWidth, origHeight, maxWidth, maxHeight int) (int, int) {
	if origWidth <= 0 || origHeight <= 0 {
		return maxWidth, maxHeight
	}
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
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return w, h
}

func rotatePixels90(src image.Image, dst *image.RGBA, degrees int) {
	bounds := src.Bounds()
	srcW := bounds.Max.X - bounds.Min.X
	srcH := bounds.Max.Y - bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			switch degrees % 360 {
			case 90:
				dst.Set(srcH-1-(y-bounds.Min.Y), x-bounds.Min.X, src.At(x, y))
			case 180:
				dst.Set(srcW-1-(x-bounds.Min.X), srcH-1-(y-bounds.Min.Y), src.At(x, y))
			case 270:
				dst.Set(y-bounds.Min.Y, srcW-1-(x-bounds.Min.X), src.At(x, y))
			}
		}
	}
}

func RotateImage(sourcePath string, outputPath string, degrees int) error {
	log := logger.GetLogger()

	if degrees%90 != 0 {
		return fmt.Errorf("rotation degrees must be multiple of 90")
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	imgCfg, format, err := image.DecodeConfig(sourceFile)
	if err != nil {
		return err
	}

	sourceFile.Seek(0, 0)

	var originalImg image.Image
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		originalImg, err = jpeg.Decode(sourceFile)
	case "png":
		originalImg, err = png.Decode(sourceFile)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return err
	}

	var rotated image.RGBA
	switch degrees % 360 {
	case 90, 270:
		rotated = *image.NewRGBA(image.Rect(0, 0, imgCfg.Height, imgCfg.Width))
		rotatePixels90(originalImg, &rotated, degrees)
	case 180:
		rotated = *image.NewRGBA(image.Rect(0, 0, imgCfg.Width, imgCfg.Height))
		rotatePixels90(originalImg, &rotated, degrees)
	default:
		rotated = *image.NewRGBA(originalImg.Bounds())
		draw.Copy(&rotated, image.Point{}, originalImg, originalImg.Bounds(), draw.Over, nil)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, &rotated, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outFile, &rotated)
	}
	if err != nil {
		return err
	}

	log.Info("Image rotated by %d degrees: %s", degrees, outputPath)
	return nil
}

func ResizeImage(sourcePath string, outputPath string, width, height int) error {
	log := logger.GetLogger()

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
	case "jpeg", "jpg":
		originalImg, err = jpeg.Decode(sourceFile)
	case "png":
		originalImg, err = png.Decode(sourceFile)
	}
	if err != nil {
		return err
	}

	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resized, resized.Bounds(), originalImg, originalImg.Bounds(), draw.Over, nil)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outFile, resized)
	}
	if err != nil {
		return err
	}

	log.Info("Image resized: %s -> %dx%d", outputPath, width, height)
	return nil
}
