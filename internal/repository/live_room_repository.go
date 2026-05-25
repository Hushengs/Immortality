package repository

import (
	"context"
	"log/slog"

	"github.com/Hushengs/Immortality/internal/app"
)

type MemoryLiveRoomRepository struct {
	logger *slog.Logger
}

func NewMemoryLiveRoomRepository(logger *slog.Logger) *MemoryLiveRoomRepository {
	return &MemoryLiveRoomRepository{logger: logger}
}

func (r *MemoryLiveRoomRepository) Save(_ context.Context, room app.LiveRoom) error {
	r.logger.Info("live room saved", "room_id", room.RoomID, "status", room.Status, "title", room.Title)
	return nil
}

var _ app.LiveRoomRepository = (*MemoryLiveRoomRepository)(nil)
