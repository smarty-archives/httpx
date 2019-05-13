package middleware

import "net/http"

type NestableHandlerFunc struct {
	wrapped func() http.Handler
}

func NewNestableHandlerFunc(wrapped func() http.Handler) *NestableHandlerFunc {
	return &NestableHandlerFunc{wrapped: wrapped}
}

func (this *NestableHandlerFunc) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.wrapped().ServeHTTP(response, request)
}

func (this *NestableHandlerFunc) Install(inner http.Handler) {
	panic("End of line.")
}
