package ignition

import (
	"fmt"
	"strings"

	"govel/ignition/interfaces"
	"govel/ignition/models"
)

// BadMethodCallSolutionProvider provides solutions for method-related errors
type BadMethodCallSolutionProvider struct{}

// GetSolutions returns solutions for method call errors
func (p *BadMethodCallSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "has no field or method") {
		methodName := extractMethodName(errMsg)
		solution := models.NewSolution(
			"Method does not exist",
			fmt.Sprintf("The method '%s' does not exist on this type.", methodName),
		)
		solution.AddLink("Effective Go - Methods", "https://golang.org/doc/effective_go.html#methods")
		solution.AddLink("Go Tour - Methods", "https://tour.golang.org/methods/1")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "cannot call pointer method") {
		solution := models.NewSolution(
			"Pointer receiver method called on value",
			"You're trying to call a method that requires a pointer receiver on a value. Use & to get the address.",
		)
		solution.AddLink("Go Tour - Pointer Receivers", "https://tour.golang.org/methods/4")
		solution.AddLink("Go FAQ - Methods on Values or Pointers", "https://golang.org/doc/faq#methods_on_values_or_pointers")
		solutions = append(solutions, solution)
	}

	return solutions
}

// NilPointerSolutionProvider provides solutions for nil pointer errors
type NilPointerSolutionProvider struct{}

// GetSolutions returns solutions for nil pointer errors
func (p *NilPointerSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "nil pointer dereference") ||
		strings.Contains(errMsg, "invalid memory address") {
		solution := models.NewSolution(
			"Nil pointer dereference",
			"You're trying to access a field or method on a nil pointer. Make sure to initialize the pointer or check if it's nil before using it.",
		)
		solution.AddLink("Go FAQ - Nil Errors", "https://golang.org/doc/faq#nil_error")
		solution.AddLink("Go Tour - Pointers", "https://tour.golang.org/moretypes/1")
		solutions = append(solutions, solution)
	}

	return solutions
}

// TypeMismatchSolutionProvider provides solutions for type errors
type TypeMismatchSolutionProvider struct{}

// GetSolutions returns solutions for type mismatch errors
func (p *TypeMismatchSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "cannot use") && strings.Contains(errMsg, "as type") {
		solution := models.NewSolution(
			"Type mismatch",
			"You're trying to use a value of one type where another type is expected. Check the types and use type conversion if necessary.",
		)
		solution.AddLink("Go Tour - Type Conversions", "https://tour.golang.org/basics/13")
		solution.AddLink("Effective Go - Conversions", "https://golang.org/doc/effective_go.html#conversions")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "interface conversion") {
		solution := models.NewSolution(
			"Invalid interface conversion",
			"You're trying to convert an interface to a type that doesn't implement it. Use type assertions or type switches to safely convert.",
		)
		solution.AddLink("Go Tour - Type Assertions", "https://tour.golang.org/methods/15")
		solution.AddLink("Effective Go - Interface Conversions", "https://golang.org/doc/effective_go.html#interface_conversions")
		solutions = append(solutions, solution)
	}

	return solutions
}

// ImportSolutionProvider provides solutions for import-related errors
type ImportSolutionProvider struct{}

// GetSolutions returns solutions for import errors
func (p *ImportSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "cannot find package") {
		packageName := extractPackageName(errMsg)
		solution := models.NewSolution(
			"Package not found",
			fmt.Sprintf("The package '%s' cannot be found. Make sure it's installed and the import path is correct.", packageName),
		)
		solution.AddLink("Go Documentation - Import Paths", "https://golang.org/doc/code.html#ImportPaths")
		solution.AddLink("Go Packages", "https://pkg.go.dev/")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "imported and not used") {
		solution := models.NewSolution(
			"Unused import",
			"You've imported a package but haven't used it. Remove unused imports or use the blank identifier _ to import for side effects only.",
		)
		solution.AddLink("Effective Go - Blank Identifier", "https://golang.org/doc/effective_go.html#blank_unused")
		solutions = append(solutions, solution)
	}

	return solutions
}

// ConcurrencySolutionProvider provides solutions for concurrency-related errors
type ConcurrencySolutionProvider struct{}

