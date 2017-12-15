package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestNestableHandlerFixture(t *testing.T) {
	gunit.Run(new(NestableHandlerFixture), t)
}

type nestingHandler interface {
	http.Handler
	Install(http.Handler)
}

type NestableHandlerFixture struct {
	*gunit.Fixture

	wrapped http.Handler
	outer   nestingHandler
}

func (this *NestableHandlerFixture) Setup() {
	this.wrapped = &NonNestingHandler{}
	this.outer = NewNestableHandler(this.wrapped)
}

func (this *NestableHandlerFixture) TestInvocationOfInstallPanics() {
	this.So(func() { this.outer.Install(this.wrapped) }, should.Panic)
}

func (this *NestableHandlerFixture) TestWrappedHandlerIsInvoked() {
	response := httptest.NewRecorder()
	this.outer.ServeHTTP(response, nil)
	this.So(response.Body.String(), should.Equal, "1")
}

///////////////////////////////////////////////////////////////////////////////

type NonNestingHandler struct{}

func (this *NonNestingHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "1")
}

////////////////////////////////////////////////////////////////////////////////
