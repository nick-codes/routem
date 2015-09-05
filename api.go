package routem

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// Constants for various HTTP Method strings
const (
	Connect Method = "CONNECT" // Connect represents the CONNECT HTTP Method.
	Delete  Method = "DELETE"  // Delete represents the DELETE HTTP Method
	Get     Method = "GET"     // Get represents the GET HTTP Method
	Head    Method = "HEAD"    // Head represents the HEAD HTTP Method
	Options Method = "OPTIONS" // Options represents the OPTIONS HTTP Method
	Patch   Method = "PATCH"   // Patch represents the PATCH HTTP Method
	Put     Method = "PUT"     // Put represents the PUT HTTP Method
	Post    Method = "POST"    // Post represents the POST HTTP Method
	Trace   Method = "TRACE"   // Trace represents the TRACE HTTP Method
)

// The DefaultTimeout used by a Router
const (
	DefaultTimeout time.Duration = 2 * time.Second // Two seconds should be enough, right?
)

type (
	// Method is the type for HTTP Methods
	Method string

	// RouteConfigurator is responsible for configuring a route
	// and reporting that configuration. This includes
	// the error handler, timeout and middleware stack
	// for a given route.
	RouteConfigurator interface {
		WithErrorHandler(ErrorHandlerFunc) RouteConfigurator
		WithTimeout(time.Duration) RouteConfigurator
		WithMiddleware(MiddlewareFunc) RouteConfigurator
		WithMiddlewares([]MiddlewareFunc) RouteConfigurator

		ErrorHandler() ErrorHandlerFunc
		Timeout() time.Duration
		Middlewares() []MiddlewareFunc
	}

	// A Routable is a Group or a Route which can be configured
	// and has a path.
	Routable interface {
		RouteConfigurator

		Path() string
	}

	// A RouteCreator knows how to construct routes and Groups,
	// and returns and Routables constructed by the creator.
	//
	// Routes returns all the Routables construted by this creator.
	//
	// With() constructs a new Route with the requested methods, path
	// and handler. If the RouteCreator is a Group then the path is
	// nested inside the prefix for the group. The Route will inherit
	// all configuration from the RouteCreator this is called on at
	// the time of the call. Further reconfiguration of the
	// RouteCreator will not effect the configuration of this route.
	//
	// WithHTTP() constructs a new Route with the requested methods,
	// path and http.Handler. This method is a convenience for
	// projects moving towards the use of Routem. Endpoints added in
	// this way still pay the Context construction cost, but for users
	// with legacy handlers this can be used to ease the transition
	// towards the use of Context in the stack since both older
	// http.Handlers and Routem Handlers can be mixed under
	// Routem. The Route will inherit all configuration from the
	// RouteCreator this is called on at the time of the call.
	// Further reconfiguration of the RouteCreator will not effect the
	// configuration of this route.
	//
	// WithGroup() Constructs a new Group with the given path.  This
	// Group inherits any configuration from the RouteCreator
	// constructing it and as with Routes further changes to the
	// configuration of the RouteCreator will not be inherited by this
	// Group. However Routes created by the group will inherit the
	// Group's configuration.
	//
	// The rest of the interface is syntactic sugar to make code more
	// readable.
	RouteCreator interface {
		Routes() []Routable

		With([]Method, string, HandlerFunc) Route
		WithHTTP([]Method, string, http.Handler) Route
		WithGroup(string) Group

		Noop(string, HandlerFunc) Route
		Connect(string, HandlerFunc) Route
		Delete(string, HandlerFunc) Route
		Get(string, HandlerFunc) Route
		Head(string, HandlerFunc) Route
		Options(string, HandlerFunc) Route
		Patch(string, HandlerFunc) Route
		Put(string, HandlerFunc) Route
		Post(string, HandlerFunc) Route
		Trace(string, HandlerFunc) Route
		Crud(string, HandlerFunc) Route
		Any(string, HandlerFunc) Route

		NoopHTTP(string, http.Handler) Route
		ConnectHTTP(string, http.Handler) Route
		DeleteHTTP(string, http.Handler) Route
		GetHTTP(string, http.Handler) Route
		HeadHTTP(string, http.Handler) Route
		OptionsHTTP(string, http.Handler) Route
		PatchHTTP(string, http.Handler) Route
		PutHTTP(string, http.Handler) Route
		PostHTTP(string, http.Handler) Route
		TraceHTTP(string, http.Handler) Route
		CrudHTTP(string, http.Handler) Route
		AnyHTTP(string, http.Handler) Route
	}

	// Runnable provides an interface for things which can be served.
	//
	// Run() serves the configured Routes. Note that Run may be called
	// multiple times to serve the same set of routes on multiple
	// addresses. Further configuration after the call to Run do not
	// effect the served routes. This makes it easy to server all the
	// same routes plus additional routes on an internal port.
	//
	// RunTLS() serves the configured Routes using the passed TLS
	// Configuration. Note that Run may be called multiple times to
	// serve the same set of routes on multiple addresses. Further
	// configuration after the call to Run do not effect the served
	// routes.
	Runnable interface {
		Run(address string) error
		RunTLS(address string, cert string, key string) error
	}

	// A Router holds default configuration for all Routes and Groups
	// and knows how to create them and run the final Server.
	Router interface {
		RouteConfigurator
		RouteCreator
		Runnable
	}

	// A Group provides a container with a prefix for
	// routes.
	Group interface {
		Routable
		RouteCreator
	}

	// A Route has methods, and a handler.
	Route interface {
		Routable

		Methods() []Method
		Handler() HandlerFunc
	}

	// HTTPError encapsulates an error with an HTTP Result code.
	HTTPError interface {
		error
		Code() int
	}

	// HandlerFunc is a Routem Handler that takes a context and
	// returns an HTTPError.  If the function returns a HTTPError
	// Routem will write the error headers to the response. The
	// HTTPError will be passed to the Error Handler function
	// configured for the route.
	HandlerFunc func(context.Context) HTTPError

	// ErrorHandlerFunc is an error handler function that handles an
	// HTTPError returned by a HandlerFunc. If an ErrorHandlerFunc
	// returns an error Routem will write an Internal Server Error
	// header.
	ErrorHandlerFunc func(HTTPError, context.Context) error

	// MiddlewareFunc provides an easy way to alter or validate an
	// inbound request. Middleware can be stacked in front of a Route
	// by configuring it as part of the routes stack.  The Middleware
	// has the ability to add information to the Context it passes on
	// to the next handler it calls in the stack.  a MiddlewareFunc is
	// responsible for wrapping the HandlerFunc passed in and calling
	// that HandlerFunc if execution of the request should continue.
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// Params provides mapping of path parameter names to the values
	// found in the request URI.
	Params map[string]string
)

// Conveniences for those who prefer to use the With() and WithHTTP()
// interfaces instead of the syntactic sugar equivalents. NoMethod can
// be used to temporarily disable a Route without having to delete the
// associated setup code.
//
// These are also used internally to construct the equivalent syntactic sugar methods.
var (
	NoMethod      = []Method{}
	ConnectMethod = []Method{Connect}
	DeleteMethod  = []Method{Delete}
	GetMethod     = []Method{Get}
	HeadMethod    = []Method{Head}
	OptionsMethod = []Method{Options}
	PatchMethod   = []Method{Patch}
	PutMethod     = []Method{Put}
	PostMethod    = []Method{Post}
	TraceMethod   = []Method{Trace}
	CrudMethod    = []Method{Delete, Get, Put, Patch, Post}
	AnyMethod     = []Method{Connect, Delete, Get, Head, Options, Patch, Put, Post, Trace}
)
