package repository

import (
	"context"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/app"
)

type MemoryDeviceRepository struct {
	logger *slog.Logger
}

func NewMemoryDeviceRepository(logger *slog.Logger) *MemoryDeviceRepository {
	return &MemoryDeviceRepository{logger: logger}
}

func (r *MemoryDeviceRepository) SaveHeartbeat(_ context.Context, heartbeat app.DeviceHeartbeat) error {
	r.logger.Info("device heartbeat recorded", "device_id", heartbeat.DeviceID, "status", heartbeat.Status)
	return nil
}

var _ app.DeviceRepository = (*MemoryDeviceRepository)(nil)
