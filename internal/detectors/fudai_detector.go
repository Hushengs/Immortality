package detectors

import (
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/Hushengs/Immortality/internal/capture"
	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/ocr"
)

type FudaiDetector struct {
	config model.RuntimeConfig
	ocr    *ocr.Service
}

func NewFudaiDetector(config model.RuntimeConfig, ocrService *ocr.Service) *FudaiDetector {
	return &FudaiDetector{
		config: config,
		ocr:    ocrService,
	}
}

func (d *FudaiDetector) Detect(ctx context.Context, screenPath string) (model.FudaiDetection, error) {
	region := d.config.ScreenRegions.FudaiIcon
	ratio, err := computeRedRatio(screenPath, region)
	if err != nil {
		return model.FudaiDetection{}, err
	}

	if ratio >= d.config.Thresholds.FudaiRedRatio {
		return model.FudaiDetection{
			Found:        true,
			MatchScore:   ratio,
			Reason:       "matched red pixel ratio in fixed area",
			ClickPoint:   region.Center(),
			SampleRegion: region,
		}, nil
	}

	tempFile, err := os.CreateTemp("", "fudai_region_*.png")
	if err != nil {
		return model.FudaiDetection{}, err
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := capture.CropImage(screenPath, region, tempPath); err != nil {
		return model.FudaiDetection{}, err
	}

	text, err := d.ocr.Recognize(ctx, tempPath, d.config.OCR.DefaultLang)
	if err != nil {
		return model.FudaiDetection{}, nil
	}

	if containsAny(text, d.config.Keywords.Fudai) {
		return model.FudaiDetection{
			Found:        true,
			MatchScore:   0.5,
			Reason:       "matched fudai keyword by OCR",
			ClickPoint:   region.Center(),
			SampleRegion: region,
		}, nil
	}

	return model.FudaiDetection{
		Found:        false,
		MatchScore:   ratio,
		Reason:       "no fudai signal detected",
		ClickPoint:   region.Center(),
		SampleRegion: region,
	}, nil
}

func computeRedRatio(screenPath string, region model.Region) (float64, error) {
	file, err := os.Open(screenPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	bounds := img.Bounds()
	startX := max(region.X, bounds.Min.X)
	startY := max(region.Y, bounds.Min.Y)
	endX := min(region.X+region.Width, bounds.Max.X)
	endY := min(region.Y+region.Height, bounds.Max.Y)

	if startX >= endX || startY >= endY {
		return 0, nil
	}

	total := 0
	red := 0
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			total++
			if r > 46000 && g < 28000 && b < 28000 {
				red++
			}
		}
	}
	if total == 0 {
		return 0, nil
	}
	return float64(red) / float64(total), nil
}

func containsAny(text string, keywords []string) bool {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
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
