package http

import (
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/http/handlers"
)

func registerRoutes(mux *nethttp.ServeMux) {
	mux.HandleFunc("GET /health", handlers.GetHealth)
}
