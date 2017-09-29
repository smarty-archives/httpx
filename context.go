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

func SetContext(request *http.Request, key, value interface{}) {
	contextValues(request)[key] = value
}

func GetContext(request *http.Request, key interface{}) interface{} {
	return contextValues(request)[key]
}
func GetContextInt(request *http.Request, key interface{}) int {
	value, _ := GetContext(request, key).(int)
	return value
}
func GetContextInt64(request *http.Request, key interface{}) int64 {
	value, _ := GetContext(request, key).(int64)
	return value
}
func GetContextUint64(request *http.Request, key interface{}) uint64 {
	value, _ := GetContext(request, key).(uint64)
	return value
}
func GetContextString(request *http.Request, key interface{}) string {
	value, _ := GetContext(request, key).(string)
	return value
}

func contextValues(request *http.Request) map[interface{}]interface{} {
	return request.Context().Value(contextNamespace).(map[interface{}]interface{})
}

const contextNamespace = "smartystreets"
