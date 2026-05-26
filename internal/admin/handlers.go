package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Hushengs/Immortality/internal/model"
)

type Handler struct {
	configService *ConfigService
}

func NewHandler(configService *ConfigService) *Handler {
	return &Handler{
		configService: configService,
	}
}

func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}

func (h *Handler) Config(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		config, err := h.configService.Load(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, config)
	case http.MethodPut:
		var config model.RuntimeConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := h.configService.Save(r.Context(), config); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"saved": true})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Records(w http.ResponseWriter, r *http.Request) {
	page, pageSize := readPagination(r)
	config, err := h.configService.Load(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	result, err := NewRecordService(config.Artifacts.RecordFile).ListRecords(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	page, pageSize := readPagination(r)
	config, err := h.configService.Load(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	result, err := NewRecordService(config.Artifacts.RecordFile).ListEvents(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	config, err := h.configService.Load(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	lines, err := NewLogService(filepath.Join(config.Artifacts.LogDir, "runtime.log")).TailLines(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": lines,
		"total": len(lines),
	})
}

func (h *Handler) Screenshots(w http.ResponseWriter, r *http.Request) {
	config, err := h.configService.Load(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	items, err := NewAssetService(config.Artifacts.ScreenshotDir).ListScreenshots()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func (h *Handler) ScreenshotFile(w http.ResponseWriter, r *http.Request) {
	config, err := h.configService.Load(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/api/screenshots/")
	name = filepath.Base(name)
	if name == "." || name == "" {
		writeError(w, http.StatusBadRequest, errors.New("invalid screenshot name"))
		return
	}
	http.ServeFile(w, r, filepath.Join(config.Artifacts.ScreenshotDir, name))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}

func readPagination(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	return page, pageSize
}
