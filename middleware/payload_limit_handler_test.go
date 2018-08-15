package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPayloadLimitHandlerFixture(t *testing.T) {
	gunit.Run(new(PayloadLimitHandlerFixture), t)
}

type PayloadLimitHandlerFixture struct {
	*gunit.Fixture

	maxSize uint64
	handler *PayloadLimitHandler
	inner   *FakeInnerHandler

	bodyBuffer  *bytes.Buffer
	requestBody *ReadCloser
	request     *http.Request
	response    *httptest.ResponseRecorder
}

func (this *PayloadLimitHandlerFixture) Setup() {
	this.maxSize = 16
	this.inner = &FakeInnerHandler{}
	this.handler = NewPayloadLimitHandler(this.maxSize)
	this.handler.Install(this.inner)

	this.bodyBuffer = bytes.NewBuffer([]byte{})
	this.requestBody = NewReadCloser(this.bodyBuffer)
	this.request, _ = http.NewRequest("PUT", "/", this.requestBody)
	this.response = httptest.NewRecorder()
}

func (this *PayloadLimitHandlerFixture) TestLargeRequestsAreRejected() {
	this.bodyBuffer.WriteString(strings.Repeat(".", int(this.maxSize+1)))

	this.handler.ServeHTTP(this.response, this.request)

	this.assertRequestRejected()
}
func (this *PayloadLimitHandlerFixture) assertRequestRejected() {
	this.So(this.response.Code, should.Equal, http.StatusRequestEntityTooLarge)
	this.So(strings.TrimSpace(this.response.Body.String()), should.Equal, "Request Entity Too Large")
	this.So(this.inner.calls, should.Equal, 0)
	this.So(this.requestBody.closed, should.Equal, 1)
}

func (this *PayloadLimitHandlerFixture) TestSmallRequestsAreAccepted() {
	expectedContents := strings.Repeat(".", int(this.maxSize-1))
	this.bodyBuffer.WriteString(expectedContents)

	this.handler.ServeHTTP(this.response, this.request)

	actualContents, _ := ioutil.ReadAll(this.request.Body)
	this.So(string(actualContents), should.Equal, expectedContents)

	this.assertRequestAccepted()
}
func (this *PayloadLimitHandlerFixture) assertRequestAccepted() {
	this.So(this.inner.calls, should.Equal, 1)
	this.So(this.inner.request, should.Equal, this.request)
	this.So(this.inner.response, should.Equal, this.response)
	this.So(this.requestBody.closed, should.Equal, 1)
}

func (this *PayloadLimitHandlerFixture) TestBodylessMethodsWithBodyRejected() {
	for method, containsBody := range httpMethods {
		if !containsBody {
			this.assertBodylessMethodWithBodyRejected(method)
		}
	}
}
func (this *PayloadLimitHandlerFixture) assertBodylessMethodWithBodyRejected(method string) {
	this.Setup()
	this.bodyBuffer.WriteString("INVALID BODY")
	this.request.Method = method

	this.handler.ServeHTTP(this.response, this.request)

	this.assertRequestRejected()
}

func (this *PayloadLimitHandlerFixture) TestBodylessMethodsBodyUnmodified() {
	for method, containsBody := range httpMethods {
		if !containsBody {
			this.assertBodylessMethodsBodyUnmodified(method)
		}
	}
}
func (this *PayloadLimitHandlerFixture) assertBodylessMethodsBodyUnmodified(method string) {
	this.Setup()
	this.request.Method = method

	this.handler.ServeHTTP(this.response, this.request)

	this.So(this.inner.request.Body, should.Equal, this.requestBody)
	this.assertRequestAccepted()
}

func (this *PayloadLimitHandlerFixture) TestMethodsAllowedBodyCanContainBody() {
	for method, containsBody := range httpMethods {
		if containsBody {
			this.assertBodyMethodsAllowedThrough(method)
		}
	}
}
func (this *PayloadLimitHandlerFixture) assertBodyMethodsAllowedThrough(method string) {
	this.Setup()
	this.request.Method = method
	this.bodyBuffer.WriteString(strings.Repeat(".", int(this.maxSize-1)))

	this.handler.ServeHTTP(this.response, this.request)

	this.assertRequestAccepted()
}

type ReadCloser struct {
	reader io.Reader
	closed int
}

func NewReadCloser(reader io.Reader) *ReadCloser         { return &ReadCloser{reader: reader} }
func (this *ReadCloser) Read(buffer []byte) (int, error) { return this.reader.Read(buffer) }
func (this *ReadCloser) Close() error                    { this.closed++; return nil }

var httpMethods = map[string]bool{
	"GET":     false,
	"HEAD":    false,
	"OPTIONS": false,
	"PUT":     true, // has a body
	"POST":    true, // has a body
	"PATCH":   true, // body is allowed
	"DELETE":  false,
	"TRACE":   false,
	"DEBUG":   false,
}
