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
	if this.handleBody(response, request) {
		request.ParseForm()
		this.inner.ServeHTTP(response, request)
	} else {
		httpx.WriteResult(response, http.StatusRequestEntityTooLarge)
	}
}

func (this *PayloadLimitHandler) handleBody(response http.ResponseWriter, request *http.Request) bool {
	body := request.Body
	valid := this.validBody(response, request)
	body.Close()
	return valid
}

func (this *PayloadLimitHandler) validBody(response http.ResponseWriter, request *http.Request) bool {
	switch request.Method {
	case "PUT", "POST", "PATCH":
		return this.bufferedBody(response, request)
	default:
		return this.assertEmptyBody(request.Body)
	}
}
func (this *PayloadLimitHandler) bufferedBody(response http.ResponseWriter, request *http.Request) bool {
	if payload, err := ioutil.ReadAll(http.MaxBytesReader(response, request.Body, int64(this.maxSize))); err != nil {
		return false
	} else if len(payload) > 0 {
		request.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
	}

	return true
}
func (this *PayloadLimitHandler) assertEmptyBody(body io.ReadCloser) bool {
	read, _ := body.Read(discardBuffer)
	return read == 0 // this method shouldn't have a body, if read > 0, they provided a body and shouldn't have
}

var discardBuffer = make([]byte, 4) // writes to buffer are not thread safe but value is never read
