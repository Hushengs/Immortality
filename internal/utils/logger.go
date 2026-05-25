package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewLogger(logDir string) (*slog.Logger, func() error, error) {
	if err := EnsureDir(logDir); err != nil {
		return nil, nil, err
	}

	logPath := filepath.Join(logDir, "runtime.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return logger, file.Close, nil
}
