package httpx

import (
	"context"
	"net/http"
)

func InitializeContext(request *http.Request) *http.Request {
	parent := request.Context()
	child := context.WithValue(parent, contextNamespace, make(Namespace))
	return request.WithContext(child)
}

func Context(request *http.Request) Namespace {
	return request.Context().Value(contextNamespace).(Namespace)
}

type Namespace map[interface{}]interface{}

func (this Namespace) Int(key string) int {
	value, _ := this[key].(int)
	return value
}
func (this Namespace) Int64(key string) int64 {
	value, _ := this[key].(int64)
	return value
}

func (this Namespace) Uint64(key string) uint64 {
	value, _ := this[key].(uint64)
	return value
}

func (this Namespace) String(key string) string {
	value, _ := this[key].(string)
	return value
}

const contextNamespace = "smartystreets"
