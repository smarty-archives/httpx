package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestAuthenticatedCORSHandlerFixture(t *testing.T) {
	gunit.Run(new(AuthenticatedCORSHandlerFixture), t)
}

type AuthenticatedCORSHandlerFixture struct {
	*gunit.Fixture

	handler  *AuthenticatedCORSHandler
	fake     *FakeInnerHandler
	response *httptest.ResponseRecorder
	request  *http.Request
}

func (this *AuthenticatedCORSHandlerFixture) Setup() {
	this.fake = &FakeInnerHandler{}
	this.handler = NewAuthenticatedCORSHandler("localhost", "smartystreets.dev", "smartystreets.com")
	this.handler.Install(this.fake)
	this.response = httptest.NewRecorder()
	this.request, _ = http.NewRequest("GET", "/", nil)
}

func (this *AuthenticatedCORSHandlerFixture) TestInnerAlwaysCalled() {
	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.fake.response, should.Equal, this.response)
	this.So(this.fake.request, should.Equal, this.request)
}

func (this *AuthenticatedCORSHandlerFixture) TestWriteHeadersOnAllowedOrigin() {
	this.assertCORSHeaders("http://localhost:8080")
	this.assertCORSHeaders("http://smartystreets.dev:5678")
	this.assertCORSHeaders("http://smartystreets.com:1234")
	this.assertCORSHeaders("https://smartystreets.com")
}
func (this *AuthenticatedCORSHandlerFixture) assertCORSHeaders(origin string) {
	this.response = httptest.NewRecorder()
	this.request.Header.Set("Origin", origin)

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header().Get("Access-Control-Allow-Credentials"), should.Equal, "true")
	this.So(this.response.Header().Get("Access-Control-Allow-Origin"), should.Equal, origin)
}

func (this *AuthenticatedCORSHandlerFixture) TestNoHeadersOnOtherOrigins() {
	this.request.Header.Set("Origin", "https://whatever.com:1234")

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Header(), should.NotContainKey, "Access-Control-Allow-Credentials")
	this.So(this.response.Header(), should.NotContainKey, "Access-Control-Allow-Origin")
}
