package httpx

import (
	"net/http/httptest"
	"testing"

	"net/http"

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
	this.request = httptest.NewRequest("GET", "/", nil)
}

func (this *ContextFixture) TestContextInitialization() {
	request := InitializeContext(this.request)
	actual := request.Context().Value(contextNamespace)
	expected := map[interface{}]interface{}{}
	this.So(actual, should.Resemble, expected)
}
