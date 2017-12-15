package middleware

import (
	"net/http"

	"github.com/smartystreets/clock"
	"github.com/smartystreets/httpx"
)

type RequestLoggingHandler struct {
	clock  *clock.Clock
	inner  http.Handler
	logger logger

	remoteAddressHeader string
}

func NewRequestLoggingHandler(inner http.Handler, remoteAddressHeader string) *RequestLoggingHandler {
	return &RequestLoggingHandler{
		inner:  inner,
		logger: new(contextLogger),

		remoteAddressHeader: remoteAddressHeader}
}

func (this *RequestLoggingHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	remoteAddress := httpx.ReadClientIPAddress(request, this.remoteAddressHeader)
	context := newContext(this.clock.UTCNow(), remoteAddress, request, response)
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
