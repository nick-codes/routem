package routem

import (
	"net/http"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

const (
	testPath = "test/path"
)

var (
	testHandler HandlerFunc = func(c context.Context) HTTPError { return nil }
)

type (
	testHTTPHandler struct {
		called *bool
	}
)

func (t testHTTPHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	*(t.called) = true
}

func assertTestRoute(t *testing.T, route Route) {
	assertTestConfig(t, route)

	assertRoute(t, route)
}

func assertDefaultRoute(t *testing.T, route Route) {
	assertDefaultConfig(t, route)

	assertRoute(t, route)
}

func assertTestRouteWithMethods(t *testing.T, route Route, methods []Method) {
	assertTestConfig(t, route)

	assertRouteWithMethods(t, route, methods)
}

func assertRoute(t *testing.T, route Route) {
	assertRouteWithMethods(t, route, GetMethod)
}

func assertRouteWithMethods(t *testing.T, route Route, methods []Method) {
	assert.Equal(t, len(methods), len(route.Methods()), "Wrong method count")
	for i, method := range methods {
		assert.Equal(t, method, route.Methods()[i], "Wrong method")
	}
	assert.Equal(t, testPath, route.Path(), "Wrong path")
	assert.NotNil(t, route.Handler(), "No Handler")
}

func assertDefaultHTTPRoute(t *testing.T, route Route, testHandler testHTTPHandler) {
	assertDefaultConfig(t, route)
	assertRoute(t, route)
	assertHTTPRoute(t, route, testHandler)
}

func assertTestHTTPRoute(t *testing.T, route Route, testHandler testHTTPHandler) {
	assertTestConfig(t, route)
	assertRoute(t, route)
	assertHTTPRoute(t, route, testHandler)
}

func assertTestHTTPRouteWithMethods(t *testing.T, route Route, testHandler testHTTPHandler, methods []Method) {
	assertTestConfig(t, route)
	assertRouteWithMethods(t, route, methods)
	assertHTTPRoute(t, route, testHandler)
}

func assertHTTPRoute(t *testing.T, route Route, testHandler testHTTPHandler) {
	response, request, params := setupContextTest(t)

	ctx, _ := newContext(route.Timeout(), request, response, params)

	err := route.Handler()(ctx)

	assert.Nil(t, err, "Error from handler?")

	assert.True(t, *(testHandler.called), "Not called!")
}

func TestNewRoute(t *testing.T) {
	config := defaultConfig()
	route := newRoute(config, GetMethod, testPath, testHandler)

	assertDefaultRoute(t, route)
}

func TestNewHttpRoute(t *testing.T) {
	config := defaultConfig()

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}

	route := newHTTPRoute(config, GetMethod, testPath, testHandler)

	assertDefaultHTTPRoute(t, route, testHandler)
}

func TestPrefix(t *testing.T) {
	config := defaultConfig()

	parts := strings.Split(testPath, "/")
	route := newRoute(config, GetMethod, parts[1], testHandler)
	route = route.Prefix(parts[0] + "/")

	assertRoute(t, route)
}
