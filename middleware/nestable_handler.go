package middleware

import "net/http"

type NestableHandler struct {
	wrapped http.Handler
}

func NewNestableHandler(wrapped http.Handler) *NestableHandler {
	return &NestableHandler{wrapped: wrapped}
}

func (this *NestableHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.wrapped.ServeHTTP(response, request)
}

func (this *NestableHandler) Install(inner http.Handler) {
	panic("End of line.")
}
