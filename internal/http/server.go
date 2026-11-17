package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/config"
)

var multiplexer = nethttp.NewServeMux()
var server *nethttp.Server

func createServer(cfg *config.Config) {
	server = &nethttp.Server{
		Addr:    ":" + cfg.Server.HTTP.Port,
		Handler: multiplexer,
	}
	registerRoutes(multiplexer)
}

func Run(cfg *config.Config) error {
	createServer(cfg)
	slog.Info("Listening on", "port", server.Addr)
	return server.ListenAndServe()
}
