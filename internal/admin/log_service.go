package admin

import (
	"bufio"
	"os"
)

type LogService struct {
	path string
}

func NewLogService(path string) *LogService {
	return &LogService{path: path}
}

func (s *LogService) TailLines(limit int) ([]string, error) {
	if limit <= 0 {
		limit = 200
	}

	file, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) <= limit {
		return lines, nil
	}
	return lines[len(lines)-limit:], nil
}
