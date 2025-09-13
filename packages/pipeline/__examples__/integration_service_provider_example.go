package main

import (
	"fmt"
	"log"

	"govel/new/pipeline"
	"govel/new/pipeline/interfaces"
	"govel/new/pipeline/providers"
	"govel/container/mocks"
)

// RequestProcessor is a sample service that uses pipelines
type RequestProcessor struct {
	hub interfaces.HubInterface
}

func NewRequestProcessor(hub interfaces.HubInterface) *RequestProcessor {
	return &RequestProcessor{hub: hub}
}

func (rp *RequestProcessor) ProcessAPIRequest(request map[string]interface{}) (interface{}, error) {
	return rp.hub.Pipe(request, "api")
}

func (rp *RequestProcessor) ProcessWebRequest(request map[string]interface{}) (interface{}, error) {
	return rp.hub.Pipe(request, "web")
}

// ValidationMiddleware validates requests
func ValidationMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	request := passable.(map[string]interface{})
	fmt.Printf("[Validation] Validating request: %v\n", request["path"])
	
	// Basic validation
	if request["path"] == "" {
		return nil, fmt.Errorf("path is required")
	}
	
	if request["method"] == "" {
		return nil, fmt.Errorf("method is required")
	}
	
	fmt.Println("[Validation] Request is valid")
	return next(passable)
}

// AuthenticationMiddleware checks authentication
func AuthenticationMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	request := passable.(map[string]interface{})
	fmt.Printf("[Authentication] Checking auth for: %v\n", request["path"])
	
	// Check for auth header
	if auth, exists := request["auth"]; !exists || auth == "" {
		return nil, fmt.Errorf("authentication required")
	}
	
	// Simulate auth validation
	if request["auth"] != "valid-token" {
		return nil, fmt.Errorf("invalid authentication token")
	}
	
	request["user_id"] = 12345
	fmt.Println("[Authentication] User authenticated successfully")
	return next(passable)
}

// CacheMiddleware handles caching
func CacheMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	request := passable.(map[string]interface{})
	fmt.Printf("[Cache] Checking cache for: %v\n", request["path"])
	
	// Simulate cache miss (always process)
	fmt.Println("[Cache] Cache miss, processing request")
	
	result, err := next(passable)
	
	if err == nil {
		fmt.Println("[Cache] Caching result")
	}
	
	return result, err
}

// ExampleIntegrationWithServiceProvider demonstrates full integration
func ExampleIntegrationWithServiceProvider() {
	fmt.Println("=== Service Provider Integration Example ===")
	
	// Create mock container
	container := mocks.NewMockContainer()
	
	// Create and register pipeline service provider
	pipelineProvider := providers.NewPipelineServiceProvider()
	
	// Bind pipeline services to container
	container.Singleton("pipeline.hub", func() interface{} {
		return pipeline.NewHub(container)
	})
	
	container.Singleton("pipeline.pipeline", func() interface{} {
		return pipeline.NewPipeline(container)
	})
	
	fmt.Println("✓ Pipeline services registered successfully")
	
	// Show what services would be provided by the service provider
	services := pipelineProvider.Provides()
	fmt.Printf("✓ Available service provider would provide: %v\n", services)
	
	// Get hub from container
	hubInterface, err := container.Make("pipeline.hub")
	if err != nil {
		log.Fatalf("Failed to resolve hub: %v", err)
	}
	
	hub, ok := hubInterface.(interfaces.HubInterface)
	if !ok {
		log.Fatal("Hub does not implement HubInterface")
	}
	
	fmt.Println("✓ Hub resolved from container")
	
	// Configure API pipeline
	hub.Pipeline("api", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
		fmt.Println("\n--- Processing API Request ---")
		
		result, err := pipeline.
			Send(passable).
			Through([]interface{}{ValidationMiddleware, AuthenticationMiddleware, CacheMiddleware}).
			Then(func(passable interface{}) interface{} {
				request := passable.(map[string]interface{})
				fmt.Printf("[API Handler] Processing %s %s for user %v\n", 
					request["method"], request["path"], request["user_id"])
				
				return map[string]interface{}{
					"status": "success",
					"data": map[string]interface{}{
						"message": "API request processed",
						"path":    request["path"],
						"user_id": request["user_id"],
					},
				}
			})
		
		if err != nil {
			fmt.Printf("[API] Request failed: %v\n", err)
			return map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		}
		
		return result
	})
	
	// Configure Web pipeline (simpler, no auth required)
	hub.Pipeline("web", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
		fmt.Println("\n--- Processing Web Request ---")
		
		result, err := pipeline.
			Send(passable).
			Through([]interface{}{ValidationMiddleware, CacheMiddleware}).
			Then(func(passable interface{}) interface{} {
				request := passable.(map[string]interface{})
				fmt.Printf("[Web Handler] Processing %s %s\n", 
					request["method"], request["path"])
				
				return map[string]interface{}{
					"status": "success",
					"data": map[string]interface{}{
						"message": "Web request processed",
						"path":    request["path"],
					},
				}
			})
		
		if err != nil {
			fmt.Printf("[Web] Request failed: %v\n", err)
			return map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		}
		
		return result
	})
	
	// Register request processor service that uses the hub
	container.Singleton("request.processor", func() interface{} {
		hubInstance, _ := container.Make("pipeline.hub")
		return NewRequestProcessor(hubInstance.(interfaces.HubInterface))
	})
	
	fmt.Println("✓ Pipelines configured and services registered")
	
	// Get request processor from container
	processorInterface, err := container.Make("request.processor")
	if err != nil {
		log.Fatalf("Failed to resolve request processor: %v", err)
	}
	
	processor := processorInterface.(*RequestProcessor)
	fmt.Println("✓ Request processor resolved from container")
	
	// Test API request (requires authentication)
	fmt.Println("\n=== Testing API Request ===")
	apiRequest := map[string]interface{}{
		"method": "GET",
		"path":   "/api/users/123",
		"auth":   "valid-token",
	}
	
	apiResult, err := processor.ProcessAPIRequest(apiRequest)
	if err != nil {
		fmt.Printf("API request failed: %v\n", err)
	} else {
		fmt.Printf("API result: %+v\n", apiResult)
	}
	
	// Test Web request (no auth required)
	fmt.Println("\n=== Testing Web Request ===")
	webRequest := map[string]interface{}{
		"method": "GET",
		"path":   "/about",
	}
	
	webResult, err := processor.ProcessWebRequest(webRequest)
	if err != nil {
		fmt.Printf("Web request failed: %v\n", err)
	} else {
		fmt.Printf("Web result: %+v\n", webResult)
	}
	
	// Test error case - invalid auth
	fmt.Println("\n=== Testing Invalid Auth ===")
	invalidRequest := map[string]interface{}{
		"method": "POST",
		"path":   "/api/users",
		"auth":   "invalid-token",
	}
	
	invalidResult, err := processor.ProcessAPIRequest(invalidRequest)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		fmt.Printf("Invalid auth result: %+v\n", invalidResult)
	}
	
	// Show container statistics
	fmt.Println("\n=== Container Statistics ===")
	stats := container.GetStatistics()
	fmt.Printf("Container stats: %+v\n", stats)
	
	bindings := container.GetBindings()
	fmt.Printf("Registered services: %v\n", getKeys(bindings))
}

// Helper function to get map keys
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	ExampleIntegrationWithServiceProvider()
}
