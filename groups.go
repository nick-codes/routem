package routem

import (
	"net/http"
)

type (
	group struct {
		config

		prefix  string
		routes  []Routable
		handler HandlerFunc
	}
)

// =-=-=-=
// Setters
// =-=-=-=

func (g *group) With(methods []Method, path string, handler HandlerFunc) *route {
	route := newRoute(g.config, methods, path, handler)

	g.routes = append(g.routes, route)

	return route
}
func (g *group) WithHTTP(methods []Method, path string, handler http.Handler) *route {
	route := newHTTPRoute(g.config, methods, path, handler)

	g.routes = append(g.routes, route)

	return route
}

func (g *group) WithGroup(path string) *group {
	group := newGroup(g.config, path)

	g.routes = append(g.routes, group)

	return group
}

// =-=-=-=-=-=-=-=-=-=
// HandlerFunc Aliases
// =-=-=-=-=-=-=-=-=-=

func (g *group) Noop(path string, handler HandlerFunc) *route {
	return g.With(NoMethod, path, handler)
}
func (g *group) Connect(path string, handler HandlerFunc) *route {
	return g.With(ConnectMethod, path, handler)
}
func (g *group) Delete(path string, handler HandlerFunc) *route {
	return g.With(DeleteMethod, path, handler)
}
func (g *group) Get(path string, handler HandlerFunc) *route {
	return g.With(GetMethod, path, handler)
}
func (g *group) Head(path string, handler HandlerFunc) *route {
	return g.With(HeadMethod, path, handler)
}
func (g *group) Options(path string, handler HandlerFunc) *route {
	return g.With(OptionsMethod, path, handler)
}
func (g *group) Patch(path string, handler HandlerFunc) *route {
	return g.With(PatchMethod, path, handler)
}
func (g *group) Put(path string, handler HandlerFunc) *route {
	return g.With(PutMethod, path, handler)
}
func (g *group) Post(path string, handler HandlerFunc) *route {
	return g.With(PostMethod, path, handler)
}
func (g *group) Trace(path string, handler HandlerFunc) *route {
	return g.With(TraceMethod, path, handler)
}
func (g *group) Crud(path string, handler HandlerFunc) *route {
	return g.With(CrudMethod, path, handler)
}
func (g *group) Any(path string, handler HandlerFunc) *route {
	return g.With(AnyMethod, path, handler)
}

// =-=-=-=-=-=-=-=-=-=-=
// http.Handler Aliases
// =-=-=-=-=-=-=-=-=-=-=

func (g *group) NoopHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(NoMethod, path, handler)
}
func (g *group) ConnectHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(ConnectMethod, path, handler)
}
func (g *group) DeleteHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(DeleteMethod, path, handler)
}
func (g *group) GetHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(GetMethod, path, handler)
}
func (g *group) HeadHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(HeadMethod, path, handler)
}
func (g *group) OptionsHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(OptionsMethod, path, handler)
}
func (g *group) PatchHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(PatchMethod, path, handler)
}
func (g *group) PutHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(PutMethod, path, handler)
}
func (g *group) PostHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(PostMethod, path, handler)
}
func (g *group) TraceHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(TraceMethod, path, handler)
}
func (g *group) CrudHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(CrudMethod, path, handler)
}
func (g *group) AnyHTTP(path string, handler http.Handler) *route {
	return g.WithHTTP(AnyMethod, path, handler)
}

// =-=-=-=
// Getters
// =-=-=-=

func (g *group) Routes() []Routable {
	return g.routes
}
func (g *group) Path() string {
	return g.prefix
}
func (g *group) Handler() HandlerFunc {
	return g.handler
}

// =-=-=-=
// Helpers
// =-=-=-=

func newGroup(defs config, prefix string) *group {
	g := &group{
		config: newConfig(defs),
		prefix: prefix,
	}

	return g
}
