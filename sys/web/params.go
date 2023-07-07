package web

import (
	"net/http"
	"strings"
)

// FromQuery retrieves the value of a specified
// URL query parameter or nil if no value is found.
func FromQuery(req *http.Request, key string) string {
	val := req.FormValue(key)
	if val == "" {
		return ""
	}

	return val
}

// FromContext retrieves a pointer of a specific type
// from the request's context or nil if no value is found.
func FromContext[T, K any](req *http.Request, key K) *T {
	val := req.Context().Value(key)

	v, ok := val.(T)
	if !ok {
		return nil
	}

	return &v
}

// FromHeader retrieves a value from the request's headers
// and optionally strips a specified prefix or nil if no header is found.
func FromHeader(req *http.Request, name string, prefix string) string {
	val := req.Header.Get(name)
	if val == "" || !strings.HasPrefix(strings.ToLower(val), prefix) {
		return ""
	}

	return val
}
