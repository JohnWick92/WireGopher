package main

import (
	"flag"
	"log/slog"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/http"
)

var cfg *config.Config
var err error

func Bootstrap() {
	loadConfig()
	runHTTPServer()
}

func loadConfig() {
	configPath := flag.String("config", "/etc/wiregopher/wiregopher.yaml", "path to config file")
	flag.Parse()
	cfg, err = config.Load(*configPath)
	if err != nil {
		config.LogAndExit(err)
	}
	slog.Info("Config loaded")
}

func runHTTPServer() {
	if err := http.Run(&cfg.Server); err != nil {
		slog.Error("Server", "error", err)
	}
}
