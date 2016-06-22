package routem

import (
	"net/http"
	"time"

	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func newContext(timeout time.Duration, request *http.Request, response http.ResponseWriter, params Params) (context.Context, context.CancelFunc) {
	return NewRequestContext(context.Background(), timeout, request, response, params)
}

func setupContextTest(t *testing.T) (http.ResponseWriter, *http.Request, Params) {
	response := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "http://test.om/", nil)

	if err != nil {
		t.Fatalf("Unable to build a request.")
	}

	params := make(Params)

	return response, request, params
}

func TestNewContext(t *testing.T) {
	response, request, params := setupContextTest(t)

	ctx, cancel := newContext(DefaultTimeout,
		request, response, params)

	assert.NotNil(t, ctx, "Context null")
	assert.NotNil(t, cancel, "Cancel null")
}

func TestResponseWriterFromContext(t *testing.T) {
	response, request, params := setupContextTest(t)

	ctx, _ := newContext(DefaultTimeout,
		request, response, params)

	testResponse := ResponseWriterFromContext(ctx)

	assert.Equal(t, response, testResponse, "Responses not equal")

	assert.Panics(t, func() { ResponseWriterFromContext(context.Background()) }, "No Panic on Empty Context")
}

func TestRequestFromContext(t *testing.T) {
	response, request, params := setupContextTest(t)

	ctx, _ := newContext(DefaultTimeout,
		request, response, params)

	testRequest := RequestFromContext(ctx)

	assert.Equal(t, request, testRequest, "Request not equal")

	assert.Panics(t, func() { RequestFromContext(context.Background()) }, "No Panic on Empty Context")
}

func TestParamsFromContext(t *testing.T) {
	response, request, params := setupContextTest(t)

	ctx, _ := newContext(DefaultTimeout,
		request, response, params)

	testParams := ParamsFromContext(ctx)

	assert.Equal(t, params, testParams, "Params not equal")

	assert.Panics(t, func() { ParamsFromContext(context.Background()) }, "No Panic on Empty Context")
}
