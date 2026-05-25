package actions

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/model"
	"github.com/Hushengs/Immortality/internal/device"
)

type RewardActions struct {
	logger *slog.Logger
	adb    *device.ADBClient
	config model.RuntimeConfig
}

func NewRewardActions(logger *slog.Logger, adb *device.ADBClient, config model.RuntimeConfig) *RewardActions {
	return &RewardActions{
		logger: logger,
		adb:    adb,
		config: config,
	}
}

func (a *RewardActions) Claim(ctx context.Context, claimRegion model.Region) error {
	if a.config.SafeMode || !a.config.AutoClaimReward {
		a.logger.Warn("skip auto reward claim because safe mode or auto claim is disabled")
		return nil
	}
	if a.config.AutoPlaceOrder {
		return errors.New("auto place order is intentionally not implemented in mvp")
	}

	a.logger.Info("claim reward", "point", claimRegion.Center())
	return a.adb.Tap(ctx, claimRegion.Center())
}
