package middleware

import (
	"net/http"
	"time"

	"github.com/smartystreets/httpx/v2"
)

// Deprecated: remove when KAFKA is in place.
type RequestLoggingHandler struct {
	now    func() time.Time
	inner  http.Handler
	logger logger

	remoteAddressHeader string
}

// Deprecated: remove when KAFKA is in place.
func NewRequestLoggingHandler(inner http.Handler, remoteAddressHeader string, now func() time.Time) *RequestLoggingHandler {
	return &RequestLoggingHandler{
		inner:  inner,
		logger: new(contextLogger),
		now: now,
		remoteAddressHeader: remoteAddressHeader}
}

func (this *RequestLoggingHandler) Install(inner http.Handler) { this.inner = inner }

func (this *RequestLoggingHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	remoteAddress := httpx.ReadClientIPAddress(request, this.remoteAddressHeader)
	context := newContext(this.now(), remoteAddress, request, response)
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

func utcNow() time.Time {
	return time.Now().UTC()
}