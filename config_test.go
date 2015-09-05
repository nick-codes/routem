package routem

import (
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	testErrorHandler  ErrorHandlerFunc = func(error HTTPError, ctx context.Context) error { return nil }
	testMiddleware    MiddlewareFunc   = func(HandlerFunc) HandlerFunc { return nil }
	testMiddlewareTwo MiddlewareFunc   = func(HandlerFunc) HandlerFunc { return nil }
)

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()
	assert.Equal(t, DefaultTimeout, config.timeout, "Config didn't have default timeout.")
	assert.Nil(t, config.errorHandler, "Default error handler provided?")
	assert.NotNil(t, config.middlewares, "Nil default middleware?")
	assert.Equal(t, 0, len(config.middlewares), "Default middleware?")
}

func TestNewConfig(t *testing.T) {
	orig := config{
		timeout:      DefaultTimeout,
		errorHandler: testErrorHandler,
		middlewares:  []MiddlewareFunc{testMiddleware},
	}
	config := newConfig(orig)
	assert.Equal(t, orig.timeout, config.timeout, "Incorrect Timeout")
	assert.Equal(t, len(orig.middlewares), len(config.middlewares), "Incorrect middlewares count")
}

func TestWithErrorHandler(t *testing.T) {
	config := defaultConfig()

	assert.Nil(t, config.errorHandler, "Incorrect default error handler")

	config.WithErrorHandler(testErrorHandler)

	assert.NotNil(t, config.errorHandler, "Incorrect error handler in struct")
	assert.NotNil(t, config.ErrorHandler(), "Incorrect error handler via function")
}

func TestWithTimeout(t *testing.T) {
	config := defaultConfig()

	assert.Equal(t, DefaultTimeout, config.timeout, "Incorrect default timeout")

	config.WithTimeout(time.Hour)

	assert.Equal(t, time.Hour, config.timeout, "Incorrect timeout in struct")
	assert.Equal(t, time.Hour, config.Timeout(), "Incorrect timeout via function")
}

func TestWithMiddleware(t *testing.T) {
	config := defaultConfig()

	assert.Equal(t, 0, len(config.middlewares), "Incorrect default midlewares")

	config.WithMiddleware(testMiddleware)

	assert.Equal(t, 1, len(config.middlewares), "Incorrect middleware count")
	assert.Equal(t, 1, len(config.Middlewares()), "Incorrect middleware count via function")
	assert.NotNil(t, config.middlewares[0], "Incorrect middleware.")
	assert.NotNil(t, config.Middlewares()[0], "Incorrect middleware via function")
}

func TestWithMiddlewares(t *testing.T) {
	config := defaultConfig()

	assert.Equal(t, 0, len(config.middlewares), "Incorrect default midlewares")

	config.WithMiddlewares([]MiddlewareFunc{testMiddleware, testMiddlewareTwo})

	assert.Equal(t, 2, len(config.middlewares), "Incorrect middleware count")
	assert.Equal(t, 2, len(config.Middlewares()), "Incorrect middleware count via function")
	assert.NotNil(t, config.middlewares[0], "Incorrect middleware zero.")
	assert.NotNil(t, config.Middlewares()[0], "Incorrect middleware zero via function")
	assert.NotNil(t, config.middlewares[1], "Incorrect middleware one.")
	assert.NotNil(t, config.Middlewares()[1], "Incorrect middleware one via function")
}