// GetSolutions returns solutions for concurrency errors
func (p *ConcurrencySolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "concurrent map") {
		solution := models.NewSolution(
			"Concurrent map access",
			"You're accessing a map concurrently from multiple goroutines without proper synchronization. Use a mutex or sync.Map for concurrent access.",
		)
		solution.AddLink("Go FAQ - Atomic Maps", "https://golang.org/doc/faq#atomic_maps")
		solution.AddLink("Go Sync Package", "https://pkg.go.dev/sync#Map")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "send on closed channel") {
		solution := models.NewSolution(
			"Sending on closed channel",
			"You're trying to send a value on a channel that has been closed. Check if the channel is closed before sending or use a select statement with a default case.",
		)
		solution.AddLink("Go Tour - Channels", "https://tour.golang.org/concurrency/4")
		solution.AddLink("Effective Go - Channels", "https://golang.org/doc/effective_go.html#channels")
		solutions = append(solutions, solution)
	}

	return solutions
}

// HTTPSolutionProvider provides solutions for HTTP-related errors
type HTTPSolutionProvider struct{}

// GetSolutions returns solutions for HTTP errors
func (p *HTTPSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "connection refused") {
		solution := models.NewSolution(
			"Connection refused",
			"The server refused the connection. Make sure the server is running and accessible on the specified address and port.",
		)
		solution.AddLink("Go HTTP Package", "https://pkg.go.dev/net/http")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "timeout") {
		solution := models.NewSolution(
			"Request timeout",
			"The request timed out. Consider increasing the timeout value or checking if the server is responding properly.",
		)
		solution.AddLink("Go HTTP Client", "https://pkg.go.dev/net/http#Client")
		solution.AddLink("Cloudflare - HTTP Timeouts", "https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/")
		solutions = append(solutions, solution)
	}

	return solutions
}

// DatabaseSolutionProvider provides solutions for database-related errors
type DatabaseSolutionProvider struct{}

// GetSolutions returns solutions for database errors
func (p *DatabaseSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface
	errMsg := err.Error()

	if strings.Contains(errMsg, "no rows in result set") {
		solution := models.NewSolution(
			"No rows found",
			"The query returned no results. Check your query conditions and ensure the data exists.",
		)
		solution.AddLink("Go SQL Package", "https://pkg.go.dev/database/sql#ErrNoRows")
		solutions = append(solutions, solution)
	}

	if strings.Contains(errMsg, "connection refused") && strings.Contains(errMsg, "database") {
		solution := models.NewSolution(
			"Database connection failed",
			"Could not connect to the database. Check your connection string, credentials, and ensure the database server is running.",
		)
		solution.AddLink("Go Database Package", "https://pkg.go.dev/database/sql")
		solutions = append(solutions, solution)
	}

	return solutions
}

// GoVelSolutionProvider provides solutions for GoVel-specific errors
type GoVelSolutionProvider struct{}

// GetSolutions returns solutions for GoVel errors
func (p *GoVelSolutionProvider) GetSolutions(err error) []interfaces.SolutionInterface {
	var solutions []interfaces.SolutionInterface

	// Check if it's a GoVel-specific error
	if govelErr, ok := err.(*GoVelError); ok {
		switch govelErr.Code {
		case "GOVEL_ROUTE_NOT_FOUND":
			solution := models.NewSolution(
				"Route not found",
				"The requested route is not registered. Check your route definitions and make sure the route is properly registered.",
			)
			solution.AddLink("Documentation", "https://govel#routing")
			solutions = append(solutions, solution)

		case "GOVEL_MIDDLEWARE_ERROR":
			solution := models.NewSolution(
				"Middleware error",
				"An error occurred in middleware processing. Check your middleware chain and ensure all middleware is properly configured.",
			)
			solution.AddLink("Documentation", "https://govel#middleware")
			solutions = append(solutions, solution)

		case "GOVEL_CONTAINER_BINDING_ERROR":
			solution := models.NewSolution(
				"Container binding error",
				"Failed to resolve a dependency from the container. Make sure the service is properly registered and all dependencies are satisfied.",
			)
			solution.AddLink("Documentation", "https://govel#dependency-injection")
			solutions = append(solutions, solution)
		}
	}

	return solutions
}

// GetDefaultSolutionProviders returns the default set of solution providers
func GetDefaultSolutionProviders() []interfaces.SolutionProviderInterface {
	return []interfaces.SolutionProviderInterface{
		&BadMethodCallSolutionProvider{},
		&NilPointerSolutionProvider{},
		&TypeMismatchSolutionProvider{},
		&ImportSolutionProvider{},
		&ConcurrencySolutionProvider{},
		&HTTPSolutionProvider{},
		&DatabaseSolutionProvider{},
		&GoVelSolutionProvider{},
	}
}

// Helper functions to extract information from error messages

func extractMethodName(errMsg string) string {
	// Simple extraction - in a real implementation, you'd use regex
	parts := strings.Split(errMsg, "'")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

func extractPackageName(errMsg string) string {
	// Simple extraction - in a real implementation, you'd use regex
	parts := strings.Split(errMsg, "\"")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}
