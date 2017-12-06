package httpx

import (
	"net"
	"net/http"
	"strings"
	"sync"
)

func NewWaitGroup(workers int) *sync.WaitGroup {
	waiter := &sync.WaitGroup{}
	waiter.Add(workers)
	return waiter
}

func ReadHeader(request *http.Request, canonicalHeaderName string) string {
	if values, contains := request.Header[canonicalHeaderName]; contains && len(values) > 0 {
		return values[0]
	} else {
		return ""
	}
}

func WriteHeader(request *http.Request, canonicalHeaderName string, value string) {
	request.Header[canonicalHeaderName] = []string{value}
}

func ReadClientIPAddress(request *http.Request, customHeader string) string {
	if len(customHeader) > 0 {
		if remoteAddress := ReadHeader(request, customHeader); len(remoteAddress) > 0 {
			return remoteAddressFromHeader(remoteAddress)
		}
	}

	if forwardedAddress := ReadHeader(request, HeaderXForwardedFor); len(forwardedAddress) > 0 {
		return remoteAddressFromHeader(forwardedAddress)
	}

	return remoteAddressFromTCP(request.RemoteAddr)
}
func remoteAddressFromHeader(value string) string {
	if index := strings.LastIndex(value, ", "); index >= 0 {
		return value[index+2:]
	}

	return value
}
func remoteAddressFromTCP(raw string) string {
	value, _, _ := net.SplitHostPort(raw)
	return value
}
