package app

import (
	"net/http"
)

type API struct{}

type RouteFunc func(mux *http.ServeMux)

func NewMux(routes ...RouteFunc) *http.ServeMux {
	mux := http.NewServeMux()

	for i := range routes {
		routes[i](mux)
	}

	return mux
}
