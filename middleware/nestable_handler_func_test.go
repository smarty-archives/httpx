package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestNestableHandlerFuncFixture(t *testing.T) {
	gunit.Run(new(NestableHandlerFuncFixture), t)
}

type NestableHandlerFuncFixture struct {
	*gunit.Fixture

	wrapped http.Handler
	outer   nestingHandler
}

func (this *NestableHandlerFuncFixture) Setup() {
	this.wrapped = &NonNestingHandler{}
	this.outer = NewNestableHandlerFunc(func() http.Handler { return this.wrapped })
}

func (this *NestableHandlerFuncFixture) TestInvocationOfInstallPanics() {
	this.So(func() { this.outer.Install(this.wrapped) }, should.Panic)
}

func (this *NestableHandlerFuncFixture) TestWrappedHandlerIsInvoked() {
	response := httptest.NewRecorder()
	this.outer.ServeHTTP(response, nil)
	this.So(response.Body.String(), should.Equal, "1")
}
