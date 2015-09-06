package routem

import (
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

func (r *router) handler() (http.Handler, error) {
	return r.factory.Handler(r.routes)
}
