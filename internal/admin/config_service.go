package admin

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/Hushengs/Immortality/internal/model"
	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	path string
	mu   sync.Mutex
}

func NewConfigService(path string) *ConfigService {
	return &ConfigService{path: path}
}

func (s *ConfigService) Load(_ context.Context) (model.RuntimeConfig, error) {
	content, err := os.ReadFile(s.path)
	if err != nil {
		return model.RuntimeConfig{}, err
	}

	var config model.RuntimeConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return model.RuntimeConfig{}, err
	}
	return config, nil
}

func (s *ConfigService) Save(_ context.Context, config model.RuntimeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path, content, 0o644)
}
