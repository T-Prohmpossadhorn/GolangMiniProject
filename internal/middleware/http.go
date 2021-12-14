// Package middleware contains HTTP middleware for servers
package middleware

import (
	"net/http"
	"strings"
)

// mimeType extracts the MIME type (media type) component from an RFC 7231 section 3.1.1.1 Content-Type value.
// The media type is returned in lower case.
func mimeType(contentType string) string {
	return strings.ToLower(strings.Split(contentType, ";")[0])
}

// JSONContentTypeValidator is a wrapper to reject requests based on the value of the Content-Type header.
// If the header is present, the MIME type component of the value must be "application/json" to be passed through.
// For other values JSONContentTypeValidator responds with http.StatusUnsupportedMediaType and aborts the request.
func JSONContentTypeValidator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Enforce application/json content-type for requests that have one set
		contentType := r.Header.Get("Content-Type")
		if len(contentType) != 0 && mimeType(contentType) != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		h.ServeHTTP(w, r)
	})
}
