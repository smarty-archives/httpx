package httpx

import (
	"context"
	"net/http"
)

type Contextual http.Request

func (this Contextual) Set(key, value interface{}) {
	request := http.Request(this)
	namespace := request.Context().Value(contextNamespace)
	values := namespace.(map[interface{}]interface{})
	values[key] = value
}

func InitializeContext(request *http.Request) *http.Request {
	parent := request.Context()
	child := context.WithValue(parent, contextNamespace, make(map[interface{}]interface{}))
	return request.WithContext(child)
}

const contextNamespace = "smartystreets"
