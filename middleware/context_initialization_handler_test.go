package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/httpx"
)

func TestContextInitializationHandlerFixture(t *testing.T) {
	gunit.Run(new(ContextInitializationHandlerFixture), t)
}

type ContextInitializationHandlerFixture struct {
	*gunit.Fixture

	handler *ContextInitializationHandler
	fake    *FakeInnerHandler

	response *httptest.ResponseRecorder
	request  *http.Request
}

func (this *ContextInitializationHandlerFixture) Setup() {
	this.fake = NewFakeInnerHandler("")
	this.handler = NewContextInitializationHandler()
	this.handler.Install(this.fake)
	this.response = httptest.NewRecorder()
	this.request = httptest.NewRequest("GET", "/42", nil)
}

func (this *ContextInitializationHandlerFixture) TestConstructor() {
	this.handler.ServeHTTP(this.response, this.request)

	this.assertRequestContextInitialized()
	this.assertResponseAndRequestPassedToInnerHandler()
}
func (this *ContextInitializationHandlerFixture) assertResponseAndRequestPassedToInnerHandler() {
	this.So(this.fake.response, should.Resemble, this.response)
	this.So(this.fake.request.URL.String(), should.Equal, "/42")
}
func (this *ContextInitializationHandlerFixture) assertRequestContextInitialized() {
	this.So(func() { httpx.Context(this.fake.request) }, should.NotPanic)
}
