package routem

import (
	"time"
)

type (
	config struct {
		errorHandler ErrorHandlerFunc
		timeout      time.Duration
		middlewares  []MiddlewareFunc
	}
)

func newConfig(defs config) config {
	return config{
		timeout:      defs.timeout,
		errorHandler: defs.errorHandler,
		middlewares:  defs.middlewares,
	}
}

func defaultConfig() config {
	return config{
		timeout:      DefaultTimeout,
		errorHandler: nil,
		middlewares:  []MiddlewareFunc{},
	}
}

func (c *config) WithErrorHandler(handler ErrorHandlerFunc) RouteConfigurator {
	c.errorHandler = handler
	return c
}

func (c *config) WithTimeout(t time.Duration) RouteConfigurator {
	c.timeout = t
	return c
}

func (c *config) WithMiddleware(middleware MiddlewareFunc) RouteConfigurator {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

func (c *config) WithMiddlewares(middlewares []MiddlewareFunc) RouteConfigurator {
	for _, middleware := range middlewares {
		c.WithMiddleware(middleware)
	}
	return c
}

func (c *config) Timeout() time.Duration {
	return c.timeout
}

func (c *config) ErrorHandler() ErrorHandlerFunc {
	return c.errorHandler
}

func (c *config) Middlewares() []MiddlewareFunc {
	return c.middlewares
}
