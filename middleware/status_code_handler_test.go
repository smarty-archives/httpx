package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStatusCodeHandlerFixture(t *testing.T) {
	gunit.Run(new(StatusCodeHandlerFixture), t)
}

type StatusCodeHandlerFixture struct {
	*gunit.Fixture
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *StatusCodeHandlerFixture) Setup() {
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
}

///////////////////////////////////////////////////////////////////////////////

func (this *StatusCodeHandlerFixture) TestResultContainsCorrectStatus() {
	this.assertResponse(404, "Not Found")
	this.assertResponse(429, "Too Many Requests")
}

func (this *StatusCodeHandlerFixture) assertResponse(statusCode int, statusText string) {
	this.assertHandlerResponse(NewStatusCodeHandler(statusCode), statusCode, statusText)
}
func (this *StatusCodeHandlerFixture) assertHandlerResponse(handler *StatusCodeHandler, statusCode int, statusText string) {
	this.Setup()

	handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, statusCode)
	this.So(this.response.Body.String(), should.EqualTrimSpace, statusText)
}

///////////////////////////////////////////////////////////////////////////////

func (this *StatusCodeHandlerFixture) TestResultContainsCorrectStatusCodeAndText() {
	this.assertHandlerResponse(NewStatusCodeAndTextHandler(200, "Hello, World 200!"), 200, "Hello, World 200!")
	this.assertHandlerResponse(NewStatusCodeAndTextHandler(599, "Hello, World 599!"), 599, "Hello, World 599!")
	this.assertHandlerResponse(NewNotFoundHandler(), 404, "Not Found")
}
