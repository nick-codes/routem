package trie

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/nick-codes/routem"

	"golang.org/x/net/context"

	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	keyType int
)

const (
	zerothKey keyType = 0
	firstKey  keyType = 1
	secondKey keyType = 2
)

type testRoute struct {
	path         string
	method       []routem.Method
	handler      routem.HandlerFunc
	errorHandler routem.ErrorHandlerFunc
}

func (t *testRoute) Handler() routem.HandlerFunc {
	if t.handler != nil {
		return t.handler
	}
	return func(ctx context.Context) routem.HTTPError {

		if 2 == ctx.Value(secondKey) {
			return nil
		}

		return routem.NewHTTPError(500, fmt.Errorf("Value 2 not in context."))
	}
}

func (t *testRoute) Methods() []routem.Method {
	if len(t.method) > 0 {
		return t.method
	}
	return routem.GetMethod
}

func (t *testRoute) Prefix(prefix string) routem.Route {
	return &testRoute{path: prefix + t.path}
}

func (t *testRoute) WithTimeout(time.Duration) routem.RouteConfigurator {
	return t
}

func (t *testRoute) WithErrorHandler(routem.ErrorHandlerFunc) routem.RouteConfigurator {
	return t
}

func (t *testRoute) WithMiddleware(routem.MiddlewareFunc) routem.RouteConfigurator {
	return t
}

func (t *testRoute) WithMiddlewares([]routem.MiddlewareFunc) routem.RouteConfigurator {
	return t
}

func (t *testRoute) Middlewares() []routem.MiddlewareFunc {
	return []routem.MiddlewareFunc{
		func(next routem.HandlerFunc) routem.HandlerFunc {
			return func(ctx context.Context) routem.HTTPError {
				ctx = context.WithValue(ctx, firstKey, 1)
				return next(ctx)
			}
		},
		func(next routem.HandlerFunc) routem.HandlerFunc {
			return func(ctx context.Context) routem.HTTPError {
				if 1 == ctx.Value(firstKey) {
					ctx = context.WithValue(ctx, secondKey, 2)
					return next(ctx)
				}
				return routem.NewHTTPError(500, fmt.Errorf("First value not in context."))
			}
		},
	}
}

func (*testRoute) Timeout() time.Duration {
	return routem.DefaultTimeout
}

func (t *testRoute) ErrorHandler() routem.ErrorHandlerFunc {
	return t.errorHandler
}

func (r *testRoute) Path() string {
	return r.path
}

// Helpers
func assertError(t *testing.T, routes []routem.Route) {
	factory := NewHandlerFactory(nil, nil)

	handler, err := factory.Handler(routes)

	t.Logf("%s", err)

	assert.Nil(t, handler)
	assert.NotNil(t, err)
}

func assertSuccess(t *testing.T, routes []routem.Route) {
	factory := NewHandlerFactory(nil, nil)

	handler, err := factory.Handler(routes)

	t.Logf("%#v", handler)

	assert.NotNil(t, handler)
	assert.Nil(t, err)
}

// Real tests begin here

func TestNewHandlerFactory(t *testing.T) {
	factory := NewHandlerFactory(nil, nil)

	assert.NotNil(t, factory)
}

func TestErrorWithNilRoutes(t *testing.T) {
	assertError(t, nil)
}

func TestErrorWithEmptyRoutes(t *testing.T) {
	routes := make([]routem.Route, 0)
	assertError(t, routes)
}

func TestErrorWithNilRoute(t *testing.T) {
	routes := make([]routem.Route, 1, 1)
	assertError(t, routes)
}

func TestErrorWithDuplicateRoutes(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
		&testRoute{path: "/test"},
	}
	assertError(t, routes)
}

func TestErrorWithDuplicateParamName(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/:name/:name"},
	}
	assertError(t, routes)
}

func TestErrorWithUnnamedParamter(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/:"},
	}
	assertError(t, routes)
}

func TestErrorWithZeroLengthRoute(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: ""},
	}
	assertError(t, routes)
}

func TestErrorWithNonSlashRoute(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "no/slash"},
	}
	assertError(t, routes)
}

func TestErrorWithDuplicateParamRoutesDifferentNames(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/:name"},
		&testRoute{path: "/:id"},
	}
	assertError(t, routes)
}

func TestErrorWithDuplicateParamRoutes(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test/:name"},
		&testRoute{path: "/test/:name"},
	}
	assertError(t, routes)
}

func TestErrorWithDoubleSlashRoute(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/:name//"},
	}
	assertError(t, routes)
}

func TestShallowSuccess(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/name"},
	}
	assertSuccess(t, routes)
}

func TestSuccessRoot(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/"},
	}
	assertSuccess(t, routes)
}

func TestDuplicateRouteDifferentMethods(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
		&testRoute{path: "/test", method: routem.PutMethod},
	}
	assertSuccess(t, routes)
}

