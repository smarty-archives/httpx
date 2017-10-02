package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func (this *ContextFixture) TestNamespaceAccessBeforeInitialization_ShouldPanic() {
	uninitialized := httptest.NewRequest("GET", "/", nil)
	this.So(func() { Context(uninitialized)["KEY"] = 1 }, should.Panic)
	this.So(func() { Context(uninitialized).Int("KEY") }, should.Panic)
	this.So(func() { Context(uninitialized).Int64("KEY") }, should.Panic)
	this.So(func() { Context(uninitialized).Uint64("KEY") }, should.Panic)
	this.So(func() { Context(uninitialized).String("KEY") }, should.Panic)
}

func (this *ContextFixture) TestContextGetValue_WhenNotThere_ProvideDefaultValue() {
	this.So(Context(this.request).Int("KEY"), should.Equal, 0)
	this.So(Context(this.request).Int64("KEY"), should.Equal, 0)
	this.So(Context(this.request).Uint64("KEY"), should.Equal, 0)
	this.So(Context(this.request).String("KEY"), should.Equal, "")
}

func (this *ContextFixture) TestAfterInitialization_ReturnsStoredValues() {
	Context(this.request)["INTERFACE"] = struct{}{}
	Context(this.request)["INT"] = int(1)
	Context(this.request)["INT64"] = int64(1)
	Context(this.request)["UINT64"] = uint64(1)
	Context(this.request)["STRING"] = "hi"

	this.So(Context(this.request)["INTERFACE"], should.Resemble, struct{}{})
	this.So(Context(this.request).Int("INT"), should.Equal, int(1))
	this.So(Context(this.request).Int64("INT64"), should.Equal, int64(1))
	this.So(Context(this.request).Uint64("UINT64"), should.Equal, uint64(1))
	this.So(Context(this.request).String("STRING"), should.Equal, "hi")
}
