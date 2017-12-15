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
	contentType := httpx.ReadHeader(request, "Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		httpx.WriteHeader(request, "Content-Type", this.override)
	}

	this.inner.ServeHTTP(response, request)
}