func assertServer(t *testing.T, routes []routem.Route, method routem.Method, url string) *httptest.ResponseRecorder {
	factory := NewHandlerFactory(context.Background(), nil)
	return assertServerFactory(t, factory, routes, method, url)
}

func assertServerFactory(t *testing.T, factory routem.HandlerFactory, routes []routem.Route, method routem.Method, url string) *httptest.ResponseRecorder {

	handler, err := factory.Handler(routes)

	assert.NotNil(t, handler)
	assert.Nil(t, err)

	response := httptest.NewRecorder()
	request, err := http.NewRequest(string(method), url, nil)
	assert.Nil(t, err)
	handler.ServeHTTP(response, request)

	return response
}

func TestSimpleRouting(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
	}

	assertServer(t, routes, routem.Get, "http://localhost/test")
}

func TestParamRouting(t *testing.T) {
	called := false
	routes := []routem.Route{
		&testRoute{
			path: "/test/:id/:name",
			handler: func(ctx context.Context) routem.HTTPError {
				params := routem.ParamsFromContext(ctx)
				assert.NotNil(t, params)
				assert.Equal(t, "5", params["id"])
				assert.Equal(t, "bill", params["name"])
				called = true
				return nil
			},
		},
	}

	response := assertServer(t, routes, routem.Get, "http://localhost/test/5/bill")

	assert.True(t, called)
	assert.Equal(t, 200, response.Code)
}

func TestTimeout(t *testing.T) {
	done := make(chan struct{})
	routes := []routem.Route{
		&testRoute{
			path: "/test",
			handler: func(ctx context.Context) routem.HTTPError {
				time.Sleep(3 * time.Second)
				assert.NotNil(t, ctx.Err())
				close(done)
				return nil
			},
		},
	}

	response := assertServer(t, routes, routem.Get, "http://localhost/test")
	assert.Equal(t, 500, response.Code)
	// Wait for the check for context error
	<-done
}

func TestRouteNotFoundMethod(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
	}

	response := assertServer(t, routes, routem.Put, "http://localhost/test")
	assert.Equal(t, 404, response.Code)
}

func TestDeepRouteNotFoundMethod(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test/a"},
		&testRoute{path: "/test/b"},
	}

	response := assertServer(t, routes, routem.Put, "http://localhost/test/c")
	assert.Equal(t, 404, response.Code)
}

func TestFactoryErrorHandler(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
	}
	factory := NewHandlerFactory(context.Background(), func(err routem.HTTPError, ctx context.Context) error {
		response := routem.ResponseWriterFromContext(ctx)
		http.Error(response, "Converted to 400", 400)
		return nil
	})
	response := assertServerFactory(t, factory, routes, routem.Put, "http://localhost/test/c")
	assert.Equal(t, 400, response.Code)
}

func TestErrorWithFactoryErrorHandler(t *testing.T) {
	routes := []routem.Route{
		&testRoute{path: "/test"},
	}
	factory := NewHandlerFactory(context.Background(), func(err routem.HTTPError, ctx context.Context) error {
		return fmt.Errorf("Error handling an error: %v", err)
	})
	response := assertServerFactory(t, factory, routes, routem.Put, "http://localhost/test/c")
	assert.Equal(t, 500, response.Code)
}

func TestRouteErrorHandler(t *testing.T) {
	routes := []routem.Route{
		&testRoute{
			path: "/test",
			handler: func(ctx context.Context) routem.HTTPError {
				return routem.NewHTTPError(304, fmt.Errorf("Moved Permanently!"))
			},
			errorHandler: func(err routem.HTTPError, ctx context.Context) error {
				response := routem.ResponseWriterFromContext(ctx)
				http.Error(response, "Converted to 400", 400)
				return nil
			},
		},
	}

	response := assertServer(t, routes, routem.Get, "http://localhost/test")
	assert.Equal(t, 400, response.Code)
}

func TestErrorRouteErrorHandler(t *testing.T) {
	routes := []routem.Route{
		&testRoute{
			path: "/test",
			handler: func(ctx context.Context) routem.HTTPError {
				return routem.NewHTTPError(304, fmt.Errorf("Moved Permanently!"))
			},
			errorHandler: func(err routem.HTTPError, ctx context.Context) error {
				return fmt.Errorf("Error handling an error: %v", err)
			},
		},
	}

	response := assertServer(t, routes, routem.Get, "http://localhost/test")
	assert.Equal(t, 500, response.Code)
}

func TestRootContextPassed(t *testing.T) {
	called := false
	routes := []routem.Route{
		&testRoute{
			path: "/test",
			handler: func(ctx context.Context) routem.HTTPError {
				assert.Equal(t, 0, ctx.Value(zerothKey))
				called = true
				return nil
			},
		},
	}

	ctx := context.WithValue(context.Background(), zerothKey, 0)
	factory := NewHandlerFactory(ctx, nil)
	assertServerFactory(t, factory, routes, routem.Get, "http://localhost/test")
	assert.True(t, called)
}
