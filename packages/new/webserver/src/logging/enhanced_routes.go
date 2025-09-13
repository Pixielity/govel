package logging

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"govel/new/routing/src"
	"govel/new/webserver/src/interfaces"
)

// DisplayRoutesClickable logs all registered routes with clickable URLs for easy testing.
// This enhanced version formats URLs to be clickable in most terminals and includes
// descriptions and example URLs with sample parameters.
func DisplayRoutesClickable(server interfaces.WebserverInterface, host string, port int) {
	// Cast to access GetRoutes()
	webserver, ok := server.(interface{ GetRoutes() *routing.RouteCollection })
	if !ok {
		fmt.Println("Error: Could not retrieve routes from server instance.")
		return
	}

	allRoutes := webserver.GetRoutes().GetAllRoutes()

	// Sort routes for consistent ordering
	sort.Slice(allRoutes, func(i, j int) bool {
		if allRoutes[i].Path != allRoutes[j].Path {
			return allRoutes[i].Path < allRoutes[j].Path
		}
		return allRoutes[i].Method.String() < allRoutes[j].Method.String()
	})

	// Build base URL
	scheme := "http"
	if port == 443 {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	if (scheme == "http" && port != 80) || (scheme == "https" && port != 443) {
		baseURL += ":" + strconv.Itoa(port)
	}

	fmt.Printf("\nAvailable Routes on %s\n", baseURL)
	fmt.Printf("Click on the URLs below to test them:\n\n")

	// Group routes by path for better readability
	groupedRoutes := groupRoutesByPath(allRoutes)

	for path, methods := range groupedRoutes {
		printRouteGroup(baseURL, path, methods)
	}

	fmt.Println("\nTesting Tips:")
	fmt.Printf("• For routes with parameters (e.g., :id), replace with actual values\n")
	fmt.Printf("• POST/PUT/PATCH routes may require JSON body data\n")
	fmt.Printf("• Use curl for programmatic testing: curl -X METHOD %s/path\n\n", baseURL)
}

// groupRoutesByPath groups routes by their path pattern
func groupRoutesByPath(routes []*routing.Route) map[string][]*routing.Route {
	grouped := make(map[string][]*routing.Route)
	for _, route := range routes {
		// Clean up the path
		cleanPath := strings.ReplaceAll(route.Path, "//", "/")
		grouped[cleanPath] = append(grouped[cleanPath], route)
	}
	return grouped
}

// printRouteGroup prints a group of routes that share the same path
func printRouteGroup(baseURL, path string, routes []*routing.Route) {
	// Convert path parameters to example values
	examplePath := convertParamsToExample(path)
	fullURL := baseURL + examplePath

	// Create method list
	var methods []string
	for _, route := range routes {
		methodColor := getMethodColor(route.Method.String())
		methods = append(methods, fmt.Sprintf("%s%s\033[0m", methodColor, route.Method.String()))
	}

	// Print the clickable URL
	fmt.Printf("  Path: %s\n", path)
	fmt.Printf("  Methods: %s\n", strings.Join(methods, ", "))
	fmt.Printf("  Test URL: \033]8;;%s\033\\%s\033]8;;\033\\\n", fullURL, fullURL)

	// Add descriptions based on path patterns
	description := getPathDescription(path, routes)
	if description != "" {
		fmt.Printf("  Description: %s\n", description)
	}

	// Add example curl commands for non-GET methods
	addCurlExamples(fullURL, routes)
	fmt.Println()
}

// convertParamsToExample converts path parameters like :id to example values
func convertParamsToExample(path string) string {
	// Common parameter replacements
	replacements := map[string]string{
		":id":      "123",
		":userId":  "456",
		":postId":  "789",
		":type":    "example",
		":format":  "json",
		":status":  "active",
		":name":    "sample",
	}

	result := path
	for param, example := range replacements {
		result = strings.ReplaceAll(result, param, example)
	}

	// Handle wildcard routes
	result = strings.ReplaceAll(result, "/*", "/test/path")

	return result
}

// getPathDescription provides human-readable descriptions for common path patterns
func getPathDescription(path string, routes []*routing.Route) string {
	switch {
	case path == "/":
		return "Welcome/Home page"
	case path == "/health":
		return "Health check endpoint"
	case strings.Contains(path, "/api/"):
		return "API endpoint"
	case strings.Contains(path, "/admin"):
		return "Admin functionality"
	case strings.Contains(path, "/:id") && len(routes) > 0:
		method := routes[0].Method.String()
		switch method {
		case "GET":
			return "Retrieve item by ID"
		case "PUT":
			return "Update item by ID"
		case "DELETE":
			return "Delete item by ID"
		default:
			return "Operate on specific item"
		}
	case strings.Contains(path, "/users"):
		return "User management"
	case strings.Contains(path, "/search"):
		return "Search functionality"
	case strings.Contains(path, "/upload"):
		return "File upload"
	case strings.Contains(path, "/download"):
		return "File download"
	case strings.Contains(path, "/errors"):
		return "Error handling demo"
	default:
		return ""
	}
}

// addCurlExamples adds curl command examples for non-GET requests
func addCurlExamples(url string, routes []*routing.Route) {
	for _, route := range routes {
		method := route.Method.String()
		switch method {
		case "POST":
			fmt.Printf("  Example: curl -X POST -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' %s\n", url)
		case "PUT":
			fmt.Printf("  Example: curl -X PUT -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' %s\n", url)
		case "PATCH":
			fmt.Printf("  Example: curl -X PATCH -H \"Content-Type: application/json\" -d '{\"key\":\"value\"}' %s\n", url)
		case "DELETE":
			fmt.Printf("  Example: curl -X DELETE %s\n", url)
		}
	}
}

// DisplayRoutesSummary provides a concise summary of all routes
func DisplayRoutesSummary(server interfaces.WebserverInterface, host string, port int) {
	webserver, ok := server.(interface{ GetRoutes() *routing.RouteCollection })
	if !ok {
		fmt.Println("Error: Could not retrieve routes from server instance.")
		return
	}

	allRoutes := webserver.GetRoutes().GetAllRoutes()
	
	// Count methods
	methodCounts := make(map[string]int)
	for _, route := range allRoutes {
		methodCounts[route.Method.String()]++
	}

	scheme := "http"
	if port == 443 {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	if (scheme == "http" && port != 80) || (scheme == "https" && port != 443) {
		baseURL += ":" + strconv.Itoa(port)
	}

	fmt.Printf("Server: %s\n", baseURL)
	fmt.Printf("Total Routes: %d\n", len(allRoutes))
	fmt.Print("Methods: ")
	
	var methodStrings []string
	for method, count := range methodCounts {
		color := getMethodColor(method)
		methodStrings = append(methodStrings, fmt.Sprintf("%s%s(%d)\033[0m", color, method, count))
	}
	fmt.Printf("%s\n", strings.Join(methodStrings, ", "))
}
