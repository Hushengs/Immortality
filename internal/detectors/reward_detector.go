package detectors

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Hushengs/Immortality/internal/capture"
	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/ocr"
)

type RewardDetector struct {
	config model.RuntimeConfig
	ocr    *ocr.Service
}

func NewRewardDetector(config model.RuntimeConfig, ocrService *ocr.Service) *RewardDetector {
	return &RewardDetector{
		config: config,
		ocr:    ocrService,
	}
}

func (d *RewardDetector) Analyze(ctx context.Context, screenPath string, layout model.PopupLayout) (model.PopupAnalysis, error) {
	prizeText, err := d.readRegion(ctx, screenPath, layout.PrizeRegion, d.config.OCR.DefaultLang)
	if err != nil {
		return model.PopupAnalysis{}, err
	}

	countdownText, err := d.readRegion(ctx, screenPath, layout.CountdownRegion, d.config.OCR.DigitsLang)
	if err != nil {
		return model.PopupAnalysis{}, err
	}

	buttonText, err := d.readRegion(ctx, screenPath, layout.ButtonRegion, d.config.OCR.DefaultLang)
	if err != nil {
		return model.PopupAnalysis{}, err
	}

	return model.PopupAnalysis{
		PrizeText:        prizeText,
		CountdownText:    countdownText,
		CountdownSeconds: parseCountdownSeconds(countdownText),
		ButtonText:       buttonText,
		NeedTask:         containsAny(prizeText+" "+buttonText, []string{"任务", "完成任务"}),
		NeedFansClub:     containsAny(prizeText+" "+buttonText, []string{"粉丝团", "粉丝"}),
		CannotJoin:       containsAny(buttonText, []string{"无法参与", "已结束", "不可参与"}),
		AlreadyJoined:    containsAny(buttonText, []string{"已参与", "参与成功"}),
		LayoutName:       layout.Name,
	}, nil
}

func (d *RewardDetector) DetectResult(ctx context.Context, screenPath string) (string, error) {
	layout := d.config.ScreenRegions.PopupPrize
	text, err := d.readRegion(ctx, screenPath, layout, d.config.OCR.DefaultLang)
	if err != nil {
		return "", err
	}

	switch {
	case containsAny(text, d.config.Keywords.Won):
		return "won", nil
	case containsAny(text, []string{"未中奖", "很遗憾", "下次再来"}):
		return "lost", nil
	default:
		return "unknown", nil
	}
}

func (d *RewardDetector) readRegion(ctx context.Context, screenPath string, region model.Region, lang string) (string, error) {
	tempFile, err := os.CreateTemp("", "ocr_region_*.png")
	if err != nil {
		return "", err
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := capture.CropImage(screenPath, region, tempPath); err != nil {
		return "", err
	}
	return d.ocr.Recognize(ctx, tempPath, lang)
}

func parseCountdownSeconds(text string) int {
	cleaned := strings.ReplaceAll(text, "：", ":")
	reMinuteSecond := regexp.MustCompile(`(\d+)\s*[:分]\s*(\d+)`)
	if match := reMinuteSecond.FindStringSubmatch(cleaned); len(match) == 3 {
		minute, _ := strconv.Atoi(match[1])
		second, _ := strconv.Atoi(match[2])
		return (minute * 60) + second
	}

	reNumber := regexp.MustCompile(`\d+`)
	values := reNumber.FindAllString(cleaned, -1)
	if len(values) == 1 {
		number, _ := strconv.Atoi(values[0])
		return number
	}
	if len(values) >= 2 {
		minute, _ := strconv.Atoi(values[0])
		second, _ := strconv.Atoi(values[1])
		return (minute * 60) + second
	}
	return 0
}
