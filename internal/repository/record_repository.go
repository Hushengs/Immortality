package repository

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/Hushengs/Immortality/internal/app"
)

type FileRecordRepository struct {
	path string
	mu   sync.Mutex
}

func NewFileRecordRepository(path string) *FileRecordRepository {
	return &FileRecordRepository{path: path}
}

func (r *FileRecordRepository) SaveRecord(_ context.Context, record app.LotteryRecord) error {
	return r.append(record)
}

func (r *FileRecordRepository) SaveEvent(_ context.Context, event app.RuntimeEvent) error {
	return r.append(event)
}

func (r *FileRecordRepository) append(v any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = file.Write(append(encoded, '\n'))
	return err
}

var _ app.RecordRepository = (*FileRecordRepository)(nil)
