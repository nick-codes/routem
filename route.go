package routem

import (
	"net/http"

	"golang.org/x/net/context"
)

type (
	route struct {
		config

		methods []Method
		path    string
		handler HandlerFunc
	}
)

// =-=-=-=
// Getters
// =-=-=-=

func (r *route) Methods() []Method {
	return r.methods
}
func (r *route) Path() string {
	return r.path
}
func (r *route) Handler() HandlerFunc {
	return r.handler
}

// =-=-=-=
// Helpers
// =-=-=-=

func newRoute(defs config, methods []Method, path string, handler HandlerFunc) Route {

	r := &route{
		config:  newConfig(defs),
		methods: methods,
		path:    path,
		handler: handler,
	}

	return r
}

func newHTTPRoute(def config, methods []Method, path string, handler http.Handler) Route {

	route := newRoute(def, methods, path, wrapHTTPHandler(handler))

	return route
}

func wrapHTTPHandler(handler http.Handler) HandlerFunc {
	return func(c context.Context) HTTPError {
		request := RequestFromContext(c)
		response := ResponseWriterFromContext(c)

		handler.ServeHTTP(response, request)

		return nil
	}
}
