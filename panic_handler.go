package httpx

import (
	"log"
	"net/http"
	"runtime/debug"
)

func PanicHandler(response http.ResponseWriter, _ *http.Request, err interface{}) {
	log.Printf(errorFormat, err, string(debug.Stack()))
	response.WriteHeader(http.StatusInternalServerError)
}

const errorFormat = `[ERROR] Recovered panic: %v
Stack Trace:
%s
------------------- END OF STACK TRACE -------------------`
