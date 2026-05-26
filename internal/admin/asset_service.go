package admin

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type AssetService struct {
	screenshotDir string
}

type ScreenshotFile struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	URL     string    `json:"url"`
}

func NewAssetService(screenshotDir string) *AssetService {
	return &AssetService{screenshotDir: screenshotDir}
}

func (s *AssetService) ListScreenshots() ([]ScreenshotFile, error) {
	entries, err := os.ReadDir(s.screenshotDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []ScreenshotFile{}, nil
		}
		return nil, err
	}

	files := make([]ScreenshotFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, ScreenshotFile{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			URL:     "/api/screenshots/" + filepath.Base(entry.Name()),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
	return files, nil
}
