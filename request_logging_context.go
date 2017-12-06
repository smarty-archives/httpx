package httpx

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/smartystreets/logging"
)

type logger interface {
	Log(*loggingContext, interface{})
}

type contextLogger struct {
	logger *logging.Logger
}

func (this *contextLogger) Log(context *loggingContext, err interface{}) {
	this.logger.Print(context.Format(err))
}

type loggingContext struct {
	response      http.ResponseWriter
	request       *http.Request
	originalURL   string
	remoteAddress string
	started       time.Time
	statusCode    int
	bytesWritten  int
}

func newContext(now time.Time, request *http.Request, response http.ResponseWriter) *loggingContext {
	return &loggingContext{
		started:     now,
		response:    response,
		request:     request,
		originalURL: request.URL.String(),
		remoteAddress: orDefault(
			request.Header.Get(HeaderRemoteAddress),
			host(request.RemoteAddr)),
	}
}
func host(address string) string {
	host, _, _ := net.SplitHostPort(address)
	return host
}

func (this *loggingContext) Header() http.Header {
	return this.response.Header()
}
func (this *loggingContext) Write(p []byte) (n int, err error) {
	n, err = this.response.Write(p)
	this.bytesWritten += n
	return n, err
}
func (this *loggingContext) WriteHeader(status int) {
	if this.statusCode == 0 && status != 0 {
		this.statusCode = status
		this.response.WriteHeader(status)
	}
}
func (this *loggingContext) canLogRequest() bool {
	return this.originalURL != "/status"
}

func (this *loggingContext) logFields(err interface{}) []interface{} {
	return []interface{}{
		this.remoteAddress,
		this.started.Format("02/Jan/2006:15:04:05"),
		this.request.Method,
		this.originalURL,
		this.request.Proto,
		this.finalStatusCode(err),
		this.bytesWritten,
		orDefault(this.request.Referer(), "-"),
		orDefault(this.request.UserAgent(), "-"),
	}
}

func (this *loggingContext) Format(err interface{}) string {
	return fmt.Sprintf(`%s - - [%s +0000] "%s %s %s" %d %d "%s" "%s"`, this.logFields(err)...)
}
func (this *loggingContext) finalStatusCode(err interface{}) int {
	if err != nil {
		return http.StatusInternalServerError
	} else if this.statusCode == 0 {
		return http.StatusOK
	} else {
		return this.statusCode
	}
}
func orDefault(value, Default string) string {
	if len(value) > 0 {
		return value
	} else {
		return Default
	}
}
