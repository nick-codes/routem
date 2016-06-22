package routem

import (
	"fmt"
	"net/http"

	_ "net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

const (
	testAddress    = "localhost:9000"
	invalidAddress = "localhost:9000000000"
	testCert       = "test/test.crt"
	testKey        = "test/test.key"
	invalidCert    = "test/invalid.crt"
	invalidKey     = "test/invalid.key"
)

type (
	testHandlerFactory struct {
		error   bool
		routes  []Route
		handler http.Handler
	}
)

func (hf *testHandlerFactory) Handler(routes []Route) (http.Handler, error) {
	hf.routes = routes
	if hf.error {
		return nil, fmt.Errorf("test error")
	}
	return hf.handler, nil
}

func TestNewRouter(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	assertDefaultConfig(t, router)
}

func TestRun(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.Run(testAddress)

	require.Nil(t, err, "Run failed")
	require.NotNil(t, srv, "Nil service")

	assert.Equal(t, testAddress, srv.Address(), "Wrong address")
	assert.True(t, srv.IsRunning(), "Not running")

	wait := make(chan error)
	go func() {
		wait <- srv.Wait()
	}()

	err = srv.Stop()

	assert.Nil(t, err, "Failed to stop")
	assert.NotNil(t, <-wait, "Wait didn't err")
	assert.False(t, srv.IsRunning(), "Still running")
}

func TestRunWithFactoryError(t *testing.T) {
	hf := &testHandlerFactory{error: true}

	router := NewRouter(hf)

	srv, err := router.Run(testAddress)

	assert.NotNil(t, err, "Run didn't return an error")
	assert.Nil(t, srv, "Run returned a service")
}

func TestRunWithInvalidAddress(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.Run(invalidAddress)

	assert.NotNil(t, err, "Run didn't return an error")
	assert.Nil(t, srv, "Returned a service.")
}

func TestRunTLS(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.RunTLS(testAddress, testCert, testKey)

	assert.Nil(t, err, "RunTLS failed")

	require.Nil(t, err, "Run failed")
	require.NotNil(t, srv, "Nil service")

	err = srv.Stop()

	assert.Nil(t, err, "Failed to stop")
}

func TestRunTLSWithFactoryError(t *testing.T) {
	hf := &testHandlerFactory{error: true}

	router := NewRouter(hf)

	srv, err := router.RunTLS(testAddress, testCert, testKey)

	assert.NotNil(t, err, "RunTLS didn't return err")
	assert.Nil(t, srv, "RunTLS returned a service")
}

func TestRunTLSWithInvalidAddress(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.RunTLS(invalidAddress, testCert, testKey)

	assert.NotNil(t, err, "Run didn't return an error")
	assert.Nil(t, srv, "Returned a service.")
}

func TestRunTLSWithInvalidCert(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.RunTLS(testAddress, invalidCert, testKey)

	assert.NotNil(t, err, "Run didn't return an error")
	assert.Nil(t, srv, "Returned a service.")
}

func TestRunTLSWithInvalidKey(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf)

	srv, err := router.RunTLS(testAddress, testCert, invalidKey)

	assert.NotNil(t, err, "Run didn't return an error")
	assert.Nil(t, srv, "Returned a service.")
}

type nonRoute struct {
	config
}

func (nr *nonRoute) Path() string {
	return ""
}

var _ Routable = &nonRoute{}

func TestFlattenErrorsOnInvalidType(t *testing.T) {
	router := &router{}

	// This is impossible using the API but make sure
	// we handle it anyway
	router.routes = append(router.routes, &nonRoute{})

	_, err := router.handler()

	assert.NotNil(t, err)
}

func TestFlattenErrorsCascadeUpOnInvalidType(t *testing.T) {
	router := &router{}

	group := router.WithGroup("/").(*group)

	// This is impossible using the API but make sure
	// we handle it anyway
	group.routes = append(group.routes, &nonRoute{})

	_, err := router.handler()

	assert.NotNil(t, err)
}

func TestFlattenWithGroup(t *testing.T) {
	hf := &testHandlerFactory{}

	router := NewRouter(hf).(*router)

	router.Get("/blah", func(ctx context.Context) HTTPError { return nil })

	group := router.WithGroup("/")
	group.Get("test", func(ctx context.Context) HTTPError { return nil })

	_, err := router.handler()

	assert.Nil(t, err)

	assert.Equal(t, 2, len(hf.routes))

	assert.Equal(t, "/blah", hf.routes[0].Path())
	assert.Equal(t, "/test", hf.routes[1].Path())
}
