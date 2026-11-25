package rewrite

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Options struct {
	UpstreamURL     string
	MiddlewareURL   string
	ProxyPath       string
	UseHTTPS        bool
	StatusOverride  map[string]int
	SkipPrefixPaths []string
}

var processableContentTypes = []string{
	"text/html",
	"application/javascript",
	"application/x-javascript",
	"text/css",
	"application/json",
	"image/",
}

func Body(resp *http.Response, opts Options) error {
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || !isProcessable(contentType) {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	bodyStr := string(body)

	for needle, code := range opts.StatusOverride {
		if strings.Contains(bodyStr, needle) {
			resp.StatusCode = code
		}
	}

	proxyURL := opts.MiddlewareURL + opts.ProxyPath

	if opts.UpstreamURL != "" {
		re := regexp.MustCompile(regexp.QuoteMeta(opts.UpstreamURL))
		bodyStr = re.ReplaceAllString(bodyStr, proxyURL)
	}

	dirs := `admin|realms|resources|theme|auth|welcome-content|js|css|img|fonts|common|lib|node_modules|api|clients|v2|static|custom|favicon|meta|login_logo|login-logo|univention|logo|branding|icons|assets`
	pattern := `(["'\(\s=])\s*\/\s*(` + dirs + `)([\/\.\? "'\)=])`
	re2 := regexp.MustCompile(pattern)
	bodyStr = re2.ReplaceAllString(bodyStr, `$1`+opts.ProxyPath+`/$2$3`)

	bodyStr = strings.ReplaceAll(bodyStr, opts.ProxyPath+opts.ProxyPath+"/", opts.ProxyPath+"/")
	for _, skip := range opts.SkipPrefixPaths {
		bodyStr = strings.ReplaceAll(bodyStr, opts.ProxyPath+skip+"/", skip+"/")
	}

	resp.Body = io.NopCloser(strings.NewReader(bodyStr))
	resp.ContentLength = int64(len(bodyStr))
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyStr)))
	return nil
}

func Cookies(resp *http.Response, opts Options) bool {
	cookieHeaders := resp.Header.Values("Set-Cookie")
	if len(cookieHeaders) == 0 {
		return false
	}

	var processedCookies []string
	modified := false

	for _, cookieHeader := range cookieHeaders {
		attributes := strings.Split(cookieHeader, ";")
		var cleanedAttributes []string

		for _, attr := range attributes {
			trimmedAttr := strings.TrimSpace(attr)
			if trimmedAttr == "" {
				continue
			}
			attrLower := strings.ToLower(trimmedAttr)

			if strings.HasPrefix(attrLower, "domain=") {
				modified = true
				continue
			}

			if !opts.UseHTTPS {
				if attrLower == "secure" || strings.HasPrefix(attrLower, "samesite=") {
					modified = true
					continue
				}
			}

			if strings.HasPrefix(attrLower, "path=") {
				originalPath := trimmedAttr[5:]
				if opts.ProxyPath != "" && !strings.HasPrefix(originalPath, opts.ProxyPath) {
					trimmedAttr = "Path=" + opts.ProxyPath + originalPath
					modified = true
				}
			}

			cleanedAttributes = append(cleanedAttributes, trimmedAttr)
		}

		if len(cleanedAttributes) > 0 {
			processedCookies = append(processedCookies, strings.Join(cleanedAttributes, "; "))
		}
	}

	if len(processedCookies) > 0 {
		resp.Header.Del("Set-Cookie")
		for _, cookie := range processedCookies {
			resp.Header.Add("Set-Cookie", cookie)
		}
	}

	return modified
}

func Location(resp *http.Response, opts Options) bool {
	location := resp.Header.Get("Location")
	if location == "" {
		return false
	}

	if opts.ProxyPath != "" && strings.HasPrefix(location, opts.ProxyPath) {
		return false
	}

	for _, skip := range opts.SkipPrefixPaths {
		if strings.HasPrefix(location, skip) {
			return false
		}
	}

	proxyURL := opts.MiddlewareURL + opts.ProxyPath

	switch {
	case opts.UpstreamURL != "" && strings.Contains(location, opts.UpstreamURL):
		location = strings.ReplaceAll(location, opts.UpstreamURL, proxyURL)
	case strings.HasPrefix(location, "/") && !strings.HasPrefix(location, "//"):
		location = opts.ProxyPath + location
	}

	resp.Header.Set("Location", location)
	return true
}

func isProcessable(contentType string) bool {
	for _, ct := range processableContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}
