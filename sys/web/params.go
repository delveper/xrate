package web

import (
	"net/http"
	"strings"
)

// FromRequest retrieves the value of a specified
// URL query parameter or nil if no value is found.
func FromRequest(req *http.Request, key string) *string {
	var val string
	if val = req.FormValue(key); val == "" {
		if val = req.URL.Query().Get(key); val == "" {
			return nil
		}
	}

	return &val
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
func FromHeader(req *http.Request, name string, prefix string) *string {
	val := req.Header.Get(name)
	if val == "" {
		return nil
	}

	if prefix != "" {
		if !(len(val) > len(prefix) && strings.ToLower(val[:len(prefix)]) == prefix) {
			return nil
		}

		v := val[len(prefix)+1:]

		return &v
	}

	return &val
}

// FromCookie retrieves the value
// of a specified cookie from the request.
func FromCookie(req *http.Request, name string) *string {
	cookie, err := req.Cookie(name)
	if err != nil {
		return nil
	}

	return &cookie.Value
}
