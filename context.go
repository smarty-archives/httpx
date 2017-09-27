package httpx

import (
	"context"
	"net/http"
)

func InitializeContext(request *http.Request) *http.Request {
	parent := request.Context()
	child := context.WithValue(parent, contextNamespace, make(map[interface{}]interface{}))
	return request.WithContext(child)
}

const contextNamespace = "smartystreets"
