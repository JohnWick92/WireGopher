package http

import (
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/http/handlers"
)

func registerRoutes(mux *nethttp.ServeMux, cfg *config.Config) {
	iamProviders := map[string]nethttp.HandlerFunc{
		"keycloak": handlers.AuthKeycloak,
	}

	selectedHandler, _ := iamProviders[*cfg.IAM.Provider]

	mux.HandleFunc("GET /login", selectedHandler)
	mux.HandleFunc("GET /health", handlers.GetHealth)
}
