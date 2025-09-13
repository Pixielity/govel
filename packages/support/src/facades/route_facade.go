package facades

import (
	routeInterfaces "govel/types/interfaces/route"
	facade "govel/support"
)

// Route provides a clean, static-like interface to the application's routing service.
//
// This facade implements the facade pattern, providing global access to the routing
// service configured in the dependency injection container. It offers a Laravel-style
// API for route registration, URL generation, middleware management, and route
// parameter handling with automatic service resolution and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved routing service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent route operations across goroutines
//   - Supports dynamic route registration and URL generation
//   - Built-in middleware pipeline and parameter binding integration
//
// Behavior:
//   - First call: Resolves route service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if route service cannot be resolved (fail-fast behavior)
//   - Automatically handles route compilation, matching, and URL generation
//
// Returns:
//   - RouteInterface: The application's routing service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "route" service is not registered in the container
//   - If the resolved service doesn't implement RouteInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Route() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Route registration and matching are thread-safe and consistent
//
// Usage Examples:
//
//	// Basic route registration
//	facades.Route().Get("/users", userController.Index)
//	facades.Route().Post("/users", userController.Create)
//	facades.Route().Put("/users/{id}", userController.Update)
//	facades.Route().Delete("/users/{id}", userController.Delete)
//
//	// Route with middleware
//	facades.Route().Get("/admin", adminController.Dashboard).Middleware("auth", "admin")
//	facades.Route().Post("/api/data", apiController.Store).Middleware("api", "throttle:60,1")
//
//	// Route groups for organization
//	facades.Route().Group(func(r RouteRegistrar) {
//	    r.Get("/profile", userController.Profile)
//	    r.Post("/settings", userController.UpdateSettings)
//	}).Prefix("/user").Middleware("auth")
//
//	facades.Route().Group(func(r RouteRegistrar) {
//	    r.Get("/users", adminController.Users)
//	    r.Get("/reports", adminController.Reports)
//	}).Prefix("/admin").Middleware("auth", "admin")
//
//	// API route groups
//	facades.Route().Group(func(r RouteRegistrar) {
//	    r.Get("/posts", apiController.GetPosts)
//	    r.Post("/posts", apiController.CreatePost)
//	    r.Get("/posts/{id}", apiController.GetPost)
//	    r.Put("/posts/{id}", apiController.UpdatePost)
//	    r.Delete("/posts/{id}", apiController.DeletePost)
//	}).Prefix("/api/v1").Middleware("api")
//
//	// Named routes for URL generation
//	facades.Route().Get("/user/{id}/profile", userController.Profile).Name("user.profile")
//	facades.Route().Post("/posts/{id}/comments", commentController.Store).Name("comments.store")
//
//	// URL generation from named routes
//	profileURL := facades.Route().URL("user.profile", map[string]interface{}{"id": 123})
//	commentURL := facades.Route().URL("comments.store", map[string]interface{}{"id": 456})
//
//	// Route parameters and constraints
//	facades.Route().Get("/users/{id}", userController.Show).Where("id", "[0-9]+")
//	facades.Route().Get("/posts/{slug}", postController.Show).Where("slug", "[a-z0-9-]+")
//	facades.Route().Get("/categories/{category}/{subcategory?}", categoryController.Show)
//
//	// Route model binding
//	facades.Route().Get("/users/{user}", userController.Show) // Automatically injects User model
//	facades.Route().Get("/posts/{post}/comments/{comment}", commentController.Show)
//
//	// Resource routes (RESTful)
//	facades.Route().Resource("posts", postController)
//	facades.Route().Resource("users.comments", commentController) // Nested resources
//	facades.Route().APIResource("api/posts", apiPostController) // API-only routes
//
//	// Subdomain routing
//	facades.Route().Group(func(r RouteRegistrar) {
//	    r.Get("/", adminController.Dashboard)
//	    r.Get("/users", adminController.Users)
//	}).Domain("admin.{domain}")
//
//	// Route caching for production performance
//	if facades.App().IsProduction() {
//	    facades.Route().Cache() // Cache compiled routes
//	}
//
//	// Route fallbacks and error handling
//	facades.Route().Fallback(func(c Context) {
//	    c.JSON(404, map[string]string{"error": "Not Found"})
//	})
//
//	// Custom route matching
//	facades.Route().Match(["GET", "POST"], "/contact", contactController.Handle)
//	facades.Route().Any("/webhook", webhookController.Handle)
//
//	// Route information and introspection
//	allRoutes := facades.Route().GetRoutes()
//	routeExists := facades.Route().HasRoute("user.profile")
//	currentRoute := facades.Route().Current()
//	routeName := facades.Route().CurrentRouteName()
//
// Advanced Routing Patterns:
//
//	// Route service providers for organization
//	type APIRoutesServiceProvider struct {}
//
//	func (p *APIRoutesServiceProvider) Boot() {
//	    facades.Route().Group(func(r RouteRegistrar) {
//	        // User management routes
//	        r.Get("/users", userController.Index)
//	        r.Post("/users", userController.Store)
//	        r.Get("/users/{user}", userController.Show)
//	        r.Put("/users/{user}", userController.Update)
//	        r.Delete("/users/{user}", userController.Destroy)
//
//	        // Authentication routes
//	        r.Post("/login", authController.Login)
//	        r.Post("/logout", authController.Logout)
//	        r.Post("/refresh", authController.Refresh)
//	    }).Prefix("/api/v1").Middleware("api")
//	}
//
//	// Conditional route registration
//	func RegisterRoutes() {
//	    // Always available routes
//	    facades.Route().Get("/", homeController.Index)
//	    facades.Route().Get("/about", pageController.About)
//
//	    // Environment-specific routes
//	    if facades.App().IsLocal() {
//	        facades.Route().Get("/debug", debugController.Index)
//	        facades.Route().Get("/profiler", profilerController.Show)
//	    }
//
//	    // Feature flag routes
//	    if facades.Config().GetBool("features.admin_panel") {
//	        facades.Route().Group(func(r RouteRegistrar) {
//	            r.Get("/dashboard", adminController.Dashboard)
//	            r.Resource("users", adminUserController)
//	        }).Prefix("/admin").Middleware("auth", "admin")
//	    }
//	}
//
//	// Route middleware stacks
//	func SetupMiddleware() {
//	    // Global middleware
//	    facades.Route().GlobalMiddleware("cors")
//	    facades.Route().GlobalMiddleware("session")
//
//	    // Route-specific middleware groups
//	    facades.Route().MiddlewareGroup("web", []string{"session", "csrf", "throttle:60,1"})
//	    facades.Route().MiddlewareGroup("api", []string{"api", "throttle:60,1", "auth:api"})
//	}
//
//	// Dynamic route registration
//	func RegisterDynamicRoutes() {
//	    // Load routes from database or configuration
//	    dynamicRoutes := loadDynamicRoutes()
//	    for _, route := range dynamicRoutes {
//	        facades.Route().Get(route.Path, route.Handler).Name(route.Name)
//	    }
//
//	    // Plugin-based routes
//	    plugins := facades.App().GetPlugins()
//	    for _, plugin := range plugins {
//	        plugin.RegisterRoutes(facades.Route())
//	    }
//	}
//
// Best Practices:
//   - Use route groups to organize related routes and apply common middleware
//   - Name important routes for URL generation and easier maintenance
//   - Use route model binding to automatically inject models
//   - Apply appropriate middleware for authentication, authorization, and rate limiting
//   - Use resource routes for RESTful APIs
//   - Cache routes in production for better performance
//   - Use route parameters with constraints for input validation
//   - Organize routes in service providers for large applications
//
// Route Organization:
//  1. Define route groups by feature or domain
//  2. Apply middleware at the group level when possible
//  3. Use descriptive route names for important routes
//  4. Keep route definitions close to related controllers
//  5. Use separate route files for different areas (web, api, admin)
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume routing always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	router, err := facade.TryResolve[RouteInterface]("route")
//	if err != nil {
//	    // Handle routing unavailability gracefully
//	    return fmt.Errorf("routing unavailable: %w", err)
//	}
//	router.Get("/path", handler)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestRouting(t *testing.T) {
//	    // Create a test router
//	    testRouter := &TestRouter{
//	        routes: make(map[string]Route),
//	    }
//
//	    // Swap the real router with test router
//	    restore := support.SwapService("route", testRouter)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Route() returns testRouter
//	    facades.Route().Get("/test", testHandler)
//
//	    // Verify route registration
//	    assert.True(t, testRouter.HasRoute("/test"))
//	    assert.Equal(t, "GET", testRouter.GetRoute("/test").Method)
//	}
//
// Container Configuration:
// Ensure the route service is properly configured in your container:
//
//	// Example router registration
//	container.Singleton("route", func() interface{} {
//	    config := router.Config{
//	        // Router configuration
//	        CaseSensitive:     false,
//	        StrictSlash:       true,
//	        UseEncodedPath:    false,
//	        HandleMethodNotAllowed: true,
//
//	        // Middleware configuration
//	        GlobalMiddleware: []string{"cors", "session"},
//	        MiddlewareGroups: map[string][]string{
//	            "web": {"session", "csrf", "throttle:60,1"},
//	            "api": {"api", "throttle:60,1", "auth:api"},
//	        },
//
//	        // Route caching
//	        CacheRoutes: facades.App().IsProduction(),
//	        CachePath:   facades.App().StoragePath("cache", "routes.cache"),
//
//	        // Route model binding
//	        ModelBindings: map[string]interface{}{
//	            "user": &User{},
//	            "post": &Post{},
//	        },
//	    }
//
//	    return router.NewRouter(config)
//	})
func Route() routeInterfaces.RouteInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "route" service from the dependency injection container
	// - Performs type assertion to RouteInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[routeInterfaces.RouteInterface](routeInterfaces.ROUTE_TOKEN)
}

// RouteWithError provides error-safe access to the routing service.
//
// This function offers the same functionality as Route() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle routing unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Route() but with error handling.
//
// Returns:
//   - RouteInterface: The resolved route instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement RouteInterface
//
// Usage Examples:
//
//	// Basic error-safe routing
//	router, err := facades.RouteWithError()
//	if err != nil {
//	    log.Printf("Router unavailable: %v", err)
//	    return fmt.Errorf("routing service not available")
//	}
//	router.Get("/test", testHandler)
//
//	// Conditional route registration
//	if router, err := facades.RouteWithError(); err == nil {
//	    router.Get("/optional", optionalHandler)
//	}
func RouteWithError() (routeInterfaces.RouteInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "route" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[routeInterfaces.RouteInterface](routeInterfaces.ROUTE_TOKEN)
}
