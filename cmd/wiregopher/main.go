package main

import (
	"flag"
	"log/slog"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/http"
)

func main() {
	configPath := flag.String("config", "/etc/wiregopher/wiregopher.yaml", "path to config file")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		config.LogAndExit(err)
	}
	if err := http.Run(cfg); err != nil {
		slog.Error("Server", "error", err)
	}
}
