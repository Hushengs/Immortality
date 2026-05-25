package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func EnsureDir(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}

func EnsureArtifactDirs(logDir, screenshotDir, recordFile string) error {
	if err := EnsureDir(logDir); err != nil {
		return err
	}
	if err := EnsureDir(screenshotDir); err != nil {
		return err
	}
	if recordFile == "" {
		return nil
	}
	return EnsureDir(filepath.Dir(recordFile))
}

func NewTimestampedPath(dir, prefix, ext string) string {
	name := fmt.Sprintf("%s_%s%s", prefix, time.Now().Format("20060102_150405"), ext)
	return filepath.Join(dir, name)
}
