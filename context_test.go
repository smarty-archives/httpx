package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestContextFixture(t *testing.T) {
	gunit.Run(new(ContextFixture), t)
}

type ContextFixture struct {
	*gunit.Fixture
	request *http.Request
}

func (this *ContextFixture) Setup() {
	this.request = InitializeContext(httptest.NewRequest("GET", "/", nil))
}

func (this *ContextFixture) TestWhatHappensWhenSetIsCalledBeforeInitialization() {
	this.So(func() { SetContext(httptest.NewRequest("GET", "/", nil), "key", "value") }, should.Panic)
}

func (this *ContextFixture) TestWhatHappensWhenGetsAreCalledBeforeInitialization() {
	this.So(func() { GetContext(httptest.NewRequest("GET", "/", nil), "key") }, should.Panic)
	this.So(func() { GetContextInt(httptest.NewRequest("GET", "/", nil), "key") }, should.Panic)
	this.So(func() { GetContextInt64(httptest.NewRequest("GET", "/", nil), "key") }, should.Panic)
	this.So(func() { GetContextUint64(httptest.NewRequest("GET", "/", nil), "key") }, should.Panic)
	this.So(func() { GetContextString(httptest.NewRequest("GET", "/", nil), "key") }, should.Panic)
}

func (this *ContextFixture) TestContextInitialization() {
	actual := this.request.Context().Value(contextNamespace)
	expected := make(map[interface{}]interface{})
	this.So(actual, should.Resemble, expected)
}

func (this *ContextFixture) TestContextSetValue_ThenAvailableOnContextMap() {
	SetContext(this.request, "KEY", "VALUE")
	values := this.request.Context().Value(contextNamespace).(map[interface{}]interface{})
	this.So(values["KEY"], should.Equal, "VALUE")
}

func (this *ContextFixture) TestContextGetValue_WhenNotThere_ProvideDefaultValue() {
	this.So(GetContext(this.request, "KEY"), should.BeNil)
	this.So(GetContextInt(this.request, "KEY"), should.Equal, int(0))
	this.So(GetContextInt64(this.request, "KEY"), should.Equal, int64(0))
	this.So(GetContextUint64(this.request, "KEY"), should.Equal, uint64(0))
	this.So(GetContextString(this.request, "KEY"), should.Equal, "")
}

func (this *ContextFixture) TestGetContext() {
	now := time.Now()
	SetContext(this.request, "time", now)
	this.So(GetContext(this.request, "time"), should.Equal, now)
}

func (this *ContextFixture) TestGetContext_Int() {
	SetContext(this.request, "INT", int(1))
	this.So(GetContextInt(this.request, "INT"), should.Equal, int(1))
}

func (this *ContextFixture) TestGetContext_Int64() {
	SetContext(this.request, "INT64", int64(1))
	this.So(GetContextInt64(this.request, "INT64"), should.Equal, int64(1))
}

func (this *ContextFixture) TestGetContext_Uint64() {
	SetContext(this.request, "UINT64", uint64(1))
	this.So(GetContextUint64(this.request, "UINT64"), should.Equal, uint64(1))
}

func (this *ContextFixture) TestGetContext_string() {
	SetContext(this.request, "STRING", string("hi"))
	this.So(GetContextString(this.request, "STRING"), should.Equal, "hi")
}
