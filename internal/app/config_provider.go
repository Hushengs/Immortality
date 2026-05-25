package app

import (
	"context"
	"os"

	"gopkg.in/yaml.v3"
)

type LocalConfigProvider struct {
	path string
}

func NewLocalConfigProvider(path string) *LocalConfigProvider {
	return &LocalConfigProvider{path: path}
}

func (p *LocalConfigProvider) Load(_ context.Context) (RuntimeConfig, error) {
	content, err := os.ReadFile(p.path)
	if err != nil {
		return RuntimeConfig{}, err
	}

	var config RuntimeConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return RuntimeConfig{}, err
	}
	return config, nil
}

var _ ConfigProvider = (*LocalConfigProvider)(nil)
