package routem

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type (
	// typed with a private type to prevent override
	contextKey int

	requestData struct {
		request  *http.Request
		response http.ResponseWriter
		params   Params
	}
)

const (
	// private to prevent requestData override
	requestKey contextKey = iota
)

func newContext(timeout time.Duration, request *http.Request, response http.ResponseWriter, params Params) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	data := requestData{
		request:  request,
		response: response,
		params:   params,
	}

	ctx = context.WithValue(ctx, requestKey, data)

	return ctx, cancel
}

// RequestFromContext returns the *http.Request stored in this Context
func RequestFromContext(c context.Context) *http.Request {
	val, ok := c.Value(requestKey).(requestData)
	if !ok {
		contextPanic()
	}
	return val.request
}

// ResponseWriterFromContext returns the http.ResponseWriter stored in this Context
func ResponseWriterFromContext(c context.Context) http.ResponseWriter {
	val, ok := c.Value(requestKey).(requestData)
	if !ok {
		contextPanic()
	}
	return val.response
}

// ParamsFromContext returns the Params stored in this Context
func ParamsFromContext(c context.Context) Params {
	val, ok := c.Value(requestKey).(requestData)
	if !ok {
		contextPanic()
	}
	return val.params
}

func contextPanic() {
	panic("Routem: WTF?! Missing request data in context!")
}
