package admin

import (
	"context"
	"io/fs"
	"net/http"
)

type Server struct {
	handler http.Handler
}

type Options struct {
	ConfigPath string
	StaticFS   fs.FS
}

func NewServer(options Options) (*Server, error) {
	configService := NewConfigService(options.ConfigPath)
	config, err := configService.Load(context.Background())
	if err != nil {
		return nil, err
	}
	_ = config
	handler := NewHandler(configService)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.Healthz)
	mux.HandleFunc("/api/config", handler.Config)
	mux.HandleFunc("/api/records", handler.Records)
	mux.HandleFunc("/api/events", handler.Events)
	mux.HandleFunc("/api/logs", handler.Logs)
	mux.HandleFunc("/api/screenshots", handler.Screenshots)
	mux.HandleFunc("/api/screenshots/", handler.ScreenshotFile)

	if options.StaticFS != nil {
		fileServer := http.FileServer(http.FS(options.StaticFS))
		mux.Handle("/", fileServer)
	}

	return &Server{handler: mux}, nil
}

func (s *Server) Handler() http.Handler {
	return s.handler
}
