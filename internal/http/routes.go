package http

import (
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/auth"
	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/http/handlers"
	"github.com/JohnWick92/WireGopher/internal/http/proxy"
)

func registerRoutes(mux *nethttp.ServeMux, cfg *config.Config) {

	mux.Handle("GET /auth", auth.HandleAuth(cfg))
	mux.HandleFunc("GET /health", handlers.GetHealth)
	mux.Handle("/iam/{remainder...}", proxy.IAMProxy(cfg))
}
