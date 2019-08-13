package httpx

import (
	"net/http"
	"strings"
)

const (
	HeaderAccept        = "Accept"
	HeaderAuthorization = "Authorization"
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
	HeaderHost          = "Host"
	HeaderOrigin        = "Origin"
	HeaderReferer       = "Referer"
	HeaderUserAgent     = "User-Agent"

	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlMaxAgeSeconds    = "Access-Control-Max-Age"

	HeaderXForwardedFor = "X-Forwarded-For"
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
		HeaderAuthorization,
		HeaderContentLength,
		HeaderContentType,
		HeaderHost,
		HeaderOrigin,
		HeaderReferer,
	}, ", ")
)

const (
	ContentTypeOctetStream    = "application/octet-stream"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)
