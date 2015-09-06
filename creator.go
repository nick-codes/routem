package routem

import (
	"net/http"
)

type (
	creator struct {
		config
		routes []Routable
	}
)

// =-=-=-=
// Setters
// =-=-=-=

func (c *creator) With(methods []Method, path string, handler HandlerFunc) Route {
	route := newRoute(c.config, methods, path, handler)

	c.routes = append(c.routes, route)

	return route
}

func (c *creator) WithHTTP(methods []Method, path string, handler http.Handler) Route {
	route := newHTTPRoute(c.config, methods, path, handler)

	c.routes = append(c.routes, route)

	return route
}

func (c *creator) WithGroup(path string) Group {
	group := newGroup(c.config, path)

	c.routes = append(c.routes, group)

	return group
}

// =-=-=-=-=-=-=-=-=-=
// HandlerFunc Aliases
// =-=-=-=-=-=-=-=-=-=

func (c *creator) Noop(path string, handler HandlerFunc) Route {
	return c.With(NoMethod, path, handler)
}
func (c *creator) Connect(path string, handler HandlerFunc) Route {
	return c.With(ConnectMethod, path, handler)
}
func (c *creator) Delete(path string, handler HandlerFunc) Route {
	return c.With(DeleteMethod, path, handler)
}
func (c *creator) Get(path string, handler HandlerFunc) Route {
	return c.With(GetMethod, path, handler)
}
func (c *creator) Head(path string, handler HandlerFunc) Route {
	return c.With(HeadMethod, path, handler)
}
func (c *creator) Options(path string, handler HandlerFunc) Route {
	return c.With(OptionsMethod, path, handler)
}
func (c *creator) Patch(path string, handler HandlerFunc) Route {
	return c.With(PatchMethod, path, handler)
}
func (c *creator) Put(path string, handler HandlerFunc) Route {
	return c.With(PutMethod, path, handler)
}
func (c *creator) Post(path string, handler HandlerFunc) Route {
	return c.With(PostMethod, path, handler)
}
func (c *creator) Trace(path string, handler HandlerFunc) Route {
	return c.With(TraceMethod, path, handler)
}
func (c *creator) Crud(path string, handler HandlerFunc) Route {
	return c.With(CrudMethod, path, handler)
}
func (c *creator) Any(path string, handler HandlerFunc) Route {
	return c.With(AnyMethod, path, handler)
}

// =-=-=-=-=-=-=-=-=-=-=
// http.Handler Aliases
// =-=-=-=-=-=-=-=-=-=-=

func (c *creator) NoopHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(NoMethod, path, handler)
}
func (c *creator) ConnectHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(ConnectMethod, path, handler)
}
func (c *creator) DeleteHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(DeleteMethod, path, handler)
}
func (c *creator) GetHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(GetMethod, path, handler)
}
func (c *creator) HeadHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(HeadMethod, path, handler)
}
func (c *creator) OptionsHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(OptionsMethod, path, handler)
}
func (c *creator) PatchHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(PatchMethod, path, handler)
}
func (c *creator) PutHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(PutMethod, path, handler)
}
func (c *creator) PostHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(PostMethod, path, handler)
}
func (c *creator) TraceHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(TraceMethod, path, handler)
}
func (c *creator) CrudHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(CrudMethod, path, handler)
}
func (c *creator) AnyHTTP(path string, handler http.Handler) Route {
	return c.WithHTTP(AnyMethod, path, handler)
}

// =-=-=-=
// Getters
// =-=-=-=

func (c *creator) Routes() []Routable {
	return c.routes
}

// =-=-=-=
// Helpers
// =-=-=-=

func newCreator(defs config) creator {
	return creator{
		config: newConfig(defs),
		routes: []Routable{},
	}
}
