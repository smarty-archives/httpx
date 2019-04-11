package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/smartystreets/httpx"
)

type PayloadLimitHandler struct {
	inner   http.Handler
	maxSize uint64
}

func NewPayloadLimitHandler(maxSize uint64) *PayloadLimitHandler {
	return &PayloadLimitHandler{maxSize: maxSize, inner: NoopHandler{}}
}
func (this *PayloadLimitHandler) Install(inner http.Handler) {
	this.inner = inner
}

func (this *PayloadLimitHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if methodAllowsBody[request.Method] {
		this.handleAllowedBody(response, request)
	} else {
		this.handleDisallowedBody(response, request)
	}
}

func (this *PayloadLimitHandler) handleAllowedBody(response http.ResponseWriter, request *http.Request) {
	if contents, err := this.readBody(response, request); err != nil {
		httpx.WriteResult(response, http.StatusRequestEntityTooLarge)
	} else {
		_ = request.Body.Close()
		request.Body = ioutil.NopCloser(bytes.NewReader(contents))
		this.inner.ServeHTTP(response, request)
	}
}
func (this *PayloadLimitHandler) readBody(response http.ResponseWriter, request *http.Request) ([]byte, error) {
	return ioutil.ReadAll(http.MaxBytesReader(response, request.Body, int64(this.maxSize)))
}

func (this *PayloadLimitHandler) handleDisallowedBody(response http.ResponseWriter, request *http.Request) {
	if !this.bodyIsEmpty(request.Body) {
		httpx.WriteResult(response, http.StatusRequestEntityTooLarge)
	} else {
		this.inner.ServeHTTP(response, request)
	}
}
func (this *PayloadLimitHandler) bodyIsEmpty(body io.ReadCloser) bool {
	read, _ := body.Read(discardBuffer)
	return read == 0
}

var discardBuffer = make([]byte, 4) // writes to buffer are not thread safe but value is never read

var methodAllowsBody = map[string]bool{
	"PUT":   true, // has a body
	"POST":  true, // has a body
	"PATCH": true, // body is allowed

	"GET":     false,
	"HEAD":    false,
	"OPTIONS": false,
	"DELETE":  false,
	"TRACE":   false,
	"DEBUG":   false,
}
