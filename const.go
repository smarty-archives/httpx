package httpx

import (
	"net/http"
	"strings"
)

const (
	HeaderUserAgent     = "User-Agent"
	HeaderAccept        = "Accept"
	HeaderReferer       = "Referer"
	HeaderOrigin        = "Origin"
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderHost          = "Host"

	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlMaxAgeSeconds    = "Access-Control-Max-Age"

	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderVia           = "Via"
)

var (
	DefaultCORSMethods = strings.Join([]string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodHead,
	}, ", ")
	DefaultCORSHeaders = strings.Join([]string{
		HeaderAccept,
		HeaderContentType,
		HeaderContentLength,
		HeaderReferer,
		HeaderOrigin,
		HeaderHost,
	}, ", ")
)

const (
	ContentTypeOctetStream    = "application/octet-stream"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)
