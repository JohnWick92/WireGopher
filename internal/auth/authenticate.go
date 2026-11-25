package auth

import (
	nethttp "net/http"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/identity/keycloak"
)

// Switch between the providers handling
func HandleAuth(cfg *config.Config) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch *cfg.IAM.Provider {
		case "keycloak":
			keycloakRedirectToLogin(cfg, w, r)
		}
	})
}

func keycloakRedirectToLogin(cfg *config.Config, w nethttp.ResponseWriter, r *nethttp.Request) {
	loginURL := keycloak.GetRedirectionURL(cfg.IAM)
	nethttp.Redirect(w, r, loginURL, nethttp.StatusFound)
}
