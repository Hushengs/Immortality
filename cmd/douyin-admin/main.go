package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Hushengs/Immortality/internal/admin"
)

func main() {
	addr := flag.String("addr", ":18080", "http listen address")
	configPath := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	staticDir := filepath.Join(filepath.Dir(*configPath), "web", "admin")

	server, err := admin.NewServer(admin.Options{
		ConfigPath: *configPath,
		StaticFS:   os.DirFS(staticDir),
	})
	if err != nil {
		log.Fatalf("init admin server failed: %v", err)
	}

	log.Printf("douyin-admin listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
