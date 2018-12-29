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
	if address, found := readCustomHeader(ReadHeader(request, customHeader)); found {
		return address
	} else if address, found = readViaForwardedFor(request); found {
		return address
	} else if address, found = readForwardedFor(request, 0); found {
		return address
	} else {
		return remoteAddressFromTCP(request.RemoteAddr)
	}
}
func readViaForwardedFor(request *http.Request) (string, bool) {
	if via := ReadHeader(request, HeaderVia); len(via) == 0 {
		return "", false
	}

	return readForwardedFor(request, 1)
}
func readCustomHeader(value string) (string, bool) {
	return value, len(value) > 0
}
func readForwardedFor(request *http.Request, depth int) (string, bool) {
	forwardedFor := ReadHeader(request, HeaderXForwardedFor)
	address := remoteAddressFromHeader(forwardedFor, depth)
	return address, len(address) > 0
}
func remoteAddressFromHeader(value string, depth int) string {
	if index := strings.LastIndex(value, ", "); index < 0 {
		return value // not found
	} else if depth == 0 {
		return value[index+2:] // last iteration, return it
	} else {
		return remoteAddressFromHeader(value[0:index], depth-1) // remove from the end
	}
}
func remoteAddressFromTCP(raw string) string {
	value, _, _ := net.SplitHostPort(raw)
	return value
}

func CalculateRequestSize(request *http.Request) int64 {
	return int64(calculatePreambleSize(request)) +
		int64(calculateHeaderSize(request)) +
		request.ContentLength + 1
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

const httpLineBreakLength = 2
