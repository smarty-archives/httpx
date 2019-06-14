package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/smartystreets/httpx"
)

type ContentHandler struct {
	inner http.Handler

	payload     []byte
	contentType string
}

func (this *ContentHandler) Install(inner http.Handler) {
	this.inner = inner
}

func NewContentStringHandler(payload, contentType string) *ContentHandler {
	return NewContentHandler([]byte(payload), contentType)
}
func NewContentHandler(payload []byte, contentType string) *ContentHandler {
	if strings.TrimSpace(contentType) == "" {
		contentType = httpx.ContentTypeOctetStream
	}

	return &ContentHandler{inner: new(NoopHandler), payload: payload, contentType: contentType}
}
func NewKeyValueContentHandler(pairs ...string) *ContentHandler {
	builder := bytes.NewBufferString(keyValueContentSyntax)
	for x := 0; x < len(pairs); x += 2 {
		_, _ = fmt.Fprintf(builder, "%s=%s\n", pairs[x], pairs[x+1])
	}
	return NewContentStringHandler(builder.String(), "text/plain; charset=utf-8")
}

func (this *ContentHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.inner.ServeHTTP(response, request)

	response.Header()[httpx.HeaderContentType] = []string{this.contentType}
	_, _ = response.Write(this.payload)
}

const keyValueContentSyntax = `# Syntax:
# Lines beginning with '#' or '//' are to be considered comments.
# Blank lines or lines consisting of only whitespace characters are to be considered comments.
# Each non-comment line contains a key string followed by an equals sign, followed by a corresponding value string.
# Example Line: package-name=1.2.3
# The number of line items and their order herein are subject to change without notice.

`
