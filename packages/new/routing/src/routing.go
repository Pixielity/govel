// Package routing - Routing-related.
package routing

import (
	"fmt"
	"govel/packages/new/webserver/src/enums"
	"govel/packages/new/webserver/src/interfaces"
)

// Route represents a single route definition with all its associated metadata.
// This type encapsulates all the information needed to register and handle a route.
type Route struct {
	// Method is the HTTP method for this route (GET, POST, etc.)
	Method enums.HTTPMethod

	// Path is the URL path pattern for this route (e.g., "/users/:id")
	Path string

	// Handler is the handler that will process requests to this route
	Handler interfaces.HandlerInterface

	// Middleware contains route-specific middleware
	Middleware []interfaces.MiddlewareInterface

	// Name is an optional name for the route (useful for URL generation)
	Name string

	// Metadata contains arbitrary data associated with the route
	Metadata map[string]interface{}
}

// NewRoute creates a new Route with the specified method, path, and handler.
//
// Parameters:
//
//	method: The HTTP method for the route
//	path: The URL path pattern
//	handler: The handler for processing requests
//
// Returns:
//
//	*Route: A new route instance
//
// Example:
//
//	route := NewRoute(enums.GET, "/users/:id", getUserHandler)
func NewRoute(method enums.HTTPMethod, path string, handler interfaces.HandlerInterface) *Route {
	return &Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: make([]interfaces.MiddlewareInterface, 0),
		Metadata:   make(map[string]interface{}),
	}
}

// WithName sets the name for the route.
//
// Parameters:
//
//	name: The name to assign to the route
//
// Returns:
//
//	*Route: The route instance for method chaining
func (r *Route) WithName(name string) *Route {
	r.Name = name
	return r
}

// WithMiddleware adds middleware to the route.
//
// Parameters:
//
//	middleware: One or more middleware implementations to add
//
// Returns:
//
//	*Route: The route instance for method chaining
func (r *Route) WithMiddleware(middleware ...interfaces.MiddlewareInterface) *Route {
	r.Middleware = append(r.Middleware, middleware...)
	return r
}

// WithMetadata sets metadata for the route.
//
// Parameters:
//
//	key: The metadata key
//	value: The metadata value
//
// Returns:
//
//	*Route: The route instance for method chaining
func (r *Route) WithMetadata(key string, value interface{}) *Route {
	r.Metadata[key] = value
	return r
}

// GetMetadata retrieves metadata by key.
//
// Parameters:
//
//	key: The metadata key to retrieve
//
// Returns:
//
//	interface{}: The metadata value, or nil if not found
//	bool: True if the key exists, false otherwise
func (r *Route) GetMetadata(key string) (interface{}, bool) {
	value, exists := r.Metadata[key]
	return value, exists
}

// String returns a string representation of the route.
// This is useful for debugging and logging.
//
// Returns:
//
//	string: String representation of the route
func (r *Route) String() string {
	name := r.Name
	if name == "" {
		name = "<unnamed>"
	}
	return fmt.Sprintf("%s %s [%s]", r.Method.String(), r.Path, name)
}

// RouteGroup represents a group of routes that share common properties.
// This is useful for organizing routes and applying common middleware or prefixes.
type RouteGroup struct {
	// Prefix is the URL prefix for all routes in this group
	Prefix string

	// Routes contains all routes in this group
	Routes []*Route

	// Middleware contains middleware that applies to all routes in the group
	Middleware []interfaces.MiddlewareInterface

	// Name is an optional name for the route group
	Name string

	// Metadata contains arbitrary data associated with the group
	Metadata map[string]interface{}
}

// NewRouteGroup creates a new route group with the specified prefix.
//
// Parameters:
//
//	prefix: The URL prefix for routes in this group
//
// Returns:
//
//	*RouteGroup: A new route group instance
//
// Example:
//
//	apiGroup := NewRouteGroup("/api/v1")
func NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		Prefix:     prefix,
		Routes:     make([]*Route, 0),
		Middleware: make([]interfaces.MiddlewareInterface, 0),
		Metadata:   make(map[string]interface{}),
	}
}

// AddRoute adds a route to the group.
//
// Parameters:
//
//	route: The route to add to the group
//
// Returns:
//
//	*RouteGroup: The route group instance for method chaining
func (rg *RouteGroup) AddRoute(route *Route) *RouteGroup {
	rg.Routes = append(rg.Routes, route)
	return rg
}

// WithName sets the name for the route group.
//
// Parameters:
//
//	name: The name to assign to the group
//
// Returns:
//
//	*RouteGroup: The route group instance for method chaining
func (rg *RouteGroup) WithName(name string) *RouteGroup {
	rg.Name = name
	return rg
}

// WithMiddleware adds middleware to the route group.
//
// Parameters:
//
//	middleware: One or more middleware implementations to add
//
// Returns:
//
//	*RouteGroup: The route group instance for method chaining
func (rg *RouteGroup) WithMiddleware(middleware ...interfaces.MiddlewareInterface) *RouteGroup {
	rg.Middleware = append(rg.Middleware, middleware...)
	return rg
}

// RouteCollection represents a collection of routes and route groups.
// This provides utilities for managing and querying routes.
type RouteCollection struct {
	// routes stores individual routes
	routes []*Route

	// groups stores route groups
	groups []*RouteGroup
}

// NewRouteCollection creates a new empty route collection.
//
// Returns:
//
//	*RouteCollection: A new route collection instance
func NewRouteCollection() *RouteCollection {
	return &RouteCollection{
		routes: make([]*Route, 0),
		groups: make([]*RouteGroup, 0),
	}
}

// AddRoute adds a route to the collection.
//
// Parameters:
//
//	route: The route to add
//
// Returns:
//
//	*RouteCollection: The collection instance for method chaining
func (rc *RouteCollection) AddRoute(route *Route) *RouteCollection {
	rc.routes = append(rc.routes, route)
	return rc
}

// AddGroup adds a route group to the collection.
//
// Parameters:
//
//	group: The route group to add
//
// Returns:
//
//	*RouteCollection: The collection instance for method chaining
func (rc *RouteCollection) AddGroup(group *RouteGroup) *RouteCollection {
	rc.groups = append(rc.groups, group)
	return rc
}

// GetAllRoutes returns all routes including those in groups.
//
// Returns:
//
//	[]*Route: All routes in the collection
func (rc *RouteCollection) GetAllRoutes() []*Route {
	allRoutes := make([]*Route, len(rc.routes))
	copy(allRoutes, rc.routes)

	// Add routes from groups
	for _, group := range rc.groups {
		allRoutes = append(allRoutes, group.Routes...)
	}

	return allRoutes
}

// FindByName finds a route by name.
//
// Parameters:
//
//	name: The route name to search for
//
// Returns:
//
//	*Route: The found route, or nil if not found
func (rc *RouteCollection) FindByName(name string) *Route {
	for _, route := range rc.GetAllRoutes() {
		if route.Name == name {
			return route
		}
	}
	return nil
}

// FilterByMethod returns all routes that match the specified HTTP method.
//
// Parameters:
//
//	method: The HTTP method to filter by
//
// Returns:
//
//	[]*Route: Routes that match the method
func (rc *RouteCollection) FilterByMethod(method enums.HTTPMethod) []*Route {
	var filtered []*Route
	for _, route := range rc.GetAllRoutes() {
		if route.Method == method {
			filtered = append(filtered, route)
		}
	}
	return filtered
}

// Count returns the total number of routes in the collection.
//
// Returns:
//
//	int: The total number of routes
func (rc *RouteCollection) Count() int {
	return len(rc.GetAllRoutes())
}
