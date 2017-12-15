package middleware

import (
	"net/http"
	"strings"

	"github.com/smartystreets/httpx"
)

type ContentHandler struct {
	payload     []byte
	contentType string
}

func NewContentStringHandler(payload, contentType string) *ContentHandler {
	return NewContentHandler([]byte(payload), contentType)
}
func NewContentHandler(payload []byte, contentType string) *ContentHandler {
	if strings.TrimSpace(contentType) == "" {
		contentType = httpx.ContentTypeOctetStream
	}

	return &ContentHandler{payload: payload, contentType: contentType}
}

func (this *ContentHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header()[httpx.HeaderContentType] = []string{this.contentType}
	response.Write(this.payload)
}
