package middleware

import (
	"net/http"

	"github.com/smartystreets/httpx"
)

type ContextInitializationHandler struct {
	inner http.Handler
}

func NewContextInitializationHandler() *ContextInitializationHandler {
	return &ContextInitializationHandler{inner: NoopHandler{}}
}

func (this *ContextInitializationHandler) Install(inner http.Handler) {
	this.inner = inner
}

func (this *ContextInitializationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.inner.ServeHTTP(response, httpx.InitializeContext(request))
}
