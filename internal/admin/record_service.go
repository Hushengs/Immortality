package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/Hushengs/Immortality/internal/model"
)

type RecordService struct {
	recordFile string
}

type PaginationResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func NewRecordService(recordFile string) *RecordService {
	return &RecordService{recordFile: recordFile}
}

func (s *RecordService) ListRecords(_ context.Context, page, pageSize int) (PaginationResult[model.LotteryRecord], error) {
	lines, err := readJSONLLines(s.recordFile)
	if err != nil {
		return PaginationResult[model.LotteryRecord]{}, err
	}

	records := make([]model.LotteryRecord, 0, len(lines))
	for _, line := range lines {
		var record model.LotteryRecord
		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}
		if record.RecordID == "" {
			continue
		}
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})

	return paginate(records, page, pageSize), nil
}

func (s *RecordService) ListEvents(_ context.Context, page, pageSize int) (PaginationResult[model.RuntimeEvent], error) {
	lines, err := readJSONLLines(s.recordFile)
	if err != nil {
		return PaginationResult[model.RuntimeEvent]{}, err
	}

	events := make([]model.RuntimeEvent, 0, len(lines))
	for _, line := range lines {
		var event model.RuntimeEvent
		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}
		if event.Type == "" || event.Message == "" {
			continue
		}
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})

	return paginate(events, page, pageSize), nil
}

func readJSONLLines(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	lines := make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, []byte(line))
	}
	return lines, scanner.Err()
}

func paginate[T any](items []T, page, pageSize int) PaginationResult[T] {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return PaginationResult[T]{Items: []T{}, Total: total, Page: page, PageSize: pageSize}
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return PaginationResult[T]{
		Items:    items[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
