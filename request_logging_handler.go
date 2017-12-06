package httpx

import (
	"net/http"

	"github.com/smartystreets/clock"
)

type RequestLoggingHandler struct {
	clock  *clock.Clock
	inner  http.Handler
	logger logger
}

func NewRequestLoggingHandler(inner http.Handler) *RequestLoggingHandler {
	return &RequestLoggingHandler{inner: inner, logger: new(contextLogger)}
}

func (this *RequestLoggingHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context := newContext(this.clock.UTCNow(), request, response)
	defer this.log(context)
	this.forward(context, request)
}

func (this *RequestLoggingHandler) forward(response http.ResponseWriter, request *http.Request) {
	if this.inner != nil {
		this.inner.ServeHTTP(response, request)
	}
}

func (this *RequestLoggingHandler) log(context *loggingContext) {
	err := recover()

	if context.canLogRequest() {
		this.logger.Log(context, err)
	}

	if err != nil {
		panic(err)
	}
}
