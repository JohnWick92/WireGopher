package rewrite

import "net/http"

// KeycloakOptions are the Keycloak-specific inputs.
type KeycloakOptions struct {
	Options
	ForceOKOnKerberos bool
}

// KeycloakRewrite applies the Keycloak-specific rewrites on top of the neutral ones.
func KeycloakRewrite(resp *http.Response, opts KeycloakOptions) error {
	delete(resp.Header, "Www-Authenticate")
	delete(resp.Header, "Content-Security-Policy")
	delete(resp.Header, "Content-Length")
	delete(resp.Header, "Content-Encoding")

	if opts.ForceOKOnKerberos {
		opts.StatusOverride = map[string]int{"Kerberos Unsupported": http.StatusOK}
	}

	opts.SkipPrefixPaths = []string{"/callback"}

	Location(resp, opts.Options)
	Cookies(resp, opts.Options)
	return Body(resp, opts.Options)
}
