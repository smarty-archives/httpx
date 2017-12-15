package middleware

import (
	"net"
	"net/http"
	"net/url"

	"github.com/smartystreets/httpx"
)

// Designed to allow cross-domain cookies to be sent by the server and
// accepted by the browser. Whenever a cross-domain request is performed
// and the response contains a Set-Cookie header, this handler must be
// registered for that route.
type AuthenticatedCORSHandler struct {
	inner          http.Handler
	allowedOrigins map[string]bool
}

func NewAuthenticatedCORSHandler(allowedOrigins ...string) *AuthenticatedCORSHandler {
	if len(allowedOrigins) == 0 {
		allowedOrigins = defaultCORSOrigins
	}
	allowed := make(map[string]bool)
	for _, origin := range allowedOrigins {
		allowed[origin] = true
	}
	return &AuthenticatedCORSHandler{
		inner:          NoopHandler{},
		allowedOrigins: allowed,
	}
}

func (this *AuthenticatedCORSHandler) Install(handler http.Handler) {
	this.inner = handler
}

func (this *AuthenticatedCORSHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	headers := response.Header()

	if origin := request.Header.Get("Origin"); this.isAllowed(origin) {
		headers.Set(httpx.HeaderAccessControlAllowOrigin, origin)
		headers.Set(httpx.HeaderAccessControlAllowCredentials, "true")
	}

	this.inner.ServeHTTP(response, request)
}

func (this *AuthenticatedCORSHandler) isAllowed(origin string) bool {
	return this.allowedOrigins[extractHostname(origin)]
}

func extractHostname(origin string) string {
	parsed, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	if host, _, err := net.SplitHostPort(parsed.Host); err == nil {
		return host
	} else {
		return parsed.Host
	}
}

var defaultCORSOrigins = []string{
	"localhost",
	"smartystreets.dev",
	"smartystreets.com",
	"storefront.smartystreets.com",
}
