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

func (this *ContextFixture) TestContextSetValue() {
	Contextual(*this.request).Set("KEY", "VALUE")
	values := this.request.Context().Value(contextNamespace).(map[interface{}]interface{})
	this.So(values["KEY"], should.Equal, "VALUE")
}
