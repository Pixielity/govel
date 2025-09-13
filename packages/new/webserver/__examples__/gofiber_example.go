package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"govel/packages/new/webserver/src"
	"govel/packages/new/webserver/src/builders"
	"govel/packages/new/webserver/src/enums"
	"govel/packages/new/webserver/src/interfaces"
	"govel/packages/new/webserver/src/logging"
	"govel/packages/new/webserver/src/types"
)
)

// GoFiber-specific example showing webserver usage with GoFiber adapter
// This example highlights performance and fiber idioms
func main() {
	fmt.Println("Starting GoFiber Webserver Example...")

	// Create webserver using GoFiber engine with specific configuration
	server := builders.Configure().
		WithEngine(enums.GoFiber).
		WithPort(8083).
		WithHost("0.0.0.0").
		Set("fiber.prefork", false).
		Set("fiber.case_sensitive", true).
		Set("fiber.strict_routing", false).
		Set("fiber.server_header", "Govel-Fiber").
		Build()

	// Welcome endpoint with Fiber branding
	server.Get("/", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"message": "Welcome to Govel with Fiber!",
			"engine":  "GoFiber",
			"tagline": "Express.js inspired, fast and concise",
			"features": []string{
				"Low memory footprint",
				"Fast routing",
				"Middleware support",
				"Built-in server utilities",
			},
			"documentation": "https://docs.gofiber.io/",
		})
	}))

	// High-performance list endpoint
	server.Get("/fiber/items", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		limit := req.QueryInt("limit", 100)
		if limit > 1000 {
			limit = 1000
		}
		items := make([]map[string]interface{}, 0, limit)
		for i := 1; i <= limit; i++ {
			items = append(items, map[string]interface{}{
				"id":    i,
				"name":  fmt.Sprintf("Item %d", i),
				"value": i * 10,
			})
		}
		return webserver.NewResponse().Json(map[string]interface{}{
			"count":  limit,
			"items":  items,
			"engine": "Fiber",
		})
	}))

	// Path parameter parsing
	server.Get("/fiber/items/:id", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		id := req.ParamInt("id")
		if id <= 0 {
			return webserver.NewResponse().Status(400).Json(map[string]string{"error": "invalid id"})
		}
		return webserver.NewResponse().Json(map[string]interface{}{
			"id":      id,
			"name":    fmt.Sprintf("Item %d", id),
			"created": time.Now().Format(time.RFC3339),
			"engine":  "Fiber",
		})
	}))

	// Form submission example
	server.Post("/fiber/submit", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		name := req.Input("name", "")
		age := req.Input("age", "0")
		ageInt, _ := strconv.Atoi(age)

		if name == "" || ageInt <= 0 {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]interface{}{
					"error":  "name and positive age are required",
					"engine": "Fiber",
				})
		}

		return webserver.NewResponse().Json(map[string]interface{}{
			"message": fmt.Sprintf("Hello %s!", name),
			"age":     ageInt,
			"engine":  "Fiber",
		})
	}))

	// Route groups and nested paths
	server.Group("/api/fiber", func(api interfaces.WebserverInterface) {
		api.Get("/stats", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			return webserver.NewResponse().Json(map[string]interface{}{
				"uptime":   "simulated",
				"requests": 12345,
				"engine":   "Fiber",
			})
		}))

		api.Get("/time", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			return webserver.NewResponse().Json(map[string]string{
				"now": time.Now().Format(time.RFC3339),
			})
		}))
	})

	// Headers and cookies
	server.Get("/fiber/headers", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		ua := req.Header("User-Agent", "unknown")
		return webserver.NewResponse().
			Header("X-Powered-By", "Govel-Fiber").
			Cookie(&http.Cookie{Name: "fiber", Value: "true", Path: "/", Expires: time.Now().Add(24 * time.Hour)}).
			Json(map[string]interface{}{
				"user_agent": ua,
				"has_cookie": req.HasCookie("fiber"),
			})
	}))

	// JSON body parsing (Fiber strength)
	server.Post("/fiber/users", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		type FiberUser struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Age      int    `json:"age"`
			Active   bool   `json:"active"`
			Metadata map[string]interface{} `json:"metadata"`
		}

		var user FiberUser
		if err := req.Json(&user); err != nil {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]interface{}{
					"error": "Invalid JSON payload",
					"details": err.Error(),
					"fiber_feature": "High-speed JSON parsing",
				})
		}

		// Basic validation
		if user.Name == "" || user.Email == "" {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]interface{}{
					"error": "Name and email are required",
					"fiber_performance": "Fast validation with minimal overhead",
				})
		}

		if user.Age < 0 || user.Age > 150 {
			return webserver.NewResponse().
				Status(400).
				Json(map[string]interface{}{
					"error": "Age must be between 0 and 150",
					"received_age": user.Age,
				})
		}

		// Create response with Fiber optimizations
		userResponse := map[string]interface{}{
			"id":       fmt.Sprintf("fiber-user-%d", time.Now().Unix()),
			"name":     user.Name,
			"email":    user.Email,
			"age":      user.Age,
			"active":   user.Active,
			"metadata": user.Metadata,
			"created_at": time.Now().Format(time.RFC3339),
			"engine":   "GoFiber",
			"fiber_features": []string{
				"Zero allocation JSON",
				"Fast struct binding",
				"Minimal memory usage",
				"Express-like middleware",
			},
		}

		return webserver.NewResponse().
			Status(201).
			Header("X-Created-With", "GoFiber").
			Json(userResponse)
	}))

	// Wildcard route matching (Fiber feature)
	server.Get("/fiber/wildcard/*", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		// In Fiber, we'd use c.Params("*") to get the wildcard part
		path := req.Path()
		wildcardPart := ""
		if len(path) > len("/fiber/wildcard/") {
			wildcardPart = path[len("/fiber/wildcard/"):]
		}

		return webserver.NewResponse().Json(map[string]interface{}{
			"message": "Fiber wildcard route matched",
			"wildcard_path": wildcardPart,
			"full_path": path,
			"fiber_feature": "Powerful route matching with wildcards",
			"examples": []string{
				"/fiber/wildcard/anything",
				"/fiber/wildcard/nested/paths",
				"/fiber/wildcard/file.ext",
			},
		})
	}))

	// Multiple parameter route
	server.Get("/fiber/users/:userId/posts/:postId", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		userId := req.Param("userId")
		postId := req.Param("postId")

		// Simulate nested resource access
		post := map[string]interface{}{
			"post_id":    postId,
			"user_id":    userId,
			"title":      fmt.Sprintf("Post %s by User %s", postId, userId),
			"content":    "This is a sample post content processed by Fiber.",
			"created_at": time.Now().Format(time.RFC3339),
			"fiber_optimization": "Fast parameter extraction",
		}

		return webserver.NewResponse().Json(map[string]interface{}{
			"post": post,
			"fiber_features": []string{
				"Multiple path parameters",
				"Fast route matching",
				"Zero allocation routing",
			},
		})
	}))

	// File upload simulation (Fiber multipart)
	server.Post("/fiber/upload", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		// Simulate file upload processing
		filename := req.Input("filename", "upload.txt")
		description := req.Input("description", "")
		category := req.Input("category", "general")

		// Fiber would handle multipart forms very efficiently
		uploadInfo := map[string]interface{}{
			"upload_id":    fmt.Sprintf("fiber-%d", time.Now().Unix()),
			"filename":     filename,
			"description":  description,
			"category":     category,
			"size":         1024 * 512, // Simulated 512KB
			"content_type": "application/octet-stream",
			"uploaded_at":  time.Now().Format(time.RFC3339),
			"engine":       "GoFiber",
			"fiber_advantages": []string{
				"Efficient multipart parsing",
				"Low memory file handling",
				"Built-in file validation",
				"Stream processing support",
			},
		}

		return webserver.NewResponse().
			Status(201).
			Header("Location", fmt.Sprintf("/fiber/files/%s", filename)).
			Json(map[string]interface{}{
				"message": "File uploaded successfully with Fiber!",
				"upload": uploadInfo,
			})
	}))

	// Performance benchmark endpoint
	server.Get("/fiber/benchmark", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		start := time.Now()

		// Simulate some CPU-intensive work
		data := make(map[string]interface{})
		for i := 0; i < 10000; i++ {
			data[fmt.Sprintf("key_%d", i)] = map[string]interface{}{
				"value":     i * 2,
				"squared":   i * i,
				"timestamp": time.Now().UnixNano(),
			}
		}

		processingTime := time.Since(start)

		return webserver.NewResponse().Json(map[string]interface{}{
			"benchmark_results": map[string]interface{}{
				"items_processed":     10000,
				"processing_time_ms":  float64(processingTime.Nanoseconds()) / 1e6,
				"items_per_second":    10000.0 / (float64(processingTime.Nanoseconds()) / 1e9),
				"memory_efficiency":   "Optimized",
			},
			"fiber_performance": map[string]interface{}{
				"routing_speed":     "Lightning fast",
				"memory_footprint":  "Minimal",
				"json_serialization": "Zero allocation",
				"benchmarks_url":    "https://docs.gofiber.io/extra/benchmarks",
			},
			"comparison": map[string]string{
				"vs_express": "10x faster",
				"vs_gin": "Similar performance",
				"vs_echo": "Competitive",
			},
		})
	}))

	// Error handling examples
	server.Get("/fiber/errors/:type", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		errorType := req.Param("type")

		switch errorType {
		case "404":
			return webserver.NewResponse().
				Status(404).
				Json(map[string]interface{}{
					"error": "Not Found",
					"message": "Fiber-style 404 error handling",
					"code": "FIBER_NOT_FOUND",
					"timestamp": time.Now().Format(time.RFC3339),
					"fiber_features": []string{"Custom error pages", "Error middleware"},
				})
		case "500":
			return webserver.NewResponse().
				Status(500).
				Json(map[string]interface{}{
					"error": "Internal Server Error",
					"message": "Fiber error recovery example",
					"code": "FIBER_INTERNAL_ERROR",
					"recovery": "Fiber's built-in recovery prevents crashes",
				})
		case "timeout":
			return webserver.NewResponse().
				Status(408).
				Json(map[string]interface{}{
					"error": "Request Timeout",
					"message": "Fiber timeout handling example",
					"code": "FIBER_TIMEOUT",
					"fiber_feature": "Built-in timeout handling",
				})
		default:
			return webserver.NewResponse().Json(map[string]interface{}{
				"available_error_types": []string{"404", "500", "timeout"},
				"example": "/fiber/errors/404",
				"message": "Specify error type in URL",
			})
		}
	}))

	// Content type handling
	server.Get("/fiber/content/:type", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		contentType := req.Param("type")
		data := map[string]interface{}{
			"message": "Fiber content type example",
			"engine": "GoFiber",
			"timestamp": time.Now().Format(time.RFC3339),
			"performance": "High",
		}

		switch contentType {
		case "json":
			return webserver.NewResponse().
				Header("Content-Type", "application/json").
				Json(data)
		case "text":
			return webserver.NewResponse().
				Header("Content-Type", "text/plain").
				Text(fmt.Sprintf("Fiber: %s at %s", data["message"], data["timestamp"]))
		case "html":
			html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>Fiber Content</title></head>
<body>
	<h1>%s</h1>
	<p>Engine: %s</p>
	<p>Timestamp: %s</p>
	<p>Performance: %s</p>
</body>
</html>`, data["message"], data["engine"], data["timestamp"], data["performance"])
			return webserver.NewResponse().
				Header("Content-Type", "text/html").
				HTML(html)
		default:
			return webserver.NewResponse().Json(map[string]interface{}{
				"available_types": []string{"json", "text", "html"},
				"example": "/fiber/content/json",
				"fiber_rendering": "Supports multiple content types with high performance",
			})
		}
	}))

	// Start server with Fiber-specific information
	fmt.Println("Starting GoFiber Webserver Example...")
	fmt.Println("GoFiber Framework Features:")
	fmt.Println("  • Express-inspired API with Go performance")
	fmt.Println("  • Extremely low memory footprint")
	fmt.Println("  • Zero allocation router")
	fmt.Println("  • Fast JSON serialization")
	fmt.Println("  • Built-in middleware support")
	fmt.Println("  • Prefork support for scaling")
	fmt.Println("")
	
	// Display all registered routes with clickable URLs
	logging.DisplayRoutesClickable(server, "0.0.0.0", 8083)
	
	fmt.Println("Press Ctrl+C to stop the server")

	// Graceful shutdown simulation
	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("\nShutting down Fiber server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Fiber server shutdown error: %v", err)
		}
	}()

	// Start listening
	if err := server.Listen(":8083"); err != nil {
		log.Printf("Fiber server failed to start: %v", err)
	}

	fmt.Println("Fiber server stopped.")
}
