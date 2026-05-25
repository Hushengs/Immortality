package capture

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hushengs/Immortality/internal/model"
)

func CropImage(sourcePath string, region model.Region, outputPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	img, format, err := image.Decode(source)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	rect := image.Rect(
		max(region.X, bounds.Min.X),
		max(region.Y, bounds.Min.Y),
		min(region.X+region.Width, bounds.Max.X),
		min(region.Y+region.Height, bounds.Max.Y),
	)
	if rect.Empty() {
		rect = bounds
	}

	cropped := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(cropped, cropped.Bounds(), img, rect.Min, draw.Src)

	target, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer target.Close()

	switch normalizedFormat(format, outputPath) {
	case "jpg", "jpeg":
		return jpeg.Encode(target, cropped, &jpeg.Options{Quality: 92})
	default:
		return png.Encode(target, cropped)
	}
}

func normalizedFormat(decodedFormat, outputPath string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(outputPath)), ".")
	if ext != "" {
		return ext
	}
	return strings.ToLower(decodedFormat)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
