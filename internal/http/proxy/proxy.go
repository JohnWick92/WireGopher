package proxy

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/JohnWick92/WireGopher/internal/config"
	"github.com/JohnWick92/WireGopher/internal/http/proxy/rewrite"
)

// IAMProxy builds the reverse proxy handler for the configured IAM provider.
func IAMProxy(cfg *config.Config) http.Handler {
	middlewareURL := cfg.Server.HTTP.URL

	var (
		upstreamURL string
		proxyPath   string
		useHTTPS    bool
		rewriteResp func(*http.Response) error
	)

	switch *cfg.IAM.Provider {
	case "keycloak":
		upstreamURL = *cfg.IAM.Keycloak.URL
		proxyPath = "/iam"
		useHTTPS = false // ⚠️ TESTING (HTTP): false | 🔒 PRODUCTION: true

		opts := rewrite.Options{
			UpstreamURL:   upstreamURL,
			MiddlewareURL: middlewareURL,
			ProxyPath:     proxyPath,
			UseHTTPS:      useHTTPS,
		}

		rewriteResp = func(resp *http.Response) error {
			return rewrite.KeycloakRewrite(resp, rewrite.KeycloakOptions{
				Options:           opts,
				ForceOKOnKerberos: true,
			})
		}

	default:
		slog.Error("unsupported iam provider", "provider", *cfg.IAM.Provider)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unsupported iam provider", http.StatusInternalServerError)
		})
	}

	target, err := url.Parse(upstreamURL)
	if err != nil {
		slog.Error("invalid iam upstream url", "url", upstreamURL, "error", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid iam upstream url", http.StatusInternalServerError)
		})
	}

	return &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("proxy error", "error", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("Proxy error: " + err.Error()))
		},
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)

			// Strip the /iam prefix from the incoming path.
			pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, "/iam")
			if pr.Out.URL.Path == "" {
				pr.Out.URL.Path = "/"
			}
			pr.Out.URL.RawQuery = pr.In.URL.RawQuery

			// Force the scheme on redirect_uri if it's missing.
			if ru := pr.Out.URL.Query().Get("redirect_uri"); ru != "" {
				scheme := "http://"
				if useHTTPS {
					scheme = "https://"
				}

				q := pr.Out.URL.Query()
				q.Set("redirect_uri", normalizeRedirectURI(ru, scheme))
				pr.Out.URL.RawQuery = q.Encode()
			}

			// Headers.
			pr.Out.Host = target.Host
			pr.Out.Header.Set("Host", target.Host)
			publicURL, err := url.Parse(middlewareURL)
			if err == nil {
				pr.Out.Header.Set("X-Forwarded-Host", publicURL.Host)
			}

			// ⚠️ TESTING (HTTP): "http" | 🔒 PRODUCTION: "https"
			if useHTTPS {
				pr.Out.Header.Set("X-Forwarded-Proto", "https")
			} else {
				pr.Out.Header.Set("X-Forwarded-Proto", "http")
			}

			pr.Out.Header.Set("Accept-Encoding", "identity")

			if xff := pr.In.Header.Get("X-Forwarded-For"); xff != "" {
				pr.Out.Header.Set("X-Forwarded-For", xff)
			} else {
				pr.Out.Header.Set("X-Forwarded-For", pr.In.RemoteAddr)
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			return rewriteResp(resp)
		},
	}
}

func normalizeRedirectURI(raw, scheme string) string {
	raw = strings.TrimSpace(raw)

	for {
		switch {
		case strings.HasPrefix(raw, "http://"):
			raw = strings.TrimPrefix(raw, "http://")
		case strings.HasPrefix(raw, "https://"):
			raw = strings.TrimPrefix(raw, "https://")
		case strings.HasPrefix(raw, "http//"):
			raw = strings.TrimPrefix(raw, "http//")
		case strings.HasPrefix(raw, "https//"):
			raw = strings.TrimPrefix(raw, "https//")
		default:
			return scheme + strings.TrimLeft(raw, "/")
		}
	}
}
