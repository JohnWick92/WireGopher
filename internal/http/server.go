package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/config"
)

var multiplexer = nethttp.NewServeMux()
var server *nethttp.Server

func createServer(cfg *config.Server) {
	server = &nethttp.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: multiplexer,
	}
	registerRoutes(multiplexer)
}

func Run(cfg *config.Server) error {
	createServer(cfg)
	slog.Info("Listening on", "port", server.Addr)
	return server.ListenAndServe()
}
