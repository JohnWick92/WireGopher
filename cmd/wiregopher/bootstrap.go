package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/firewall/opnsense"
	"github.com/JohnWick92/WireGopher/internal/http"
)

var cfg *config.Config
var err error

func Bootstrap() {
	loadConfig()
	setSelectedFirewall()
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

func setSelectedFirewall() {
	switch cfg.Firewall.FirewallModel {
	case "opnsense":
		opnsense.Run(cfg.Firewall.OPNsense)
	default:
		cancelAndExit("Unsupported firewall method")

	}
}

func cancelAndExit(msg string) {
	fmt.Printf("The program exited with the following message: %s\n", msg)
	os.Exit(1)
}
