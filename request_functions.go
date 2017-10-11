package httpx

import (
	"net/http"
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
