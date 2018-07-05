package httpx

import (
	"net"
	"net/http"
	"net/url"
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

func CalculateRequestSize(request *http.Request) int {
	return calculatePreambleSize(request) +
		calculateHeaderSize(request) +
		calculateBodySize(request) +
		1 // off by one
}
func calculatePreambleSize(request *http.Request) (size int) {
	size += len(request.Method)
	size += 1 // space
	size += calculateURLSize(request.URL)
	size += 1 // space
	size += len(request.Proto)
	return size + httpLineBreakLength
}
func calculateURLSize(value *url.URL) int {
	pathSize := len(value.Path)

	if querySize := len(value.RawQuery); querySize > 0 {
		return pathSize + querySize
	} else {
		return pathSize
	}
}
func calculateHeaderSize(request *http.Request) (size int) {
	for key, values := range request.Header {
		for _, value := range values {
			size += calculateHeaderLineSize(key, value)
		}
	}

	return size + httpLineBreakLength
}
func calculateHeaderLineSize(key, value string) (size int) {
	size += len(key)
	size += 2 // colon and space characters
	size += len(value)
	return size + httpLineBreakLength
}
func calculateBodySize(request *http.Request) int {
	return int(request.ContentLength)
}

const httpLineBreakLength = 2
