package middleware

import (
	"net/http"
	"strings"
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
		contentType = "application/octet-stream"
	}

	return &ContentHandler{payload: payload, contentType: contentType}
}

func (this *ContentHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header()["Content-Type"] = []string{this.contentType}
	response.Write(this.payload)
}
