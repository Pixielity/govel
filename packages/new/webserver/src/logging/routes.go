
package logging

import (
	"fmt"
	"sort"
	"strings"

	"govel/new/routing"
	"govel/new/webserver/interfaces"
)

// DisplayRoutes logs all registered routes in a structured and easy-to-read format.
// This utility inspects the route collection, sorts the routes, and prints them
// to the console with their HTTP method, path, and name.
func DisplayRoutes(server interfaces.WebserverInterface) {
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

	fmt.Println("Available Routes:")
	for _, route := range allRoutes {
		// Clean up the path for better readability
		path := strings.ReplaceAll(route.Path, "//", "/")

		// Prepare route name (if available)
		name := ""
		if route.Name != "" {
			name = fmt.Sprintf("[%s]", route.Name)
		}

		// Format and print the route
		fmt.Printf("  %s%-8s %s %s\n", getMethodColor(route.Method.String()), route.Method.String(), path, name)
	}
	fmt.Println()
}

// getMethodColor returns an ANSI color code for a given HTTP method.
// This helps to visually distinguish methods in the console output.
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[32m" // Green
	case "POST":
		return "\033[34m" // Blue
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[35m" // Magenta
	case "OPTIONS":
		return "\033[36m" // Cyan
	default:
		return "\033[0m" // Reset
	}
}

