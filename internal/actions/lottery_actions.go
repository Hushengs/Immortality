package actions

import (
	"context"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/device"
)

type LotteryActions struct {
	logger *slog.Logger
	adb    *device.ADBClient
}

func NewLotteryActions(logger *slog.Logger, adb *device.ADBClient) *LotteryActions {
	return &LotteryActions{
		logger: logger,
		adb:    adb,
	}
}

func (a *LotteryActions) OpenFudaiDetail(ctx context.Context, detection model.FudaiDetection) error {
	a.logger.Info("open fudai detail", "reason", detection.Reason, "score", detection.MatchScore)
	return a.adb.Tap(ctx, detection.ClickPoint)
}

func (a *LotteryActions) ClickJoin(ctx context.Context, region model.Region) error {
	point := region.Center()
	a.logger.Info("click join button", "point", point)
	return a.adb.Tap(ctx, point)
}
