package httpx

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestReadClientAddressFixture(t *testing.T) {
	gunit.Run(new(ReadClientAddressFixture), t)
}

type ReadClientAddressFixture struct {
	*gunit.Fixture
	request *http.Request
}

func (this *ReadClientAddressFixture) TestCorrectRemoteIPAddressSetOnRequest() {
	this.assertOriginIP("[::1]:81", "", "::1")
	this.assertOriginIP("[::1]:80", "1.2.3.4", "1.2.3.4")
	this.assertOriginIP("1.2.3.4:1234", "", "1.2.3.4")
	this.assertOriginIP("1.2.3.4:1234", "::1, 5.6.7.8", "5.6.7.8")
	this.assertOriginIP("1.2.3.4:1234", "a.b.c.d, e.f.g.h, i.j.k.l", "i.j.k.l")
}

func (this *ReadClientAddressFixture) assertOriginIP(remoteAddress, forwardedAddress, expectedAddress string) {
	this.setupRemoteRequest(remoteAddress, forwardedAddress)
	this.So(ReadClientIPAddress(this.request, ""), should.Equal, expectedAddress)
}
func (this *ReadClientAddressFixture) setupRemoteRequest(remoteAddress, forwardedAddress string) {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.request.RemoteAddr = remoteAddress
	WriteHeader(this.request, HeaderXForwardedFor, forwardedAddress)
}

func (this *ReadClientAddressFixture) TestPreferTrustedHeaderForIPAddressWhenAvailable() {
	this.setupRemoteRequest("1.2.3.4", "5.6.7.8")
	WriteHeader(this.request, "X-Remote-Address", "a.b.c.d")

	this.So(ReadClientIPAddress(this.request, "X-Remote-Address"), should.Equal, "a.b.c.d")
}

/////////////////////////////////////////////////////////////

func TestRequestFixture(t *testing.T) {
	gunit.Run(new(RequestFixture), t)
}

type RequestFixture struct {
	*gunit.Fixture
}

func (this *RequestFixture) TestCalculateRequestBodySize() {
	buffer := bytes.NewBuffer([]byte{})
	io.WriteString(buffer, "Line 1\n")
	io.WriteString(buffer, "Line 2\n")
	io.WriteString(buffer, "Line 3")

	request, _ := http.NewRequest("POST", "https://user:pass@domain.com:1234/path/doc?query=value", buffer)
	request.Header["Host"] = []string{"domain.com:1234"}
	request.Header["User-Agent"] = []string{"my user agent"}
	request.Header["Content-Type"] = []string{"application/json", "application/xml"}

	raw, _ := httputil.DumpRequest(request, true)
	this.So(CalculateRequestSize(request), should.Equal, len(raw))
}
