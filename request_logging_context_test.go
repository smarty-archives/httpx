package httpx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestLoggingContextFixture(t *testing.T) {
	gunit.Run(new(LoggingContextFixture), t)
}

type LoggingContextFixture struct {
	*gunit.Fixture

	now      time.Time
	request  *http.Request
	response *httptest.ResponseRecorder
	err      interface{}
	context  *loggingContext
}

func (this *LoggingContextFixture) Setup() {
	this.response = httptest.NewRecorder()
	this.request = httptest.NewRequest("GET", "/", nil)
	this.request.RemoteAddr = "0.0.0.0:0"
}

func (this *LoggingContextFixture) createContext() {
	this.context = newContext(this.now, this.request, this.response)
}
func (this *LoggingContextFixture) formatContext() []interface{} {
	return this.context.logFields(this.err)
}
func (this *LoggingContextFixture) string(i int) string   { return this.formatContext()[i].(string) }
func (this *LoggingContextFixture) int(i int) int         { return this.formatContext()[i].(int) }
func (this *LoggingContextFixture) remoteAddress() string { return this.string(0) }
func (this *LoggingContextFixture) timeStamp() string     { return this.string(1) }
func (this *LoggingContextFixture) requestMethod() string { return this.string(2) }
func (this *LoggingContextFixture) originalURL() string   { return this.string(3) }
func (this *LoggingContextFixture) requestProto() string  { return this.string(4) }
func (this *LoggingContextFixture) statusCode() int       { return this.int(5) }
func (this *LoggingContextFixture) bytesWritten() int     { return this.int(6) }
func (this *LoggingContextFixture) referringURL() string  { return this.string(7) }
func (this *LoggingContextFixture) userAgent() string     { return this.string(8) }

func (this *LoggingContextFixture) TestRemoteAddressLogged() {
	this.request.RemoteAddr = "1.2.3.4:5"
	this.createContext()
	this.So(this.remoteAddress(), should.Equal, `1.2.3.4`)
}

func (this *LoggingContextFixture) TestRemoteAddressFromLoadBalancerOverridesRemoteAddressOnRequest() {
	this.request.RemoteAddr = "1.2.3.4:5"
	this.request.Header.Set("X-Security-Remote-Address", "RemoteAddress")
	this.createContext()
	this.So(this.remoteAddress(), should.Equal, `RemoteAddress`)
}

func (this *LoggingContextFixture) TestStartTimeLogged() {
	this.now = time.Unix(0, 0)
	this.createContext()
	this.So(this.timeStamp(), should.Equal, `31/Dec/1969:17:00:00`)
}

func (this *LoggingContextFixture) TestRequestMethodLogged() {
	this.request.Method = "POST"
	this.createContext()
	this.So(this.requestMethod(), should.Equal, "POST")
}

func (this *LoggingContextFixture) TestRequestProtoLogged() {
	this.request.Proto = "PROTO"
	this.createContext()
	this.So(this.requestProto(), should.Equal, "PROTO")
}

func (this *LoggingContextFixture) TestFinalStatusCode_AfterPanic() {
	this.err = "PANIC!"
	this.createContext()
	this.So(this.statusCode(), should.Equal, http.StatusInternalServerError)
}

func (this *LoggingContextFixture) TestFinalStatusCode_SetPreviouslyByApplication() {
	this.createContext()
	this.context.WriteHeader(http.StatusTeapot)
	this.So(this.statusCode(), should.Equal, http.StatusTeapot)
}

func (this *LoggingContextFixture) TestFinalStatusCode_Unspecified() {
	this.createContext()
	this.So(this.statusCode(), should.Equal, http.StatusOK)
}

func (this *LoggingContextFixture) TestBytesWritten() {
	this.createContext()
	fmt.Fprint(this.context, "Hello")
	fmt.Fprint(this.context, "World")
	this.So(this.bytesWritten(), should.Equal, len("HelloWorld"))
}

func (this *LoggingContextFixture) TestReferringURL_Present() {
	this.request.Header.Set("Referer", "the-referer")
	this.createContext()
	this.So(this.referringURL(), should.Equal, "the-referer")
}

func (this *LoggingContextFixture) TestReferringURL_Absent() {
	this.request.Header.Set("Referer", "")
	this.createContext()
	this.So(this.referringURL(), should.Equal, "-")
}

func (this *LoggingContextFixture) TestUserAgent_Present() {
	this.request.Header.Set("User-Agent", "the-user-agent")
	this.createContext()
	this.So(this.userAgent(), should.Equal, "the-user-agent")
}

func (this *LoggingContextFixture) TestUserAgent_Absent() {
	this.request.Header.Set("User-Agent", "")
	this.createContext()
	this.So(this.userAgent(), should.Equal, "-")
}
