package routem

import (
	"net/http"

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

func assertRoute(t *testing.T, route *route) {
	assertDefaultConfig(t, route.config)

	assert.Equal(t, 1, len(route.Methods()), "Wrong method count")
	assert.Equal(t, Get, route.Methods()[0], "Wrong method")
	assert.Equal(t, testPath, route.Path(), "Wrong path")
	assert.NotNil(t, route.Handler(), "No Handler")
}

func TestNewRoute(t *testing.T) {
	config := defaultConfig()
	route := newRoute(config, GetMethod, testPath, testHandler)

	assertRoute(t, route)
}

func TestNewHttpRoute(t *testing.T) {
	config := defaultConfig()
	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := newHTTPRoute(config, GetMethod, testPath, testHandler)

	assertDefaultConfig(t, route.config)

	response, request, params := setupContextTest(t)

	ctx, _ := newContext(route.Timeout(), request, response, params)

	err := route.Handler()(ctx)

	assert.Nil(t, err, "Error from handler?")

	assert.True(t, called, "Not called!")
}
