package middleware

import "net/http"

type FakeInnerHandler struct {
	calls        int
	request      *http.Request
	response     http.ResponseWriter
	responseBody string
}

func NewFakeInnerHandler(responseBody string) *FakeInnerHandler {
	return &FakeInnerHandler{responseBody: responseBody}
}

func (this *FakeInnerHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.calls++
	this.response = response
	this.request = request
	if len(this.responseBody) > 0 {
		response.Write([]byte(this.responseBody))
	}
}
