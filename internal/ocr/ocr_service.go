package ocr

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Hushengs/Immortality/internal/model"
)

type Service struct {
	tesseractPath string
	config        model.OCRConfig
}

func NewService(tesseractPath string, config model.OCRConfig) *Service {
	return &Service{
		tesseractPath: tesseractPath,
		config:        config,
	}
}

func (s *Service) Recognize(ctx context.Context, imagePath string, lang string) (string, error) {
	if lang == "" {
		lang = s.config.DefaultLang
	}

	if s.config.PreprocessEnabled {
		preprocessedPath, err := preprocessImageFile(imagePath, s.threshold())
		if err == nil {
			defer os.Remove(preprocessedPath)
			imagePath = preprocessedPath
		}
	}

	args := []string{
		imagePath,
		"stdout",
		"-l", lang,
		"--psm", strconv.Itoa(s.config.PSM),
	}

	cmd := exec.CommandContext(ctx, s.tesseractPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return normalizeOCRText(output), nil
}

func normalizeOCRText(raw []byte) string {
	text := strings.TrimSpace(string(bytes.ReplaceAll(raw, []byte("\r"), nil)))
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, " ")
}

func (s *Service) threshold() uint8 {
	if s.config.ThresholdValue == 0 {
		return 160
	}
	return s.config.ThresholdValue
}
