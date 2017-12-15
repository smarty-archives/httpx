package middleware

import "net/http"

type NoopHandler struct{}

func (this NoopHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}
