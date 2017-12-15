package middleware

import "net/http"

type StatusCodeHandler struct {
	statusCode int
	statusText string
}

func NewNotFoundHandler() *StatusCodeHandler {
	return NewStatusCodeHandler(http.StatusNotFound)
}

func NewStatusCodeHandler(statusCode int) *StatusCodeHandler {
	return NewStatusCodeAndTextHandler(statusCode, http.StatusText(statusCode))
}
func NewStatusCodeAndTextHandler(statusCode int, statusText string) *StatusCodeHandler {
	return &StatusCodeHandler{statusCode: statusCode, statusText: statusText}
}

func (this *StatusCodeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	http.Error(response, this.statusText, this.statusCode)
}
