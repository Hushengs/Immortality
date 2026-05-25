package device

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Hushengs/Immortality/internal/model"
)

type Manager struct {
	logger *slog.Logger
	adb    *ADBClient
}

func NewManager(logger *slog.Logger, adb *ADBClient) *Manager {
	return &Manager{
		logger: logger,
		adb:    adb,
	}
}

func (m *Manager) SelectDevice(ctx context.Context, preferredSerial string) (model.Device, *ADBClient, error) {
	devices, err := m.adb.ListDevices(ctx)
	if err != nil {
		return model.Device{}, nil, err
	}
	if len(devices) == 0 {
		return model.Device{}, nil, errors.New("no adb devices detected")
	}

	for _, item := range devices {
		if preferredSerial != "" && item.Serial != preferredSerial {
			continue
		}
		item.LastSeenAt = time.Now()
		m.logger.Info("selected device", "serial", item.Serial)
		return item, m.adb.WithSerial(item.Serial), nil
	}

	if preferredSerial != "" {
		return model.Device{}, nil, errors.New("preferred device serial not found")
	}

	chosen := devices[0]
	chosen.LastSeenAt = time.Now()
	m.logger.Info("selected first available device", "serial", chosen.Serial)
	return chosen, m.adb.WithSerial(chosen.Serial), nil
}
