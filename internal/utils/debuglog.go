package utils

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

const debugLogPath = "debug-1fa3ce.log"

var debugLogMu sync.Mutex

type debugLogEntry struct {
	SessionID    string         `json:"sessionId"`
	RunID        string         `json:"runId,omitempty"`
	HypothesisID string         `json:"hypothesisId"`
	Location     string         `json:"location"`
	Message      string         `json:"message"`
	Data         map[string]any `json:"data,omitempty"`
	Timestamp    int64          `json:"timestamp"`
}

func WriteDebugLog(runID, hypothesisID, location, message string, data map[string]any) {
	debugLogMu.Lock()
	defer debugLogMu.Unlock()

	file, err := os.OpenFile(debugLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	entry := debugLogEntry{
		SessionID:    "1fa3ce",
		RunID:        runID,
		HypothesisID: hypothesisID,
		Location:     location,
		Message:      message,
		Data:         data,
		Timestamp:    time.Now().UnixMilli(),
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return
	}

	_, _ = file.Write(append(encoded, '\n'))
}
