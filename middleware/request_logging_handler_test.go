package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRequestLoggingHandlerFixture(t *testing.T) {
	gunit.Run(new(RequestLoggingHandlerFixture), t)
}

type RequestLoggingHandlerFixture struct {
	*gunit.Fixture

	handler  *RequestLoggingHandler
	response *httptest.ResponseRecorder
	request  *http.Request
	inner    *FakeHandler
	now      time.Time
	clock    func() time.Time

	loggedContext  *loggingContext
	loggedPanicErr interface{}
}

func (this *RequestLoggingHandlerFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/not-status", nil)
	this.response = httptest.NewRecorder()
	this.inner = NewFakeHandler()
	this.now = time.Now()
	this.clock = func() time.Time { return this.now }
	this.handler = NewRequestLoggingHandler(this.inner, "X-Remote-Address", this.clock)
	this.handler.logger = this
}

func (this *RequestLoggingHandlerFixture) Log(context *loggingContext, panicErr interface{}) {
	this.loggedContext = context
	this.loggedPanicErr = panicErr
}

func (this *RequestLoggingHandlerFixture) TestStatusRequestsNotLogged() {
	this.inner.statusCode = http.StatusTeapot
	this.request.URL.Path = "/status"
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext, should.BeNil)
	this.So(this.loggedPanicErr, should.BeNil)
	this.So(this.response.Code, should.Equal, http.StatusTeapot)
}

func (this *RequestLoggingHandlerFixture) TestRequestLogged() {
	this.inner.statusCode = http.StatusTeapot
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext.request, should.Equal, this.request)
	this.So(this.loggedContext.response, should.Equal, this.response)
	this.So(this.loggedContext.statusCode, should.Equal, http.StatusTeapot)
	this.So(this.loggedContext.started, should.Equal, this.now)
	this.So(this.loggedPanicErr, should.BeNil)
	this.So(this.response.Code, should.Equal, http.StatusTeapot)
}

func (this *RequestLoggingHandlerFixture) TestNilInnerHandler() {
	this.handler.inner = nil
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext.request, should.Equal, this.request)
	this.So(this.loggedContext.response, should.Equal, this.response)
	this.So(this.loggedContext.statusCode, should.Equal, 0)
	this.So(this.loggedContext.started, should.Equal, this.now)
	this.So(this.loggedPanicErr, should.BeNil)
}

func (this *RequestLoggingHandlerFixture) TestInnerHandlerPanics() {
	this.inner.panicMessage = "GOPHERS!"
	this.So(func() { this.handler.ServeHTTP(this.response, this.request) }, should.PanicWith, this.inner.panicMessage)
	this.So(this.loggedPanicErr, should.Equal, this.inner.panicMessage)
}

func (this *RequestLoggingHandlerFixture) TestMultipleStatusCodeWrites_FirstWins() {
	this.handler.inner = new(FakeHandler_MultipleStatusCodes)
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext.statusCode, should.Equal, 200)
	this.So(this.response.Code, should.Equal, http.StatusOK)
}

func (this *RequestLoggingHandlerFixture) TestResponseBodyCounted() {
	this.inner.response = strings.Repeat(" ", 100)
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext.bytesWritten, should.Equal, 100)
}

func (this *RequestLoggingHandlerFixture) TestResponseHeadersWritten() {
	this.inner.headers["Hello"] = "World"
	this.handler.ServeHTTP(this.response, this.request)
	this.So(this.loggedContext.response.Header().Get("Hello"), should.Equal, "World")
}

/////////////////////////////////////////////////////////////////////////////////

type FakeHandler struct {
	request      *http.Request
	statusCode   int
	response     string
	panicMessage string
	headers      map[string]string
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{headers: make(map[string]string)}
}

func (this *FakeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if this.panicMessage != "" {
		panic(this.panicMessage)
	}

	this.request = request

	for key, value := range this.headers {
		response.Header().Set(key, value)
	}
	response.WriteHeader(this.statusCode)
	response.Write([]byte(this.response))
}

type FakeHandler_MultipleStatusCodes struct{}

func (this *FakeHandler_MultipleStatusCodes) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	for x := 200; x < 300; x++ {
		response.WriteHeader(x) // only the first should stick
	}
}
