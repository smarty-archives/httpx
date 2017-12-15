package middleware

import (
	"net/http"

	"github.com/smartystreets/httpx"
)

type HeadersHandler struct {
	inner       http.Handler
	headers     map[string]string
	alwaysWrite bool
}

func DefaultCORSHeadersHandler() *HeadersHandler {
	return OriginCORSHeadersHandler("*")
}
func OriginCORSHeadersHandler(origin string) *HeadersHandler {
	return CORSHeadersHandler(map[string]string{
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Methods":     "GET, PUT, POST, DELETE, HEAD",
		"Access-Control-Allow-Headers":     "Accept, Content-Type, Content-Length, Referer, Origin, Host",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Max-Age":           "600",
	})
}

func CORSHeadersHandler(headers map[string]string) *HeadersHandler {
	return newHeadersHandler(headers, true)
}

func BrowserHeadersHandler(headers map[string]string) *HeadersHandler {
	return newHeadersHandler(headers, false)
}

func NewHeadersHandler(headers map[string]string) *HeadersHandler {
	return newHeadersHandler(headers, true)
}

func newHeadersHandler(headers map[string]string, alwaysWrite bool) *HeadersHandler {
	return &HeadersHandler{headers: headers, alwaysWrite: alwaysWrite, inner: NoopHandler{}}
}

func (this *HeadersHandler) Install(inner http.Handler) {
	this.inner = inner
}

func (this *HeadersHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if this.canWriteHeaders(request) {
		this.writeHeaders(response)
	}

	this.inner.ServeHTTP(response, request)
}

func (this *HeadersHandler) canWriteHeaders(request *http.Request) bool {
	return this.alwaysWrite ||
		len(httpx.ReadHeader(request, "Origin")) > 0 ||
		len(httpx.ReadHeader(request, "Referer")) > 0
}

func (this *HeadersHandler) writeHeaders(response http.ResponseWriter) {
	headers := response.Header()
	for key, value := range this.headers {
		if len(value) > 0 {
			headers.Set(key, value)
		}
	}
}
