package integration

import (
	"context"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/app"
)

type NoopManagementClient struct {
	logger *slog.Logger
}

func NewNoopManagementClient(logger *slog.Logger) *NoopManagementClient {
	return &NoopManagementClient{logger: logger}
}

func (c *NoopManagementClient) ReportEvent(_ context.Context, event app.RuntimeEvent) error {
	c.logger.Debug("skip remote event report", "type", event.Type)
	return nil
}

func (c *NoopManagementClient) ReportHeartbeat(_ context.Context, heartbeat app.DeviceHeartbeat) error {
	c.logger.Debug("skip remote heartbeat", "device_id", heartbeat.DeviceID)
	return nil
}

var _ app.ManagementReporter = (*NoopManagementClient)(nil)
