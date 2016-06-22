package trie

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nick-codes/routem"

	"golang.org/x/net/context"
)

// Compile time type assertions
var _ http.Handler = &rootNode{}
var _ routem.HandlerFactory = &factory{}

type routeInfo struct {
	route   routem.Route
	params  map[int]string
	handler routem.HandlerFunc
}

type rootNode struct {
	ctx          context.Context
	errorHandler routem.ErrorHandlerFunc
	node
}

type node struct {
	path     string
	routes   map[routem.Method]*routeInfo
	children map[string]*node
}

type factory struct {
	ctx          context.Context
	errorHandler routem.ErrorHandlerFunc
}

// Constructs a new handler factory which uses a trie data structure
// to quickly look up routes.
//
// All routes will be passed a context
// derived from the context passed to the factory. If no context is
// passed then context.Background() is used as the root context.
//
// If an ErrorHandlerFunc is provided and the route does not have a route
// specific error handler that handler will be called if a route
// returns an error. Otherwise a 500 error will be returned to the client.
func NewHandlerFactory(ctx context.Context, errorHandler routem.ErrorHandlerFunc) routem.HandlerFactory {
	if ctx == nil {
		ctx = context.Background()
	}
	return &factory{
		ctx:          ctx,
		errorHandler: errorHandler,
	}
}

func (f *factory) Handler(routes []routem.Route) (http.Handler, error) {

	if routes == nil || len(routes) == 0 {
		return nil, fmt.Errorf("Received no routes")
	}

	root := &rootNode{
		ctx:          f.ctx,
		errorHandler: f.errorHandler,
		node: node{
			path:     "",
			children: make(map[string]*node),
			routes:   make(map[routem.Method]*routeInfo),
		},
	}

	for _, route := range routes {
		if route == nil {
			return nil, fmt.Errorf("Received a nil route.")
		}

		if len(route.Path()) == 0 {
			return nil, fmt.Errorf("Received a zero length path.")
		}

		if strings.Contains(route.Path(), "//") {
			return nil, fmt.Errorf("Route contains an invalid path: %s", route.Path())
		}

		if !strings.HasPrefix(route.Path(), "/") {
			return nil, fmt.Errorf("Route does not begin with a slash: %s", route.Path())
		}

		parts := strings.Split(route.Path(), "/")
		params := make(map[string]struct{}, len(parts))

		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				_, exists := params[part]
				if exists {
					return nil, fmt.Errorf("Route has duplicate parameter: %s", part)
				}
				params[part] = struct{}{}
			}
		}

		inserted, err := root.insert(parts, route, 0, nil)

		if err != nil {
			return nil, err
		}

		// This should never happen
		if !inserted {
			return nil, fmt.Errorf("An unknown error occured.")
		}
	}

	return root, nil
}

func newNode(path string) (*node, error) {
	var ret *node
	var err error

	if strings.HasPrefix(path, ":") {
		paramName := strings.TrimPrefix(path, ":")

		if len(paramName) == 0 {
			err = fmt.Errorf("Found an un-named parameter: %s", path)
		} else {
			ret = &node{
				path:     ":",
				routes:   make(map[routem.Method]*routeInfo),
				children: make(map[string]*node),
			}
		}
	} else {
		ret = &node{
			path:     path,
			routes:   make(map[routem.Method]*routeInfo),
			children: make(map[string]*node),
		}
	}

	return ret, err
}

func (n *node) insert(parts []string, route routem.Route, depth int, params map[int]string) (bool, error) {

	inserted := false
	var err error

	// Is this a parameter segment?
	thisPath := parts[0]
	if strings.HasPrefix(thisPath, ":") {
		if params == nil {
			params = make(map[int]string, len(parts))
		}
		params[depth] = strings.TrimPrefix(thisPath, ":")
		thisPath = ":"
	}

	// Does this belong in this sub-tree?
	if thisPath == n.path {

		// Is this the path leaf?
		if len(parts) == 1 {

			// Do we already have a route here?
			for _, method := range route.Methods() {
				if n.routes[method] != nil {
					err = fmt.Errorf("Duplicate route: %s - %s", route.Path(), n.routes[method].route.Path())
				}
			}

			// No? Then do the insert
			if err == nil {
				// Build the middleware stack
				handler := route.Handler()
				middlewares := route.Middlewares()
				for i := len(middlewares) - 1; i >= 0; i-- {
					handler = middlewares[i](handler)
				}

				// Remember all the info for the route
				info := &routeInfo{
					route:   route,
					params:  params,
					handler: handler,
				}

				// Set it on the various methods
				for _, method := range route.Methods() {
					n.routes[method] = info
				}

				inserted = true
			}
		} else {

			// Check if we can insert in any existing children
			for _, child := range n.children {
				inserted, err = child.insert(parts[1:], route, depth+1, params)
				if inserted || err != nil {
					break
				}
			}

			// Okay, then make a new child and insert
			if !inserted && err == nil {
				newChild, err := newNode(parts[1])

				if err == nil {
					n.children[newChild.path] = newChild
					inserted, err = newChild.insert(parts[1:], route, depth+1, params)
				}
			}

		}
	}

	return inserted, err
}

var routeNotFoundError = routem.NewHTTPError(http.StatusNotFound, fmt.Errorf("No Such Route"))

func (n *node) find(parts []string, method routem.Method) (*routeInfo, routem.HTTPError) {
	var info *routeInfo = nil
	var err routem.HTTPError = nil

	// Did we fish our wish?
	if n.path == ":" || parts[0] == n.path {

		// Did we run out of parts?
		if len(parts) == 1 {
			info = n.routes[method]
		} else {

			// Search all the children
			subParts := parts[1:]
			for _, child := range n.children {
				info, err = child.find(subParts, method)

				// If we found something return it up the stack
				if info != nil {
					break
				}
			}

		}
	}

	if info == nil {
		err = routeNotFoundError
	}

	return info, err
}

func routeParams(route *routeInfo, parts []string) routem.Params {
	var params routem.Params

	if route != nil && route.params != nil {
		params = make(routem.Params, len(parts))

		for index, param := range route.params {
			// Route cannot match unless it is of sufficient length, so we are sure
			// index  < len(parts) at this point.
			params[param] = parts[index]
		}
	}

	return params
}

func (root *rootNode) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")
	routeInfo, err := root.find(parts, routem.Method(request.Method))

	timeout := routem.DefaultTimeout
	if err == nil {
		timeout = routeInfo.route.Timeout()
	}

	ctx, cancel := routem.NewRequestContext(
		root.ctx, timeout, request, response,
		routeParams(routeInfo, parts))

	defer cancel()

	if err == nil {
		complete := make(chan routem.HTTPError)
		go func() {
			complete <- routeInfo.handler(ctx)
		}()

		select {
		case <-ctx.Done():
			err = routem.NewHTTPError(408, fmt.Errorf("Request Timed Out!"))
		case err = <-complete:
		}
	}

	if err != nil {
		var errErr error

		if routeInfo != nil && routeInfo.route.ErrorHandler() != nil {
			errErr = routeInfo.route.ErrorHandler()(err, ctx)
		} else if root.errorHandler != nil {
			errErr = root.errorHandler(err, ctx)
		} else if err == routeNotFoundError {
			http.Error(response, fmt.Sprintf("Route Not Found: %s", errErr), http.StatusNotFound)
		} else {
			errErr = err
		}

		if errErr != nil {
			http.Error(response, fmt.Sprintf("Internal Server Error: %s", errErr), http.StatusInternalServerError)
		}
	}
}
