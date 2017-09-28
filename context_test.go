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

func (this *ContextFixture) TestContextInitialization() {
	actual := this.request.Context().Value(contextNamespace)
	expected := map[interface{}]interface{}{}
	this.So(actual, should.Resemble, expected)
}

func (this *ContextFixture) TestContextSetValue_ThenAvailableOnContextMap() {
	Contextual(*this.request).Set("KEY", "VALUE")
	values := this.request.Context().Value(contextNamespace).(map[interface{}]interface{})
	this.So(values["KEY"], should.Equal, "VALUE")
}

func (this *ContextFixture) TestContextGetValue_WhenNotThere_ProvideDefaultValue() {
	values := Contextual(*this.request)

	this.So(values.Int("KEY"), should.Equal, 0)
	this.So(values.Int64("KEY"), should.Equal, 0)
	this.So(values.Uint64("KEY"), should.Equal, 0)
	this.So(values.String("KEY"), should.Equal, "")
}

func (this *ContextFixture) TestContextSetValue_Int() {
	values := Contextual(*this.request)
	values.Set("INT", int(1))
	this.So(values.Int("INT"), should.Equal, int(1))
}

func (this *ContextFixture) TestContextSetValue_Int64() {
	values := Contextual(*this.request)
	values.Set("INT64", int64(1))
	this.So(values.Int64("INT64"), should.Equal, int64(1))
}

func (this *ContextFixture) TestContextSetValue_Uint64() {
	values := Contextual(*this.request)
	values.Set("UINT64", uint64(1))
	this.So(values.Uint64("UINT64"), should.Equal, uint64(1))
}

func (this *ContextFixture) TestContextSetValue_string() {
	values := Contextual(*this.request)
	values.Set("STRING", "hi")
	this.So(values.String("STRING"), should.Equal, "hi")
}
