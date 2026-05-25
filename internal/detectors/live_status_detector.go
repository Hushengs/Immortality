package detectors

import (
	"context"
	"os"

	"github.com/Hushengs/Immortality/internal/capture"
	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/ocr"
)

type LiveStatusDetector struct {
	config model.RuntimeConfig
	ocr    *ocr.Service
}

func NewLiveStatusDetector(config model.RuntimeConfig, ocrService *ocr.Service) *LiveStatusDetector {
	return &LiveStatusDetector{
		config: config,
		ocr:    ocrService,
	}
}

func (d *LiveStatusDetector) Detect(ctx context.Context, screenPath string) (model.PageStatus, error) {
	tempFile, err := os.CreateTemp("", "live_status_*.png")
	if err != nil {
		return model.PageStatus{}, err
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := capture.CropImage(screenPath, d.config.ScreenRegions.LiveStatus, tempPath); err != nil {
		return model.PageStatus{}, err
	}

	text, err := d.ocr.Recognize(ctx, tempPath, d.config.OCR.DefaultLang)
	if err != nil {
		return model.PageStatus{}, err
	}

	status := model.PageStatus{
		Kind:           "live_room",
		IsLiveRoom:     true,
		IsLiveEnded:    containsAny(text, d.config.Keywords.LiveEnded),
		IsVerification: containsAny(text, d.config.Keywords.Verification),
		IsListPage:     containsAny(text, []string{"直播", "关注", "推荐"}),
		IsFollowPage:   containsAny(text, []string{"关注", "正在直播"}),
	}
	if status.IsVerification {
		status.Kind = "verification"
		status.IsLiveRoom = false
	}
	if status.IsLiveEnded {
		status.Kind = "live_ended"
		status.IsLiveRoom = false
	}

	return status, nil
}
