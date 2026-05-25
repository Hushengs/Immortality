package detectors

import (
	"context"
	"os"

	"github.com/Hushengs/Immortality/internal/capture"
	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/ocr"
)

type VerificationDetector struct {
	config model.RuntimeConfig
	ocr    *ocr.Service
}

func NewVerificationDetector(config model.RuntimeConfig, ocrService *ocr.Service) *VerificationDetector {
	return &VerificationDetector{
		config: config,
		ocr:    ocrService,
	}
}

func (d *VerificationDetector) Detect(ctx context.Context, screenPath string) (bool, string, error) {
	tempFile, err := os.CreateTemp("", "verification_*.png")
	if err != nil {
		return false, "", err
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := capture.CropImage(screenPath, d.config.ScreenRegions.Verification, tempPath); err != nil {
		return false, "", err
	}

	text, err := d.ocr.Recognize(ctx, tempPath, d.config.OCR.DefaultLang)
	if err != nil {
		return false, "", err
	}

	return containsAny(text, d.config.Keywords.Verification), text, nil
}
