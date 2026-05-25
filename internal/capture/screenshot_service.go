package capture

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/Hushengs/Immortality/internal/device"
	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/utils"
)

type ScreenshotService struct {
	logger        *slog.Logger
	adb           *device.ADBClient
	screenshotDir string
	debug         model.DebugConfig
}

func NewScreenshotService(logger *slog.Logger, adb *device.ADBClient, screenshotDir string, debug model.DebugConfig) *ScreenshotService {
	return &ScreenshotService{
		logger:        logger,
		adb:           adb,
		screenshotDir: screenshotDir,
		debug:         debug,
	}
}

func (s *ScreenshotService) Capture(ctx context.Context, prefix string) (string, error) {
	path := utils.NewTimestampedPath(s.screenshotDir, prefix, ".png")
	if s.debug.StaticScreenPath != "" {
		if err := copyFile(s.debug.StaticScreenPath, path); err != nil {
			return "", err
		}
		s.logger.Info("static screen copied", "source", s.debug.StaticScreenPath, "path", path)
		return path, nil
	}

	if err := s.adb.Screencap(ctx, path); err != nil {
		return "", err
	}
	s.logger.Info("screen captured", "path", path)
	return path, nil
}

func copyFile(sourcePath, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}
