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
	allowedOrigins map[string]struct{}
}

func NewAuthenticatedCORSHandler(allowedOrigins ...string) *AuthenticatedCORSHandler {
	if len(allowedOrigins) == 0 {
		allowedOrigins = defaultCORSOrigins
	}

	allowed := make(map[string]struct{})
	this := &AuthenticatedCORSHandler{inner: NoopHandler{}, allowedOrigins: allowed}
	return this.AppendAllowedOrigin(allowedOrigins...)
}

func (this *AuthenticatedCORSHandler) AppendAllowedOrigin(allowedOrigins ...string) *AuthenticatedCORSHandler {
	for _, value := range allowedOrigins {
		this.allowedOrigins[value] = struct{}{}
	}

	return this
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
	_, contains := this.allowedOrigins[extractHostname(origin)]
	return contains
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

var defaultCORSOrigins = []string{"localhost", "smartystreets.dev", "smartystreets.com"}
