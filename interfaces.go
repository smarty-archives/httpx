package httpx

import "net/http"

type WaitGroup interface {
	Add(delta int)
	Wait()
	Done()
}

type Sender interface {
	Send(interface{}) interface{}
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
