package middleware

import (
	"net/http"
	"strings"

	"github.com/smartystreets/httpx"
)

type ContentTypeOverrideHandler struct {
	inner    http.Handler
	override string
}

func NewContentTypeOverrideHandler(overrideType string) *ContentTypeOverrideHandler {
	return &ContentTypeOverrideHandler{
		inner:    NoopHandler{},
		override: overrideType,
	}
}

func (this *ContentTypeOverrideHandler) Install(handler http.Handler) {
	this.inner = handler
}

func (this *ContentTypeOverrideHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	contentType := httpx.ReadHeader(request, httpx.HeaderContentType)
	if strings.Contains(contentType, httpx.ContentTypeFormURLEncoded) {
		httpx.WriteHeader(request, httpx.HeaderContentType, this.override)
	}

	this.inner.ServeHTTP(response, request)
}
