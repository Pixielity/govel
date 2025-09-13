package main

import (
	"context"
	"fmt"
	"log"
	"time"

	webserver "govel/packages/new/webserver/src"
	"govel/packages/new/webserver/src/builders"
	"govel/packages/new/webserver/src/enums"
	"govel/packages/new/webserver/src/interfaces"
	"govel/packages/new/webserver/src/logging"
	"govel/packages/new/webserver/src/types"
)

// Basic webserver example showing unified API usage
// This example demonstrates core functionality using the builder pattern
func main() {
	fmt.Println("Starting Basic Webserver Example...")

	// Create webserver using builder pattern with GoFiber engine
	server := builders.Configure().
		WithEngine(enums.GoFiber).
		WithPort(8080).
		WithHost("localhost").
		Build()

	// Basic GET route
	server.Get("/", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]string{
			"message": "Hello from Govel Webserver!",
			"version": "1.0.0",
			"engine":  "GoFiber",
		})
	}))

	// GET route with path parameter
	server.Get("/users/:id", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		id := req.Param("id")
		return webserver.NewResponse().Json(map[string]interface{}{
			"user_id": id,
			"name":    fmt.Sprintf("User %s", id),
			"email":   fmt.Sprintf("user%s@example.com", id),
		})
	}))

	// POST route with JSON body
	server.Post("/users", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		name := req.Input("name", "Unknown")
		email := req.Input("email", "")

		if email == "" {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]string{
					"error": "Email is required",
				})
		}

		return webserver.NewResponse().
			Status(201).
			Json(map[string]interface{}{
				"id":      "123",
				"name":    name,
				"email":   email,
				"created": time.Now().Format(time.RFC3339),
			})
	}))

	// PUT route for updates
	server.Put("/users/:id", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		id := req.Param("id")
		name := req.Input("name")
		email := req.Input("email")

		updateData := map[string]interface{}{
			"id":      id,
			"updated": time.Now().Format(time.RFC3339),
		}

		if name != "" {
			updateData["name"] = name
		}
		if email != "" {
			updateData["email"] = email
		}

		return webserver.NewResponse().Json(updateData)
	}))

	// DELETE route
	server.Delete("/users/:id", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		id := req.Param("id")
		return webserver.NewResponse().Json(map[string]interface{}{
			"message":    fmt.Sprintf("User %s deleted successfully", id),
			"deleted_at": time.Now().Format(time.RFC3339),
		})
	}))

	// Query parameters example
	server.Get("/search", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		query := req.Query("q", "")
		limit := req.QueryInt("limit", 10)
		page := req.QueryInt("page", 1)
		includeDeleted := req.QueryBool("include_deleted", false)

		if query == "" {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]string{
					"error": "Query parameter 'q' is required",
				})
		}

		return webserver.NewResponse().Json(map[string]interface{}{
			"query":           query,
			"limit":           limit,
			"page":            page,
			"include_deleted": includeDeleted,
			"results":         []string{"result1", "result2", "result3"},
		})
	}))

	// Health check endpoint
	server.Get("/health", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(time.Now()).String(),
		})
	}))

	// Route groups example
	server.Group("/api/v1", func(api interfaces.WebserverInterface) {
		// API endpoints under /api/v1
		api.Get("/version", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			return webserver.NewResponse().Json(map[string]string{
				"api_version": "v1",
				"server":      "govel",
			})
		}))

		api.Get("/status", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			return webserver.NewResponse().Json(map[string]interface{}{
				"api":       "v1",
				"status":    "operational",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}))
	})

	// Nested route groups
	server.Group("/admin", func(admin interfaces.WebserverInterface) {
		admin.Group("/users", func(users interfaces.WebserverInterface) {
			users.Get("/", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
				return webserver.NewResponse().Json(map[string]string{
					"message": "Admin users list",
				})
			}))

			users.Get("/:id", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
				id := req.Param("id")
				return webserver.NewResponse().Json(map[string]interface{}{
					"message": "Admin user details",
					"user_id": id,
				})
			}))
		})
	})

	// Content type examples
	server.Get("/text", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Text("Hello, World! This is plain text.")
	}))

	server.Get("/html", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		html := `
<!DOCTYPE html>
<html>
<head>
	<title>Govel Webserver</title>
</head>
<body>
	<h1>Welcome to Govel!</h1>
	<p>This is an HTML response from the Govel webserver.</p>
</body>
</html>`
		return webserver.NewResponse().HTML(html)
	}))

	// Header manipulation example
	server.Get("/headers", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		userAgent := req.Header("User-Agent", "Unknown")
		host := req.Header("Host", "Unknown")

		return webserver.NewResponse().
			Header("X-Custom-Header", "Govel-Server").
			Header("X-Request-ID", fmt.Sprintf("%d", time.Now().Unix())).
			Json(map[string]interface{}{
				"your_user_agent":    userAgent,
				"your_host":          host,
				"custom_headers_set": []string{"X-Custom-Header", "X-Request-ID"},
			})
	}))

	// Redirect example
	server.Get("/redirect", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Redirect("/", 302)
	}))

	// Start the server
	fmt.Println("Starting Basic Webserver Example...")
	fmt.Println("")
	
	// Display all registered routes with clickable URLs
	logging.DisplayRoutesClickable(server, "localhost", 8080)
	
	fmt.Println("Press Ctrl+C to stop the server")

	// Handle graceful shutdown
	go func() {
		// In a real application, you'd set up signal handling here
		time.Sleep(30 * time.Second) // Simulate running for 30 seconds
		fmt.Println("\nShutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Start listening
	if err := server.Listen(":8080"); err != nil {
		log.Printf("Server failed to start: %v", err)
	}

	fmt.Println("Server stopped.")
}
