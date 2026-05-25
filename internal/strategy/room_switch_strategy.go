package strategy

import (
	"time"

	"github.com/Hushengs/Immortality/internal/model"
)

type RoomSwitchStrategy struct {
	maxWait time.Duration
}

func NewRoomSwitchStrategy(maxWaitMinutes int) *RoomSwitchStrategy {
	return &RoomSwitchStrategy{
		maxWait: time.Duration(maxWaitMinutes) * time.Minute,
	}
}

func (s *RoomSwitchStrategy) ShouldSwitch(config model.RuntimeConfig, lastNoFudaiAt time.Time) bool {
	if !config.SwitchRoom {
		return false
	}
	if lastNoFudaiAt.IsZero() {
		return false
	}
	return time.Since(lastNoFudaiAt) >= s.maxWait
}
