package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestContentTypeOverrideHandlerFixture(t *testing.T) {
	gunit.Run(new(ContentTypeOverrideHandlerFixture), t)
}

type ContentTypeOverrideHandlerFixture struct {
	*gunit.Fixture

	handler  *ContentTypeOverrideHandler
	fake     *FakeInnerHandler
	response *httptest.ResponseRecorder
	request  *http.Request
}

func (this *ContentTypeOverrideHandlerFixture) Setup() {
	this.fake = &FakeInnerHandler{}
	this.handler = NewContentTypeOverrideHandler("application/json")
	this.handler.Install(this.fake)
	this.response = httptest.NewRecorder()
	this.request, _ = http.NewRequest("POST", "/", nil)
}

func (this *ContentTypeOverrideHandlerFixture) TestFormContentTypeOverridden() {
	this.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(this.fake.calls, should.Equal, 1)
	this.So(this.fake.request, should.Equal, this.request)
	this.So(this.fake.response, should.Equal, this.response)
}

func (this *ContentTypeOverrideHandlerFixture) TestStartsWithFormContentTypeOverridden() {
	this.request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF8")

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(this.fake.calls, should.Equal, 1)
	this.So(this.fake.request, should.Equal, this.request)
	this.So(this.fake.response, should.Equal, this.response)
}

func (this *ContentTypeOverrideHandlerFixture) TestContainsFormContentTypeOverridden() {
	this.request.Header.Set("Content-Type", "random; application/x-www-form-urlencoded; charset=UTF8")

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(this.fake.calls, should.Equal, 1)
	this.So(this.fake.request, should.Equal, this.request)
	this.So(this.fake.response, should.Equal, this.response)
}

func (this *ContentTypeOverrideHandlerFixture) TestOtherContentTypesUnchanged() {
	this.request.Header.Set("Content-Type", "application/octet-stream")

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.request.Header.Get("Content-Type"), should.Equal, "application/octet-stream")
	this.So(this.fake.calls, should.Equal, 1)
	this.So(this.fake.request, should.Equal, this.request)
	this.So(this.fake.response, should.Equal, this.response)
}
