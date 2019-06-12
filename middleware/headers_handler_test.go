package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/smartystreets/httpx"
)

func TestHeadersHandlerFixture(t *testing.T) {
	gunit.Run(new(HeadersHandlerFixture), t)
}

type HeadersHandlerFixture struct {
	*gunit.Fixture

	handler *HeadersHandler
	inner   *FakeInnerHandler

	headers map[string]string

	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *HeadersHandlerFixture) Setup() {
	this.reset()
}
func (this *HeadersHandlerFixture) reset() {
	this.headers = map[string]string{}
	this.inner = &FakeInnerHandler{}
	this.handler = NewHeadersHandler(this.headers)
	this.handler.Install(this.inner)
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
}

func (this *HeadersHandlerFixture) TestHeadersAdded() {
	this.headers["key1"] = "value1"
	this.headers["key2"] = "value2"

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("key1"), should.Equal, "value1")
	this.So(this.response.Header().Get("key2"), should.Equal, "value2")
}

func (this *HeadersHandlerFixture) TestInnerHandlerCalled() {
	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.inner.calls, should.Equal, 1)
	this.So(this.inner.request, should.Equal, this.request)
	this.So(this.inner.response, should.Equal, this.response)
	this.So(this.response.Code, should.Equal, http.StatusOK)
}

func (this *HeadersHandlerFixture) TestNilInnerHandler_NOT_Called() {
	this.handler = NewHeadersHandler(this.headers)

	this.headers["key1"] = "value1"
	this.headers["key2"] = "value2"

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("key1"), should.Equal, "value1")
	this.So(this.response.Header().Get("key2"), should.Equal, "value2")
	this.So(this.inner.calls, should.Equal, 0)
}

func (this *HeadersHandlerFixture) TestCORSHeaders() {
	this.headers["cors-header-1"] = "value1"
	this.headers["cors-header-2"] = "value2"
	this.handler = CORSHeadersHandler(this.headers)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("cors-header-1"), should.Equal, "value1")
	this.So(this.response.Header().Get("cors-header-2"), should.Equal, "value2")
	this.So(this.inner.calls, should.Equal, 0)
}

func (this *HeadersHandlerFixture) TestBrowserHeadersNotWrittenWhenReferrerAndOriginAreAbsent() {
	this.headers["key1"] = "value1"
	this.handler = BrowserHeadersHandler(this.headers)
	this.handler.Install(this.inner)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("key1"), should.BeEmpty)
	this.So(this.inner.calls, should.Equal, 1)
}

func (this *HeadersHandlerFixture) TestBrowserHeadersWrittenWhenReferrerPresent() {
	this.headers["key1"] = "value1"
	httpx.WriteHeader(this.request, "Referer", "some-value")
	this.handler = BrowserHeadersHandler(this.headers)
	this.handler.Install(this.inner)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("key1"), should.Equal, "value1")
	this.So(this.inner.calls, should.Equal, 1)
}

func (this *HeadersHandlerFixture) TestBrowserHeadersWrittenWhenOriginPresent() {
	this.headers["key1"] = "value1"
	httpx.WriteHeader(this.request, "Origin", "some-value")
	this.handler = BrowserHeadersHandler(this.headers)
	this.handler.Install(this.inner)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("key1"), should.Equal, "value1")
	this.So(this.inner.calls, should.Equal, 1)
}

func (this *HeadersHandlerFixture) TestDefaultCORSHeadersWritten() {
	this.response.Header().Set("Access-Control-Allow-Origin", "overwritten")
	this.handler = DefaultCORSHeadersHandler()

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("Access-Control-Allow-Origin"), should.Equal, "*")
	this.So(this.response.Header().Get("Access-Control-Allow-Methods"), should.Equal, "GET, PUT, POST, DELETE, HEAD")
	this.So(this.response.Header().Get("Access-Control-Allow-Headers"), should.Equal, "Accept, Content-Type, Content-Length, Referer, Origin, Host")
	this.So(this.response.Header().Get("Access-Control-Allow-Credentials"), should.Equal, "true")
	this.So(this.response.Header().Get("Access-Control-Max-Age"), should.Equal, "600")
}

func (this *HeadersHandlerFixture) TestOnlyPopulatedHeadersAreWritten() {
	this.response.Header().Set("Access-Control-Allow-Origin", "preserved")
	this.headers["Access-Control-Allow-Origin"] = ""
	this.handler = CORSHeadersHandler(this.headers)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("Access-Control-Allow-Origin"), should.Equal, "preserved")
}
