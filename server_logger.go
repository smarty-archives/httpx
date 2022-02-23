package httpx

import (
	"bytes"
	"log"
)

type serverLogger struct{ logger }

func newServerLogger() *log.Logger {
	return log.New(&serverLogger{logger: log.Default()}, "", 0)
}

func (this *serverLogger) Write(buffer []byte) (int, error) {
	length := len(buffer)

	if length == 0 {
		return 0, nil
	}

	buffer = buffer[0 : length-1] // trim trailing line break

	if bytes.HasSuffix(buffer, ignoredLogStatements) {
		return length, nil // no-op
	}

	this.logger.Printf("[WARN] %s", buffer)
	return length, nil
}

var ignoredLogStatements = []byte("golang.org/issue/25192") // "http: URL query contains semicolon, which is no longer a supported separator; parts of the query may be stripped when parsed; see golang.org/issue/25192"

type logger interface {
	Printf(string, ...interface{})
}
