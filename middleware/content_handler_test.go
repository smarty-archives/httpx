package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestContentHandlerFixture(t *testing.T) {
	gunit.Run(new(ContentHandlerFixture), t)
}

type ContentHandlerFixture struct {
	*gunit.Fixture
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *ContentHandlerFixture) Setup() {
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
}

func (this *ContentHandlerFixture) TestByteResponseHasCorrectContent() {
	handler := NewContentHandler([]byte("Hello, World!"), "text/plain")

	handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Body.String(), should.Equal, "Hello, World!")
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "text/plain")
}

func (this *ContentHandlerFixture) TestByteResponseHasContentType() {
	handler := NewContentHandler([]byte("Hello, World!"), "")

	handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Body.String(), should.Equal, "Hello, World!")
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/octet-stream")
}

func (this *ContentHandlerFixture) TestStringResponseHasCorrectContent() {
	handler := NewContentStringHandler("Hello, World!", "text/plain")

	handler.ServeHTTP(this.response, this.request)

	this.So(this.response.Body.String(), should.Equal, "Hello, World!")
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "text/plain")
}
