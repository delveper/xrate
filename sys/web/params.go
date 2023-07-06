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

// FromHeader retrieves a value from the request's headers
// and optionally strips a specified prefix or nil if no header is found.
func FromHeader(req *http.Request, name string, prefix string) string {
	val := req.Header.Get(name)
	if val == "" || !strings.HasPrefix(strings.ToLower(val), prefix) {
		return ""
	}

	return val
}
