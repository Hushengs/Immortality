package actions

import (
	"context"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/device"
)

type LiveNavigation struct {
	logger *slog.Logger
	adb    *device.ADBClient
	config model.RuntimeConfig
}

func NewLiveNavigation(logger *slog.Logger, adb *device.ADBClient, config model.RuntimeConfig) *LiveNavigation {
	return &LiveNavigation{
		logger: logger,
		adb:    adb,
		config: config,
	}
}

func (n *LiveNavigation) Back(ctx context.Context) error {
	n.logger.Info("navigate back")
	return n.adb.KeyEvent(ctx, "4")
}

func (n *LiveNavigation) SwitchRoom(ctx context.Context) error {
	region := n.config.ScreenRegions.SwitchSwipe
	from := model.Point{X: region.X, Y: region.Y}
	to := model.Point{X: region.X, Y: region.Y - region.Height}
	n.logger.Info("switch room by swipe", "from", from, "to", to)
	return n.adb.Swipe(ctx, from, to, 260)
}

func (n *LiveNavigation) EnterFirstRoom(ctx context.Context) error {
	point := n.config.ScreenRegions.FirstRoomTap.Center()
	if point.X == 0 && point.Y == 0 {
		point = model.Point{X: n.config.ScreenRegions.FirstRoomTap.X, Y: n.config.ScreenRegions.FirstRoomTap.Y}
	}
	n.logger.Info("enter first room", "point", point)
	return n.adb.Tap(ctx, point)
}
