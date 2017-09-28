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

type Contextual http.Request

func (this Contextual) Set(key, value interface{}) {
	this.getContextValues()[key] = value
}

func (this Contextual) getContextValues() map[interface{}]interface{} {
	request := http.Request(this)
	namespace := request.Context().Value(contextNamespace)
	return namespace.(map[interface{}]interface{})
}

func (this Contextual) Int(key string) int {
	value, _ := this.getContextValues()[key].(int)
	return value
}
func (this Contextual) Int64(key string) int64 {
	value, _ := this.getContextValues()[key].(int64)
	return value
}

func (this Contextual) Uint64(key string) uint64 {
	value, _ := this.getContextValues()[key].(uint64)
	return value
}

func (this Contextual) String(key string) string {
	value, _ := this.getContextValues()[key].(string)
	return value
}

const contextNamespace = "smartystreets"
