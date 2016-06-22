package routem

import (
	"fmt"
	"net/http"
)

type (
	router struct {
		creator
		factory HandlerFactory
	}
)

// =-=-=-=-=-=
// Constructor
// =-=-=-=-=-=

// NewRouter constructs a new Router which can be used to setup and run
// the routem Router.
func NewRouter(factory HandlerFactory) Router {
	return &router{
		creator: newCreator(defaultConfig()),
		factory: factory,
	}
}

// =-=-=-=
// Execute
// =-=-=-=

func (r *router) Run(address string) (Service, error) {
	handler, err := r.handler()

	if err != nil {
		return nil, err
	}

	s := newService(address, handler)

	err = s.run()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *router) RunTLS(address, certFile, keyFile string) (Service, error) {
	handler, err := r.handler()

	if err != nil {
		return nil, err
	}

	s := newService(address, handler)

	err = s.runTLS(certFile, keyFile)

	if err != nil {
		return nil, err
	}

	return s, nil
}

// =-=-=-=
// Helpers
// =-=-=-=

func flatten(prefix string, routes []Routable) ([]Route, error) {
	flat := make([]Route, 0, len(routes))
	for _, route := range routes {
		group, isGroup := route.(Group)
		if isGroup {
			groupRoutes, err := flatten(prefix+group.Path(), group.Routes())
			if err != nil {
				return nil, err
			}
			for _, gr := range groupRoutes {
				flat = append(flat, gr)
			}
		} else {
			route, isRoute := route.(Route)
			// This is impossible but just in case
			if !isRoute {
				return nil, fmt.Errorf("Found Routable not a Group or Route? WTF! %v", route)
			}
			flat = append(flat, route.Prefix(prefix))
		}
	}

	return flat, nil
}

func (r *router) handler() (http.Handler, error) {
	routes, err := flatten("", r.Routes())

	if err == nil {
		return r.factory.Handler(routes)
	}

	return nil, err
}
