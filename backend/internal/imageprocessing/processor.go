package imageprocessing

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
)

type ImageProcessor struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		MaxWidth:  1920,
		MaxHeight: 1920,
		Quality:   85,
	}
}

// ProcessImage resizes and optimizes an uploaded image
func (p *ImageProcessor) ProcessImage(src io.Reader, contentType string) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize if needed
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > p.MaxWidth || height > p.MaxHeight {
		img = imaging.Fit(img, p.MaxWidth, p.MaxHeight, imaging.Lanczos)
	}

	// Auto-orient based on EXIF data
	img = imaging.AutoOrientation(img)

	// Encode optimized image
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: p.Quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
	default:
		// Convert to JPEG for other formats
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: p.Quality})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// CreateThumbnail creates a thumbnail version of the image
func (p *ImageProcessor) CreateThumbnail(src io.Reader, size int) ([]byte, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Create thumbnail
	thumbnail := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)

	// Encode
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

// ValidateImage checks if the file is a valid image
func (p *ImageProcessor) ValidateImage(src io.Reader) (bool, string, error) {
	_, format, err := image.Decode(src)
	if err != nil {
		return false, "", err
	}

	validFormats := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
	}

	if !validFormats[format] {
		return false, format, fmt.Errorf("unsupported image format: %s", format)
	}

	return true, format, nil
}
