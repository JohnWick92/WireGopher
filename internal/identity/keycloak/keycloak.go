package keycloak

import (
	"net/url"

	"github.com/JohnWick92/WireGopher/internal/config"
)

func GetRedirectionURL(cfg *config.IAM) string {
	values := url.Values{}
	values.Set("client_id", *cfg.Keycloak.ClientID)
	values.Set("response_type", "code")
	values.Set("redirect_uri", *cfg.Keycloak.RedirectURI)
	values.Set("scope", "openid profile")

	return "/iam/realms/" + *cfg.Keycloak.Realm +
		"/protocol/openid-connect/auth?" + values.Encode()
}
